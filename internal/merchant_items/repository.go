package merchantitems

import (
	"context"

	"github.com/citadel-corp/belimang/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, item *MerchantItems) (err error)
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
