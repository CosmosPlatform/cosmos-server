package auth

import (
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToModelUser(username, email string) *model.User
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToModelUser(username, email string) *model.User {
	return &model.User{
		Username: &username,
		Email:    &email,
	}
}
