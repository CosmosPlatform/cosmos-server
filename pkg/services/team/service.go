package team

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	errorUtils "errors"
	"fmt"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/services/team Service

type Service interface {
	GetAllTeams(ctx context.Context) ([]*model.Team, error)
	DeleteTeam(ctx context.Context, name string) error
	InsertTeam(ctx context.Context, team *model.Team) error
}

type teamService struct {
	storageService storage.Service
	translator     Translator
}

func NewTeamService(storageService storage.Service, translator Translator) Service {
	return &teamService{
		storageService: storageService,
		translator:     translator,
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

func (s *teamService) InsertTeam(ctx context.Context, team *model.Team) error {
	objTeam := s.translator.ToObjTeam(team)
	err := s.storageService.InsertTeam(ctx, objTeam)
	if err != nil {
		if errorUtils.Is(err, storage.ErrAlreadyExists) {
			return errors.NewConflictError(fmt.Sprint("team with name ", team.Name, " already exists"))
		}
		return err
	}
	return nil
}

func (s *teamService) DeleteTeam(ctx context.Context, name string) error {
	err := s.storageService.DeleteTeam(ctx, name)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError(fmt.Sprintf("team with name %s not found", name))
		}
		return errors.NewInternalServerError(fmt.Sprintf("failed to delete team with name %s: %v", name, err))
	}
	return nil
}
