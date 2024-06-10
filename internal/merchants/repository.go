package merchants

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/citadel-corp/belimang/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, merchant *Merchants) (err error)
	ListByUIDs(ctx context.Context, ids []string) ([]*Merchants, error)
	List(ctx context.Context, filter ListMerchantsPayload) (merchants []Merchants, err error)
	GetByUID(ctx context.Context, uid string) (merchant *Merchants, err error)
	ListByDistance(ctx context.Context, filter ListMerchantsByDistancePayload) (merchants []Merchants, err error)
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

func (d *dbRepository) List(ctx context.Context, filter ListMerchantsPayload) (merchants []Merchants, err error) {
	merchants = make([]Merchants, 0)

	q := `
		SELECT m.uid, m.name, m.merchant_category, m.image_url, m.location_lat, m.location_lng, m.created_at
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
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("LOWER(m.name) LIKE $%d ", paramNo)
		paramNo += 1
		params = append(params, "%"+strings.ToLower(filter.Name)+"%")
	}
	if filter.MerchantCategory != "" {
		q += whereOrAnd(paramNo)
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

	for rows.Next() {
		m := Merchants{}
		err = rows.Scan(&m.UID, &m.Name, &m.Category, &m.ImageURL, &m.Lat, &m.Lng, &m.CreatedAt)
		if err != nil {
			return
		}
		merchants = append(merchants, m)
	}
	return
}

func (d *dbRepository) ListByDistance(ctx context.Context, filter ListMerchantsByDistancePayload) (merchants []Merchants, err error) {
	merchants = make([]Merchants, 0)

	q := `
		SELECT earth_distance(
			ll_to_earth(m.location_lat, m.location_lng),
			ll_to_earth($1, $2)
		) as distance, 
		m.uid, m.name, m.merchant_category, m.image_url, m.location_lat, m.location_lng, m.created_at
		FROM merchants m 
	`

	paramNo := 3
	params := make([]interface{}, 0)
	params = append(params, filter.Lat)
	params = append(params, filter.Lng)

	if filter.MerchantUID != "" {
		q += fmt.Sprintf("WHERE m.uid = $%d ", paramNo)
		paramNo += 1
		params = append(params, filter.MerchantUID)
	}
	if filter.Name != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("LOWER(m.name) LIKE $%d ", paramNo)
		paramNo += 1
		params = append(params, "%"+strings.ToLower(filter.Name)+"%")
	}
	if filter.MerchantCategory != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("m.merchant_category = $%d ", paramNo)
		paramNo += 1
		params = append(params, filter.MerchantCategory)
	}

	q += " ORDER BY distance asc"

	q += fmt.Sprintf(" OFFSET $%d LIMIT $%d", paramNo, paramNo+1)
	params = append(params, filter.Offset)
	params = append(params, filter.Limit)

	rows, err := d.db.DB().QueryContext(ctx, q, params...)
	if err != nil {
		return
	}

	for rows.Next() {
		m := Merchants{}
		var distance float64
		err = rows.Scan(&distance, &m.UID, &m.Name, &m.Category, &m.ImageURL, &m.Lat, &m.Lng, &m.CreatedAt)
		if err != nil {
			return
		}
		merchants = append(merchants, m)
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

func whereOrAnd(paramNo int) string {
	if paramNo == 1 {
		return "WHERE "
	}
	return "AND "
}
