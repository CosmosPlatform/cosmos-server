package application

import "cosmos-server/api"

type Translator interface {
	ToCreateApplicationResponse(name, description, team string) *api.CreateApplicationResponse
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToCreateApplicationResponse(name, description, team string) *api.CreateApplicationResponse {
	return &api.CreateApplicationResponse{
		Application: &api.Application{
			Name:        name,
			Description: description,
			Team:        team,
		},
	}
}
