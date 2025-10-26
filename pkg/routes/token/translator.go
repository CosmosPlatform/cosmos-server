package token

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToGetTokenResponse(modelTokens []*model.Token) *api.GetTokensResponse
	ToApiToken(modelToken *model.Token) *api.Token
	ToApiTokens(modelTokens []*model.Token) []*api.Token
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToGetTokenResponse(modelTokens []*model.Token) *api.GetTokensResponse {
	return &api.GetTokensResponse{
		Tokens: t.ToApiTokens(modelTokens),
	}
}

func (t *translator) ToApiToken(modelToken *model.Token) *api.Token {
	if modelToken == nil {
		return nil
	}

	teamName := ""
	if modelToken.Team != nil {
		teamName = modelToken.Team.Name
	}

	return &api.Token{
		Name: modelToken.Name,
		Team: teamName,
	}
}

func (t *translator) ToApiTokens(modelTokens []*model.Token) []*api.Token {
	apiTokens := make([]*api.Token, 0, len(modelTokens))
	for _, modelToken := range modelTokens {
		apiTokens = append(apiTokens, t.ToApiToken(modelToken))
	}

	return apiTokens
}
