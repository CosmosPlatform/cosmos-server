package auth

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

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
		ID:    userObj.ID,
		Email: userObj.Email,
	}
}

func (t *translator) ToUserObj(userModel *model.User, encryptedPassword string) *obj.User {
	return &obj.User{
		ID:                userModel.ID,
		Username:          userModel.Username,
		Email:             userModel.Email,
		EncryptedPassword: encryptedPassword,
		Role:              userModel.Role,
	}
}
