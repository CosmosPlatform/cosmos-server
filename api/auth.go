package api

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type AuthenticateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *AuthenticateRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

type AuthenticateResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
