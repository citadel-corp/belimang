package order

import (
	merchantitems "github.com/citadel-corp/belimang/internal/merchant_items"
	"github.com/citadel-corp/belimang/internal/merchants"
)

type CalculateOrderEstimateResponse struct {
	TotalPrice                     int    `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateID           string `json:"calculatedEstimateId"`
}

type CreateOrderResponse struct {
	OrderID string `json:"orderId"`
}

type SearchOrderResponse struct {
	OrderID string                    `json:"orderId"`
	Orders  SearchOrderDetailResponse `json:"orders"`
}

type SearchOrderDetailResponse struct {
	Merchant merchants.MerchantsResponse     `json:"merchant"`
	Items    []SearchOrderDetailItemResponse `json:"items"`
}

type SearchOrderDetailItemResponse struct {
	merchantitems.MerchantItemResponse
	Quantity int `json:"quantity"`
}
