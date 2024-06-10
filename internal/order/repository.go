package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/citadel-corp/belimang/internal/common/db"
)

type Repository interface {
	InsertCalculatedEstimate(ctx context.Context, calculatedEstimate *CalculatedEstimate) error
	GetCalculatedEstimate(ctx context.Context, id string) (*CalculatedEstimate, error)
	InsertOrder(ctx context.Context, order *Order) error
	InsertOrderItem(ctx context.Context, orderItem *OrderItem) error
	ListOrdersByUserID(ctx context.Context, userID string) (*Order, error)
	ListOrderItemsByOrderID(ctx context.Context, orderID string) ([]*OrderItem, error)
	SearchOrder(ctx context.Context, req SearchOrderPayload, userID string) ([]*searchOrderQueryResult, error)
	SearchOrderItemMerchants(ctx context.Context, req SearchOrderPayload, userID string) ([]*searchOrderItemMerchantsQueryResult, error)
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

// ListOrdersByUserID implements Repository.
func (d *dbRepository) ListOrdersByUserID(ctx context.Context, userID string) (*Order, error) {
	q := `
	    SELECT id, calculated_estimate_id, user_id
		FROM orders
		WHERE user_id = ?;
	`
	row := d.db.DB().QueryRowContext(ctx, q, userID)
	o := &Order{}
	err := row.Scan(&o.ID, &o.CalculatedEstimateID, &o.UserID)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// ListOrderItemsByOrderID implements Repository.
func (d *dbRepository) ListOrderItemsByOrderID(ctx context.Context, orderID string) ([]*OrderItem, error) {
	q := `
	    SELECT id, order_id, merchant_id, items
		FROM order_items
		WHERE order_id = ?;
	`
	rows, err := d.db.DB().QueryContext(ctx, q, orderID)
	if err != nil {
		return nil, err
	}
	res := make([]*OrderItem, 0)
	for rows.Next() {
		o := &OrderItem{}
		err = rows.Scan(&o.ID, &o.OrderID, &o.MerchantID, &o.Items)
		if err != nil {
			return nil, err
		}
		res = append(res, o)
	}

	return res, nil
}

// SearchOrder implements Repository.
func (d *dbRepository) SearchOrder(ctx context.Context, req SearchOrderPayload, userID string) ([]*searchOrderQueryResult, error) {
	paramNo := 1
	listQuery := `
		SELECT oi.order_id, oi.merchant_id, oi.items, mi.uid, mi.name, mi.item_category, mi.price, order_item_detail.quantity, mi.image_url, mi.created_at, m.uid, m.name, m.merchant_category, m.image_url, m.location_lat, m.location_lng, m.created_at
		FROM order_items oi
		CROSS JOIN LATERAL jsonb_array_elements(oi.items) as order_item_detail
		INNER JOIN merchant_items mi on order_item_detail.item_id = mi.uid
		INNER JOIN merchants m on oi.merchant_id = m.id
		WHERE `
	params := make([]interface{}, 0)
	if req.MerchantID != "" {
		listQuery += fmt.Sprintf("oi.merchant_id = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.MerchantID)
	}
	if req.Name != "" {
		listQuery += fmt.Sprintf("m.name ILIKE %%$%d%% AND ", paramNo)
		paramNo += 1
		params = append(params, req.Name)
		listQuery += fmt.Sprintf("mi.name ILIKE %%$%d%% AND ", paramNo)
		paramNo += 1
		params = append(params, req.Name)
	}
	if req.MerchantCategory != "" {
		listQuery += fmt.Sprintf("m.merchant_category = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.MerchantCategory)
	}
	listQuery += fmt.Sprintf("oi.user_id = $%d ", paramNo)
	params = append(params, userID)
	if strings.HasSuffix(listQuery, "AND ") {
		listQuery, _ = strings.CutSuffix(listQuery, "AND ")
	}
	listQuery += fmt.Sprintf(" LIMIT %d OFFSET %d;", req.Limit, req.Offset)
	rows, err := d.db.DB().QueryContext(ctx, listQuery, params...)
	if err != nil {
		return nil, err
	}
	res := make([]*searchOrderQueryResult, 0)
	for rows.Next() {
		o := &searchOrderQueryResult{}
		err = rows.Scan(&o.ItemID, &o.ItemName, &o.ItemCategory, &o.ItemPrice, &o.ItemQuantity, &o.ItemImageURL, &o.ItemCreatedAt, &o.MerchantID, &o.MerchantName, &o.MerchantCategory, &o.MerchantImageURL, &o.MerchantLat, &o.MerchantLong, &o.MerchantLong)
		if err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

// SearchOrderItemMerchants implements Repository.
func (d *dbRepository) SearchOrderItemMerchants(ctx context.Context, req SearchOrderPayload, userID string) ([]*searchOrderItemMerchantsQueryResult, error) {
	paramNo := 1
	listQuery := `
		SELECT o.id, oi.items, mi.uid, mi.name, mi.item_category, mi.price, order_item_detail.quantity, mi.image_url, mi.created_at, m.uid, m.name, m.merchant_category, m.image_url, m.location_lat, m.location_lng, m.created_at
		FROM order_items oi
		INNER JOIN orders o on oi.order_id = o.id
		INNER JOIN merchants m on oi.merchant_id = m.uid
		WHERE `
	params := make([]interface{}, 0)
	if req.MerchantID != "" {
		listQuery += fmt.Sprintf("m.uid = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.MerchantID)
	}
	if req.Name != "" {
		listQuery += fmt.Sprintf("m.name ILIKE %%$%d%% AND ", paramNo)
		paramNo += 1
		params = append(params, req.Name)
	}
	if req.MerchantCategory != "" {
		listQuery += fmt.Sprintf("m.merchant_category = $%d AND ", paramNo)
		paramNo += 1
		params = append(params, req.MerchantCategory)
	}
	listQuery += fmt.Sprintf("oi.user_id = $%d ", paramNo)
	params = append(params, userID)
	if strings.HasSuffix(listQuery, "AND ") {
		listQuery, _ = strings.CutSuffix(listQuery, "AND ")
	}
	listQuery += fmt.Sprintf(" LIMIT %d OFFSET %d;", req.Limit, req.Offset)
	rows, err := d.db.DB().QueryContext(ctx, listQuery, params...)
	if err != nil {
		return nil, err
	}
	res := make([]*searchOrderItemMerchantsQueryResult, 0)
	for rows.Next() {
		o := &searchOrderItemMerchantsQueryResult{}
		err = rows.Scan(&o.OrderID, &o.OrderItems, &o.OrderID, &o.MerchantID, &o.MerchantName, &o.MerchantCategory, &o.MerchantImageURL, &o.MerchantLat, &o.MerchantLong, &o.MerchantCreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

type searchOrderItemMerchantsQueryResult struct {
	OrderID           string
	OrderItems        Item
	MerchantID        string
	MerchantName      string
	MerchantCategory  string
	MerchantImageURL  string
	MerchantLat       float64
	MerchantLong      float64
	MerchantCreatedAt time.Time
}

type searchOrderQueryResult struct {
	OrderID           string
	OrderItems        Item
	MerchantID        string
	MerchantName      string
	MerchantCategory  string
	MerchantImageURL  string
	MerchantLat       float64
	MerchantLong      float64
	MerchantCreatedAt time.Time
	ItemID            string
	ItemName          string
	ItemCategory      string
	ItemPrice         int
	ItemQuantity      int
	ItemImageURL      string
	ItemCreatedAt     time.Time
}
