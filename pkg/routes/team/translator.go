package team

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToGetTeamsResponse(teams []*model.Team) *api.GetTeamsResponse
	ToInsertTeamResponse(name, description string) *api.InsertTeamResponse
	ToModelTeam(name, description string) *model.Team
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToGetTeamsResponse(teams []*model.Team) *api.GetTeamsResponse {
	var apiTeams []*api.Team
	for _, team := range teams {
		apiTeams = append(apiTeams, &api.Team{
			Name:        team.Name,
			Description: team.Description,
		})
	}
	return &api.GetTeamsResponse{
		Teams: apiTeams,
	}
}

func (t *translator) ToInsertTeamResponse(name, description string) *api.InsertTeamResponse {
	return &api.InsertTeamResponse{
		Team: &api.Team{
			Name:        name,
			Description: description,
		},
	}
}

func (t *translator) ToModelTeam(name, description string) *model.Team {
	return &model.Team{
		Name:        name,
		Description: description,
	}
}
