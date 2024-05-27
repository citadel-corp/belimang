package merchants

import (
	validations "github.com/citadel-corp/belimang/internal/common/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"long"`
}

type CreateMerchantPayload struct {
	Name     string           `json:"name"`
	Category MerchantCategory `json:"merchantCategory"`
	ImageURL string           `json:"imageUrl"`
	Location Location         `json:"location"`
}

func (p CreateMerchantPayload) Validate() error {
	// lat, lng := fmt.Sprintf("%f", p.Location.Lat), fmt.Sprintf("%f", p.Location.Lng)
	if err := p.Location.Validate(); err != nil {
		return err
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(MinName, MaxName)),
		validation.Field(&p.Category, validation.Required, validation.In(MerchantCategories...)),
		validation.Field(&p.ImageURL, validation.Required, validations.ImgUrlValidationRule),
		validation.Field(&p.Location, validation.Required),
		// validation.Field(&p.Location.Lat, validation.Required, is.Latitude),
		// validation.Field(&lng, validation.Required, is.Longitude),
	)
}

func (p Location) Validate() error {
	if !validations.LatitudeValidation(p.Lat) {
		return validation.NewError("latitude", "latitude is not valid")
	}
	if !validations.LongitudeValidation(p.Lng) {
		return validation.NewError("longitude", "longitude is not valid")
	}

	return nil
}
