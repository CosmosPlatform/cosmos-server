package application

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/storage"
	"cosmos-server/pkg/storage/obj"
	errorUtils "errors"
)

type Service interface {
	AddApplication(ctx context.Context, name, description, team string) error
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
	serviceObj := &obj.Application{
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
		serviceObj.TeamID = &teamIDInt
	}

	err := s.storageService.InsertApplication(ctx, serviceObj)
	if err != nil {
		if errorUtils.Is(err, storage.ErrAlreadyExists) {
			return errors.NewConflictError("application with name " + name + " already exists")
		}
		return errors.NewInternalServerError("failed to insert application: " + err.Error())
	}

	s.logger.Infof("Application %s added successfully", name)
	return nil
}
