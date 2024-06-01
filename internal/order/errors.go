package order

import "errors"

var (
	ErrValidationFailed     = errors.New("validation failed")
	ErrStartingPointInvalid = errors.New("starting point must be exactly 1")
	ErrSomeMerchantNotFound = errors.New("some merchants are not found")
	ErrSomeItemNotFound     = errors.New("some items are not found")
	ErrDistanceTooFar       = errors.New("distance too far")
)
