package order

type CalculateOrderEstimateResponse struct {
	TotalPrice                     int    `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateID           string `json:"calculatedEstimateId"`
}
