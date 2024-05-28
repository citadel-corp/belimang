package merchants

import "errors"

var (
	ErrMerchantNotFound = errors.New("merchant not found")
	ErrValidationFailed = errors.New("validation failed")
)
