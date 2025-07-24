package team

import (
	"context"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
)

type Service interface {
	GetAllTeams(ctx context.Context) ([]*model.Team, error)
}

type teamService struct {
	storageService storage.Service
	translator     Translator
}

func NewTeamService(storageService storage.Service) Service {
	return &teamService{
		storageService: storageService,
	}
}

func (s *teamService) GetAllTeams(ctx context.Context) ([]*model.Team, error) {
	teams, err := s.storageService.GetTeamsWithFilter(ctx, "")
	if err != nil {
		return nil, err
	}

	teamModels := s.translator.ToModelTeams(teams)

	return teamModels, nil
}
