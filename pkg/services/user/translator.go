package user

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

//go:generate mockgen -destination=./mock/translator_mock.go -package=mock cosmos-server/pkg/services/user Translator

type Translator interface {
	ToUserObj(userModel *model.User, encryptedPassword string) *obj.User
	ToUserModel(userObj *obj.User) *model.User
	ToUserModels(userObjs []*obj.User) []*model.User
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

func (t *translator) ToUserModel(userObj *obj.User) *model.User {
	return &model.User{
		Username: userObj.Username,
		Email:    userObj.Email,
		Role:     userObj.Role,
	}
}

func (t *translator) ToUserModels(userObjs []*obj.User) []*model.User {
	var userModels []*model.User
	for _, userObj := range userObjs {
		userModels = append(userModels, t.ToUserModel(userObj))
	}
	return userModels
}
