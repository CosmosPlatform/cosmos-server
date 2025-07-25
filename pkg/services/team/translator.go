package team

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

type Translator interface {
	ToModelTeam(team *obj.Team) *model.Team
	ToModelTeams(teams []*obj.Team) []*model.Team
	ToObjTeam(team *model.Team) *obj.Team
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToModelTeam(team *obj.Team) *model.Team {
	return &model.Team{
		Name:        team.Name,
		Description: team.Description,
	}
}

func (t *translator) ToModelTeams(teams []*obj.Team) []*model.Team {
	var modelTeams []*model.Team
	for _, team := range teams {
		modelTeams = append(modelTeams, t.ToModelTeam(team))
	}
	return modelTeams
}

func (t *translator) ToObjTeam(team *model.Team) *obj.Team {
	return &obj.Team{
		Name:        team.Name,
		Description: team.Description,
	}
}
