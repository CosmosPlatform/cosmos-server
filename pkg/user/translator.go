package user

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

type Translator interface {
	ToUserObj(userModel *model.User, encryptedPassword string) *obj.User
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToUserObj(userModel *model.User, encryptedPassword string) *obj.User {
	return &obj.User{
		Username:          userModel.Username,
		Email:             userModel.Email,
		EncryptedPassword: encryptedPassword,
		Role:              userModel.Role,
	}
}
