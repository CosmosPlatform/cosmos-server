package api

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type AuthenticateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r *AuthenticateRequest) Validate() string {
	return validation.ValidateStruct(&AuthenticateRequest{},
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	).Error()
}
