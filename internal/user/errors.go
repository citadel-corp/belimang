package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrWrongPassword     = errors.New("wrong password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrValidationFailed  = errors.New("validation failed")
)
