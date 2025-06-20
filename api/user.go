package api

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (r *RegisterUserRequest) Validate() error {
	return validation.ValidateStruct(&RegisterUserRequest{},
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Email, validation.Required, validation.Length(5, 100), is.EmailFormat),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 100)),
	)
}

type RegisterUserResponse struct {
	User User `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
