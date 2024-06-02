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
