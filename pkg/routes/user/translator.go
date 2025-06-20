package user

import "cosmos-server/api"

type Translator interface {
	ToRegisterUserResponse(userID, username, email string) *api.RegisterUserResponse
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToRegisterUserResponse(userID, username, email string) *api.RegisterUserResponse {
	return &api.RegisterUserResponse{
		User: api.User{
			ID:       userID,
			Username: username,
			Email:    email,
		},
	}
}
