package api

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type AuthenticateRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
