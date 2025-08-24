package application

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"cosmos-server/pkg/storage/obj"
	errorUtils "errors"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/services/application Service

type Service interface {
	AddApplication(ctx context.Context, name, description, team string) error
	GetApplication(ctx context.Context, name string) (*model.Application, error)
	GetApplicationsWithFilter(ctx context.Context, filter string) ([]*model.Application, error)
	DeleteApplication(ctx context.Context, name string) error
}

type applicationService struct {
	storageService storage.Service
	translator     Translator
	logger         log.Logger
}

func NewApplicationService(storageService storage.Service, translator Translator, logger log.Logger) Service {
	return &applicationService{
		storageService: storageService,
		translator:     translator,
		logger:         logger,
	}
}

func (s *applicationService) AddApplication(ctx context.Context, name, description, team string) error {
	applicationObj := &obj.Application{
		Name:        name,
		Description: description,
	}

	if team != "" {
		teamObj, err := s.storageService.GetTeamWithName(ctx, team)
		if err != nil {
			if errorUtils.Is(err, storage.ErrNotFound) {
				return errors.NewNotFoundError("team not found")
			}
			return errors.NewInternalServerError("failed to retrieve team: " + err.Error())
		}
		teamIDInt := int(teamObj.ID)
		applicationObj.TeamID = &teamIDInt
	}

	err := s.storageService.InsertApplication(ctx, applicationObj)
	if err != nil {
		if errorUtils.Is(err, storage.ErrAlreadyExists) {
			return errors.NewConflictError("application with name " + name + " already exists")
		}
		return errors.NewInternalServerError("failed to insert application: " + err.Error())
	}

	s.logger.Infof("Application %s added successfully", name)
	return nil
}

func (s *applicationService) GetApplication(ctx context.Context, name string) (*model.Application, error) {
	application, err := s.storageService.GetApplicationWithName(ctx, name)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return nil, errors.NewNotFoundError("application not found")
		}
		return nil, errors.NewInternalServerError("failed to retrieve application: " + err.Error())
	}

	return s.translator.ToApplicationModel(application), nil
}

func (s *applicationService) GetApplicationsWithFilter(ctx context.Context, filter string) ([]*model.Application, error) {
	applications, err := s.storageService.GetApplicationsWithFilter(ctx, filter)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to retrieve applications: " + err.Error())
	}

	return s.translator.ToApplicationModels(applications), nil
}

func (s *applicationService) DeleteApplication(ctx context.Context, name string) error {
	err := s.storageService.DeleteApplicationWithName(ctx, name)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError("application not found")
		}
		return errors.NewInternalServerError("failed to delete application: " + err.Error())
	}

	s.logger.Infof("Application %s deleted successfully", name)
	return nil
}
