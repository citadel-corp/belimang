package order

import "time"

type CalculatedEstimate struct {
	ID                    string
	UserID                string
	TotalPrice            int
	Lat                   float64
	Long                  float64
	EstimatedDeliveryTime int
	Ordered               bool
	CreatedAt             time.Time
}
