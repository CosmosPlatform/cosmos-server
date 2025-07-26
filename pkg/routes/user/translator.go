package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToRegisterUserResponse(username, email, role string) *api.RegisterUserResponse
	ToGetUsersResponse(users []*model.User) *api.GetUsersResponse
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToRegisterUserResponse(username, email, role string) *api.RegisterUserResponse {
	return &api.RegisterUserResponse{
		User: api.User{
			Username: username,
			Email:    email,
			Role:     role,
		},
	}
}

func (t *translator) ToGetUsersResponse(users []*model.User) *api.GetUsersResponse {
	apiUsers := make([]*api.User, 0)
	for _, user := range users {
		apiUsers = append(apiUsers, &api.User{
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		})
	}
	return &api.GetUsersResponse{
		Users: apiUsers,
	}
}
