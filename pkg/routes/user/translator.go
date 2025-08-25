package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

//go:generate mockgen -destination=./mock/translator_mock.go -package=mock cosmos-server/pkg/routes/user Translator

type Translator interface {
	ToRegisterUserResponse(username, email, role string) *api.RegisterUserResponse
	ToGetUsersResponse(users []*model.User) *api.GetUsersResponse
	ToGetCurrentUserResponse(userModel *model.User) *api.GetCurrentUserResponse
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

func (t *translator) ToGetCurrentUserResponse(userModel *model.User) *api.GetCurrentUserResponse {
	var apiTeam *api.Team
	if userModel.Team != nil {
		apiTeam = &api.Team{
			Name:        userModel.Team.Name,
			Description: userModel.Team.Description,
		}
	}
	return &api.GetCurrentUserResponse{
		User: api.User{
			Username: userModel.Username,
			Email:    userModel.Email,
			Role:     userModel.Role,
			Team:     apiTeam,
		},
	}
}
