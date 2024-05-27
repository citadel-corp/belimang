package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateUserPayload struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	UserType UserType `json:"-"`
}

func (p CreateUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Username, validation.Required, validation.Length(MinUsername, MaxUsername)),
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
		validation.Field(&p.UserType, validation.Required, validation.In(UserTypes...)),
	)
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p LoginPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Username, validation.Required, validation.Length(5, 15)),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}
