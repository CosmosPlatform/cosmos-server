package auth

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

//go:generate mockgen -destination=./mock/translator_mock.go -package=mock cosmos-server/pkg/auth Translator

type Translator interface {
	ToUserModel(userObj *obj.User) *model.User
	ToUserObj(userModel *model.User, encryptedPassword string) *obj.User
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToUserModel(userObj *obj.User) *model.User {
	return &model.User{
		Email:    userObj.Email,
		Username: userObj.Username,
		Role:     userObj.Role,
	}
}

func (t *translator) ToUserObj(userModel *model.User, encryptedPassword string) *obj.User {
	return &obj.User{
		Username:          userModel.Username,
		Email:             userModel.Email,
		EncryptedPassword: encryptedPassword,
		Role:              userModel.Role,
	}
}
