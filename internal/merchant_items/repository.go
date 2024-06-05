package merchantitems

import (
	"context"
	"fmt"
	"strings"

	"github.com/citadel-corp/belimang/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, item *MerchantItems) (err error)
	List(ctx context.Context, filter ListMerchantItemsPayload) (items []MerchantItems, err error)
	ListByUIDs(ctx context.Context, uids []string) ([]*MerchantItems, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, item *MerchantItems) (err error) {
	createItemQuery := `
		INSERT INTO merchant_items (
            uid, merchant_id, name, item_category, price, image_url
        ) VALUES (
            $1, $2, $3, $4, $5, $6
        )
    `

	_, err = d.db.DB().ExecContext(ctx, createItemQuery, item.UID, item.MerchantID, item.Name,
		item.Category, item.Price, item.ImageURL)

	return
}

func (d *dbRepository) List(ctx context.Context, filter ListMerchantItemsPayload) (items []MerchantItems, err error) {
	items = make([]MerchantItems, 0)

	q := `
		SELECT mi.uid, mi.name, mi.merchant_id, mi.item_category, mi.price, mi.image_url, mi.created_at
		FROM merchant_items mi
	`

	paramNo := 1
	params := make([]interface{}, 0)
	if filter.MerchantID != 0 {
		q += "JOIN merchants m ON mi.merchant_id = m.id "
		q += fmt.Sprintf("WHERE m.id = $%d ", paramNo)
		paramNo += 1
		params = append(params, filter.MerchantID)
	}
	if filter.ItemUID != "" {
		q += fmt.Sprintf("WHERE mi.uid = $%d ", paramNo)
		paramNo += 1
		params = append(params, filter.ItemUID)
	}
	if filter.Name != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("LOWER(mi.name) LIKE $%d ", paramNo)
		paramNo += 1
		params = append(params, "%"+strings.ToLower(filter.Name)+"%")
	}
	if filter.ProductCategory != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("mi.item_category = $%d ", paramNo)
		paramNo += 1
		params = append(params, filter.ProductCategory)
	}

	orderBy := "desc"
	if filter.CreatedAtSort == "asc" {
		orderBy = "asc"
	}

	q += fmt.Sprintf(" ORDER BY mi.created_at %s", orderBy)

	q += fmt.Sprintf(" OFFSET $%d LIMIT $%d", paramNo, paramNo+1)
	params = append(params, filter.Offset)
	params = append(params, filter.Limit)

	fmt.Println("query 1:", q)

	rows, err := d.db.DB().QueryContext(ctx, q, params...)
	if err != nil {
		return
	}

	for rows.Next() {
		m := MerchantItems{}
		err = rows.Scan(&m.UID, &m.Name, &m.MerchantID, &m.Category, &m.Price, &m.ImageURL, &m.CreatedAt)
		if err != nil {
			return
		}
		items = append(items, m)
	}
	return
}

// ListByUID implements Repository.
func (d *dbRepository) ListByUIDs(ctx context.Context, uids []string) ([]*MerchantItems, error) {
	if len(uids) == 0 {
		return make([]*MerchantItems, 0), nil
	}
	q := `
		SELECT mi.uid, mi.name, mi.merchant_id, mi.item_category, mi.price, mi.image_url, mi.created_at
		FROM merchant_items mi
		WHERE mi.uid IN(
	`

	for i, v := range uids {
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
	res := make([]*MerchantItems, 0)
	for rows.Next() {
		m := &MerchantItems{}
		err = rows.Scan(&m.UID, &m.Name, &m.MerchantID, &m.Category, &m.Price, &m.ImageURL, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func whereOrAnd(paramNo int) string {
	if paramNo == 1 {
		return "WHERE "
	}
	return "AND "
}
