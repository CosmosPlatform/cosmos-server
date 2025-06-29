package user

import "cosmos-server/api"

type Translator interface {
	ToRegisterUserResponse(username, email, role string) *api.RegisterUserResponse
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
