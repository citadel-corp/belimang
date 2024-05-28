package merchants

import (
	"context"

	"github.com/citadel-corp/belimang/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, merchant *Merchants) (err error)
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
