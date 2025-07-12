package auth

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

//go:generate mockgen -destination=./mock/translator_mock.go -package=mock cosmos-server/pkg/routes/auth Translator

type Translator interface {
	ToModelUser(username, email string) *model.User

	ToAuthenticateResponse(user *model.User, token string) *api.AuthenticateResponse
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToModelUser(username, email string) *model.User {
	return &model.User{
		Username: username,
		Email:    email,
	}
}

func (t *translator) ToAuthenticateResponse(user *model.User, token string) *api.AuthenticateResponse {
	return &api.AuthenticateResponse{
		Token: token,
		User: api.User{
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}
}
