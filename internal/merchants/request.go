package merchants

import (
	validations "github.com/citadel-corp/belimang/internal/common/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Location struct {
	Lat *float64 `json:"lat"`
	Lng *float64 `json:"long"`
}

type CreateMerchantPayload struct {
	Name     string           `json:"name"`
	Category MerchantCategory `json:"merchantCategory"`
	ImageURL string           `json:"imageUrl"`
	Location *Location        `json:"location"`
}

func (p CreateMerchantPayload) Validate() error {
	// lat, lng := fmt.Sprintf("%f", p.Location.Lat), fmt.Sprintf("%f", p.Location.Lng)
	if err := validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(MinName, MaxName)),
		validation.Field(&p.Category, validation.Required, validation.In(MerchantCategories...)),
		validation.Field(&p.ImageURL, validation.Required, validations.ImgUrlValidationRule),
		validation.Field(&p.Location, validation.NotNil),
		// validation.Field(&p.Location.Lat, validation.Required, is.Latitude),
		// validation.Field(&lng, validation.Required, is.Longitude),
	); err != nil {
		return err
	}
	return p.Location.Validate()
}

func (p Location) Validate() error {
	if err := validation.ValidateStruct(&p,
		validation.Field(&p.Lat, validation.NotNil),
		validation.Field(&p.Lng, validation.NotNil),
	); err != nil {
		return err
	}
	if !validations.LatitudeValidation(*p.Lat) {
		return validation.NewError("latitude", "latitude is not valid")
	}
	if !validations.LongitudeValidation(*p.Lng) {
		return validation.NewError("longitude", "longitude is not valid")
	}

	return nil
}

type ListMerchantsPayload struct {
	MerchantUID      string           `schema:"merchantId" binding:"omitempty"`
	Name             string           `schema:"name" binding:"omitempty"`
	MerchantCategory MerchantCategory `schema:"merchantCategory"`
	CreatedAtSort    string           `schema:"createdAt" binding:"omitempty"`
	Limit            int              `schema:"limit" binding:"omitempty"`
	Offset           int              `schema:"offset" binding:"omitempty"`
}

func (p ListMerchantsPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.MerchantUID),
		validation.Field(&p.Name),
		validation.Field(&p.MerchantCategory, validation.In(MerchantCategories...)),
		validation.Field(&p.CreatedAtSort, validation.In([]interface{}{"asc", "desc"}...)),
	)
}

type ListMerchantsByDistancePayload struct {
	MerchantUID      string           `schema:"merchantId" binding:"omitempty"`
	Name             string           `schema:"name" binding:"omitempty"`
	MerchantCategory MerchantCategory `schema:"merchantCategory"`
	Lat              string
	Lng              string
	Limit            int `schema:"limit" binding:"omitempty"`
	Offset           int `schema:"offset" binding:"omitempty"`
}

func (p ListMerchantsByDistancePayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.MerchantUID),
		validation.Field(&p.Name),
		validation.Field(&p.MerchantCategory, validation.In(MerchantCategories...)),
		validation.Field(&p.Lat, validation.Required, is.Latitude),
		validation.Field(&p.Lng, validation.Required, is.Longitude),
	)
}
