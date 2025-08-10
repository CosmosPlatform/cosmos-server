package application

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToCreateApplicationResponse(name, description, team string) *api.CreateApplicationResponse
	ToGetApplicationResponse(applicationObj *model.Application) *api.Application
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
			Team:        &api.Team{Name: team},
		},
	}
}

func (t *translator) ToGetApplicationResponse(applicationModel *model.Application) *api.Application {
	return &api.Application{
		Name:        applicationModel.Name,
		Description: applicationModel.Description,
		Team:        t.ToApiTeam(applicationModel.Team),
	}
}

func (t *translator) ToApiTeam(teamModel *model.Team) *api.Team {
	if teamModel == nil {
		return nil
	}
	return &api.Team{
		Name:        teamModel.Name,
		Description: teamModel.Description,
	}
}
