package team

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToGetTeamsResponse(teams []*model.Team) *api.GetTeamsResponse
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
