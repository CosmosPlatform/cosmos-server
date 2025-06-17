package api

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type AuthenticateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r *AuthenticateRequest) Validate() error {
	return validation.ValidateStruct(&AuthenticateRequest{},
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

type AuthenticateResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
