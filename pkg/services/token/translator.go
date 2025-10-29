package token

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

type Translator interface {
	ToModelToken(obj *obj.Token) *model.Token
	ToModelTokens(objs []*obj.Token) []*model.Token
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}
func (t *translator) ToModelToken(obj *obj.Token) *model.Token {
	if obj == nil {
		return nil
	}

	return &model.Token{
		Name:           obj.Name,
		EncryptedValue: obj.EncryptedValue,
		Team:           t.ToModelTeam(obj.Team),
	}
}

func (t *translator) ToModelTokens(objs []*obj.Token) []*model.Token {
	models := make([]*model.Token, 0, len(objs))
	for _, obj := range objs {
		models = append(models, t.ToModelToken(obj))
	}
	return models
}

func (t *translator) ToModelTeam(obj *obj.Team) *model.Team {
	if obj == nil {
		return nil
	}

	return &model.Team{
		Name:        obj.Name,
		Description: obj.Description,
	}
}
