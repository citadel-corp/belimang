package order

import "time"

type Order struct {
	ID                   string
	CalculatedEstimateID string
	UserID               string
	CreatedAt            time.Time
}

type OrderItem struct {
	ID         string
	UserID     string
	OrderID    string
	MerchantID string
	Items      Items
}
