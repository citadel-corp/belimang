package order

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/belimang/internal/common/db"
)

type Repository interface {
	InsertCalculatedEstimate(ctx context.Context, calculatedEstimate *CalculatedEstimate) error
	GetCalculatedEstimate(ctx context.Context, id string) (*CalculatedEstimate, error)
	InsertOrder(ctx context.Context, order *Order) error
	InsertOrderItem(ctx context.Context, orderItem *OrderItem) error
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
            id, user_id, total_price, user_location_lat, user_location_lng, estimated_delivery_time, ordered, merchants, items
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        );
	`
	_, err := d.db.DB().ExecContext(ctx, q, calculatedEstimate.ID, calculatedEstimate.UserID, calculatedEstimate.TotalPrice, calculatedEstimate.Lat, calculatedEstimate.Long, calculatedEstimate.EstimatedDeliveryTime, calculatedEstimate.Ordered, calculatedEstimate.Merchants, calculatedEstimate.Items)
	if err != nil {
		return err
	}
	return nil
}

// GetCalculatedEstimate implements Repository.
func (d *dbRepository) GetCalculatedEstimate(ctx context.Context, id string) (*CalculatedEstimate, error) {
	q := `
	    SELECT id, user_id, total_price, user_location_lat, user_location_lng, estimated_delivery_time, ordered, merchants, items
		FROM calculated_estimates
        WHERE id = ?;
	`
	row := d.db.DB().QueryRowContext(ctx, q, id)
	calculatedEstimate := &CalculatedEstimate{}
	err := row.Scan(&calculatedEstimate.ID, &calculatedEstimate.UserID, &calculatedEstimate.TotalPrice, &calculatedEstimate.Lat, &calculatedEstimate.Long, &calculatedEstimate.EstimatedDeliveryTime, &calculatedEstimate.Ordered, &calculatedEstimate.Ordered, &calculatedEstimate.Merchants, &calculatedEstimate.Items)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCalculatedEstimateNotFound
	}
	if err != nil {
		return nil, err
	}
	return calculatedEstimate, nil
}

// InsertOrder implements Repository.
func (d *dbRepository) InsertOrder(ctx context.Context, order *Order) error {
	q := `
	    INSERT INTO orders (
            id, calculated_estimate_id, user_id
        ) VALUES (
            $1, $2, $3
        );
	`
	_, err := d.db.DB().ExecContext(ctx, q, order.ID, order.CalculatedEstimateID, order.UserID)
	if err != nil {
		return err
	}
	return nil
}

// InsertOrderItem implements Repository.
func (d *dbRepository) InsertOrderItem(ctx context.Context, orderItem *OrderItem) error {
	q := `
	    INSERT INTO order_items (
            id, order_id, merchant_id, items
        ) VALUES (
            $1, $2, $3, $4
        );
	`
	_, err := d.db.DB().ExecContext(ctx, q, orderItem.ID, orderItem.OrderID, orderItem.MerchantID, orderItem.Items)
	if err != nil {
		return err
	}
	return nil
}
