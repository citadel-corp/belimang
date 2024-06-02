package merchantitems

import (
	validations "github.com/citadel-corp/belimang/internal/common/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateMerchantItemPayload struct {
	Name            string       `json:"name"`
	ProductCategory ItemCategory `json:"productCategory"`
	Price           int          `json:"price"`
	ImageURL        string       `json:"imageUrl"`
	MerchantID      string
}

func (p CreateMerchantItemPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(MinName, MaxName)),
		validation.Field(&p.ProductCategory, validation.Required, validation.In(ProductCategories...)),
		validation.Field(&p.Price, validation.Required, validation.Min(1)),
		validation.Field(&p.ImageURL, validation.Required, validations.ImgUrlValidationRule),
	)
}

type ListMerchantItemsPayload struct {
	ItemUID         string       `schema:"itemId" binding:"omitempty"`
	Name            string       `schema:"name" binding:"omitempty"`
	ProductCategory ItemCategory `schema:"productCategory" binding:"omitempty"`
	CreatedAtSort   string       `schema:"createdAt" binding:"omitempty"`
	Limit           int          `schema:"limit" binding:"omitempty"`
	Offset          int          `schema:"offset" binding:"omitempty"`
	MerchantUID     string
	MerchantID      uint64
}

func (p ListMerchantItemsPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ItemUID),
		validation.Field(&p.Name),
		validation.Field(&p.ProductCategory, validation.In(ProductCategories...)),
		validation.Field(&p.CreatedAtSort, validation.In([]interface{}{"asc", "desc"}...)),
		validation.Field(&p.MerchantUID, validation.Required),
	)
}
