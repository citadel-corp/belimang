package merchants

import (
	"context"
	"fmt"

	"github.com/citadel-corp/belimang/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, merchant *Merchants) (err error)
	ListByUIDs(ctx context.Context, ids []string) ([]*Merchants, error)
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
