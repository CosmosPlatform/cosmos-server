package api

import (
	"cosmos-server/pkg/services/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role,omitempty"`
}

func (r *RegisterUserRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Email, validation.Required, validation.Length(5, 100), is.EmailFormat),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 100)),
		validation.Field(&r.Role, validation.Required, validation.In(user.AdminUserRole, user.RegularUserRole)),
	)
}

type RegisterUserResponse struct {
	User User `json:"user"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Team     *Team  `json:"team,omitempty"`
}

type GetUsersResponse struct {
	Users []*User `json:"users"`
}

type GetCurrentUserResponse struct {
	User User `json:"user"`
}
