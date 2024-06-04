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
	OrderID    string
	MerchantID string
	ItemIDs    []string
}
