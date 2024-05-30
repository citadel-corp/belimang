package order

import (
	"context"

	"github.com/citadel-corp/belimang/internal/common/db"
)

type Repository interface {
	InsertCalculatedEstimate(ctx context.Context, calculatedEstimate *CalculatedEstimate) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// InsertCalculatedEstimate implements Repository.
func (d *dbRepository) InsertCalculatedEstimate(ctx context.Context, calculatedEstimate *CalculatedEstimate) error {
	q := `
	    INSERT INTO calculated_estimates  (
            id, user_id, total_price, user_location_lat, user_location_lng, estimated_delivery_time, ordered
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7
        )
	`
	_, err := d.db.DB().ExecContext(ctx, q, calculatedEstimate.ID, calculatedEstimate.UserID, calculatedEstimate.TotalPrice, calculatedEstimate.Lat, calculatedEstimate.Long, calculatedEstimate.EstimatedDeliveryTime, calculatedEstimate.Ordered)
	if err != nil {
		return err
	}
	return nil
}
