package order

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CalculateOrderEstimateRequest struct {
	UserLocation UserLocationRequest `json:"userLocation"`
	Orders       []OrderRequest      `json:"orders"`
}

func (p CalculateOrderEstimateRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserLocation, validation.Required),
		validation.Field(&p.Orders, validation.Required),
	)
}

type UserLocationRequest struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (p UserLocationRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Lat, validation.Required),
		validation.Field(&p.Long, validation.Required),
	)
}

type OrderRequest struct {
	MerchantID      string             `json:"merchantId"`
	IsStartingPoint *bool              `json:"isStartingPoint"`
	Items           []OrderItemRequest `json:"items"`
}

func (p OrderRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.MerchantID, validation.Required),
		validation.Field(&p.IsStartingPoint, validation.NotNil),
	)
}

type OrderItemRequest struct {
	ItemID   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}

func (p OrderItemRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ItemID, validation.Required),
		validation.Field(&p.Quantity, validation.Required),
	)
}

type CreateOrderRequest struct {
	CalculatedEstimateID string `json:"calculatedEstimateId"`
}

func (p CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.CalculatedEstimateID, validation.Required),
	)
}

type SearchOrderPayload struct {
	MerchantID       string `schema:"merchantId" binding:"omitempty"`
	Name             string `schema:"name" binding:"omitempty"`
	MerchantCategory string `schema:"merchantCategory " binding:"omitempty"`
	Limit            int    `schema:"limit" binding:"omitempty"`
	Offset           int    `schema:"offset" binding:"omitempty"`
}
