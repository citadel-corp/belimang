package merchants

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/citadel-corp/belimang/internal/common/db"
	"github.com/citadel-corp/belimang/internal/common/response"
)

type Repository interface {
	Create(ctx context.Context, merchant *Merchants) (err error)
	ListByUIDs(ctx context.Context, ids []string) ([]*Merchants, error)
	List(ctx context.Context, filter ListMerchantsPayload) (merchants []Merchants, pagination *response.Pagination, err error)
	GetByUID(ctx context.Context, uid string) (merchant *Merchants, err error)
	ListByDistance(ctx context.Context, filter ListMerchantsByDistancePayload) (merchantWithItem []MerchantsWithItem, pagination *response.Pagination, err error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, merchant *Merchants) (err error) {
	createMerchantQuery := `
	    INSERT INTO merchants (
            uid, name, merchant_category, image_url, location_lat, location_lng
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        )
	`
	_, err = d.db.DB().ExecContext(ctx, createMerchantQuery, merchant.UID, merchant.Name, merchant.Category, merchant.ImageURL, merchant.Lat, merchant.Lng)
	if err != nil {
		return
	}
	return
}

// ListByUIDs implements Repository.
func (d *dbRepository) ListByUIDs(ctx context.Context, ids []string) ([]*Merchants, error) {
	q := `
	    SELECT id, uid, name, merchant_category, image_url, location_lat, location_lng
		FROM merchants
		WHERE uid IN (
	`
	for i, v := range ids {
		if i > 0 {
			q += ","
		}
		q += fmt.Sprintf("'%s'", v)

	}

	q += ");"
	rows, err := d.db.DB().QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*Merchants, 0)
	for rows.Next() {
		m := &Merchants{}
		err := rows.Scan(&m.ID, &m.UID, &m.Name, &m.Category, &m.ImageURL, &m.Lat, &m.Lng)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func (d *dbRepository) List(ctx context.Context, filter ListMerchantsPayload) (merchants []Merchants, pagination *response.Pagination, err error) {
	merchants = make([]Merchants, 0)

	q := `
		SELECT COUNT(*) OVER() AS total_count, m.uid, m.name, m.merchant_category, m.image_url, m.location_lat, m.location_lng, m.created_at
		FROM merchants m
	`

	paramNo := 1
	params := make([]interface{}, 0)
	if filter.MerchantUID != "" {
		q += fmt.Sprintf("WHERE m.uid = $%d ", paramNo)
		paramNo += 1
		params = append(params, filter.MerchantUID)
	}
	if filter.Name != "" {
		q += whereOrAnd(paramNo, 1)
		q += fmt.Sprintf("LOWER(m.name) LIKE $%d ", paramNo)
		paramNo += 1
		params = append(params, "%"+strings.ToLower(filter.Name)+"%")
	}
	if filter.MerchantCategory != "" {
		q += whereOrAnd(paramNo, 1)
		q += fmt.Sprintf("m.merchant_category = $%d ", paramNo)
		paramNo += 1
		params = append(params, filter.MerchantCategory)
	}

	orderBy := "desc"
	if filter.CreatedAtSort == "asc" {
		orderBy = "asc"
	}

	q += fmt.Sprintf(" ORDER BY m.created_at %s", orderBy)

	q += fmt.Sprintf(" OFFSET $%d LIMIT $%d", paramNo, paramNo+1)
	params = append(params, filter.Offset)
	params = append(params, filter.Limit)

	rows, err := d.db.DB().QueryContext(ctx, q, params...)
	if err != nil {
		return
	}

	pagination = &response.Pagination{}
	pagination.Limit = filter.Limit
	pagination.Offset = filter.Offset

	for rows.Next() {
		m := Merchants{}
		err = rows.Scan(&pagination.Total, &m.UID, &m.Name, &m.Category, &m.ImageURL, &m.Lat, &m.Lng, &m.CreatedAt)
		if err != nil {
			return
		}
		merchants = append(merchants, m)
	}
	return
}

func (d *dbRepository) ListByDistance(ctx context.Context, filter ListMerchantsByDistancePayload) (merchantWithItem []MerchantsWithItem, pagination *response.Pagination, err error) {
	merchantWithItem = make([]MerchantsWithItem, 0)

	q := `
		SELECT COUNT(*) OVER() AS total_count,
		m.*,
		COALESCE(mi.uid, ''),
		COALESCE(mi.name, ''),
		COALESCE(mi.merchant_id, 0),
		mi.item_category,
		COALESCE(mi.price, 0),
		COALESCE(mi.image_url, ''),
		mi.created_at
		FROM (
			SELECT earth_distance(
				ll_to_earth(location_lat, location_lng),
				ll_to_earth($1, $2)
			) as distance,
			id, uid, name, merchant_category, image_url, location_lat, location_lng, created_at
			FROM merchants
			ORDER BY distance ASC
			OFFSET $3 LIMIT $4
		) AS m
		LEFT JOIN merchant_items mi ON m.id = mi.merchant_id

	`

	paramNo := 5
	params := make([]interface{}, 0)
	params = append(params, filter.Lat)
	params = append(params, filter.Lng)
	params = append(params, filter.Offset)
	params = append(params, filter.Limit)

	if filter.MerchantUID != "" {
		q += whereOrAnd(paramNo, 5)
		q += fmt.Sprintf("m.uid = $%d ", paramNo)
		paramNo += 1
		params = append(params, filter.MerchantUID)
	}
	if filter.Name != "" {
		q += whereOrAnd(paramNo, 5)
		q += fmt.Sprintf("m.name ILIKE '%%%s%%' OR mi.name ILIKE '%%%s%%'", filter.Name, filter.Name)
	}
	if filter.MerchantCategory != "" {
		q += whereOrAnd(paramNo, 5)
		q += fmt.Sprintf("m.merchant_category = $%d ", paramNo)
		params = append(params, filter.MerchantCategory)
	}

	rows, err := d.db.DB().QueryContext(ctx, q, params...)
	if err != nil {
		return
	}

	pagination = &response.Pagination{}
	pagination.Limit = filter.Limit
	pagination.Offset = filter.Offset

	for rows.Next() {
		m := Merchants{}
		mi := MerchantItems{}
		var distance float64
		err = rows.Scan(&pagination.Total, &distance, &m.ID, &m.UID, &m.Name, &m.Category, &m.ImageURL, &m.Lat, &m.Lng, &m.CreatedAt,
			&mi.UID, &mi.Name, &mi.MerchantID, &mi.Category, &mi.Price, &mi.ImageURL, &mi.CreatedAt)
		if err != nil {
			return
		}
		merchantWithItem = append(merchantWithItem, MerchantsWithItem{
			ID:        m.ID,
			UID:       m.UID,
			Name:      m.Name,
			Category:  m.Category,
			ImageURL:  m.ImageURL,
			Lat:       m.Lat,
			Lng:       m.Lng,
			CreatedAt: m.CreatedAt,
			Item:      mi,
		})
	}
	return
}

func (d *dbRepository) GetByUID(ctx context.Context, uid string) (merchant *Merchants, err error) {
	getMerchantQuery := `
		SELECT id, uid, name, merchant_category, image_url, location_lat, location_lng, created_at
		FROM merchants
		WHERE uid = $1
	`

	row := d.db.DB().QueryRowContext(ctx, getMerchantQuery, uid)
	m := Merchants{}
	err = row.Scan(&m.ID, &m.UID, &m.Name, &m.Category, &m.ImageURL, &m.Lat, &m.Lng, &m.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrMerchantNotFound
		}
		return
	}

	merchant = &m

	return
}

func whereOrAnd(paramNo int, targetParamNo int) string {
	if paramNo == targetParamNo {
		return "WHERE "
	}
	return "AND "
}
