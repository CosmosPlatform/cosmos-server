package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToRegisterUserResponse(username, email, role string) *api.RegisterUserResponse
	ToGetUsersResponse(users []*model.User) *api.GetUsersResponse
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToRegisterUserResponse(username, email, role string) *api.RegisterUserResponse {
	return &api.RegisterUserResponse{
		User: api.User{
			Username: username,
			Email:    email,
			Role:     role,
		},
	}
}

func (t *translator) ToGetUsersResponse(users []*model.User) *api.GetUsersResponse {
	apiUsers := make([]*api.User, 0)
	for _, user := range users {
		var apiTeam *api.Team
		if user.Team != nil {
			apiTeam = &api.Team{
				Name:        user.Team.Name,
				Description: user.Team.Description,
			}
		}
		apiUsers = append(apiUsers, &api.User{
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Team:     apiTeam,
		})
	}
	return &api.GetUsersResponse{
		Users: apiUsers,
	}
}
