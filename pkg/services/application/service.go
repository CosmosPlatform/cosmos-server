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
	AddApplication(ctx context.Context, name, description, team string, gitInformation *model.GitInformation, monitoringInformation *model.MonitoringInformation) error
	GetApplication(ctx context.Context, name string) (*model.Application, error)
	GetApplicationsByTeam(ctx context.Context, team string) ([]*model.Application, error)
	GetApplicationsWithFilter(ctx context.Context, filter string) ([]*model.Application, error)
	DeleteApplication(ctx context.Context, name string) error
	UpdateApplication(ctx context.Context, name string, updateData *model.ApplicationUpdate) (*model.Application, error)

	GetApplicationsToMonitor(ctx context.Context) ([]*model.Application, error)
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

func (s *applicationService) AddApplication(ctx context.Context, name, description, team string, gitInformation *model.GitInformation, monitoringInformation *model.MonitoringInformation) error {
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

	if gitInformation != nil {
		applicationObj.GitProvider = gitInformation.Provider
		applicationObj.GitRepositoryOwner = gitInformation.RepositoryOwner
		applicationObj.GitRepositoryName = gitInformation.RepositoryName
		applicationObj.GitRepositoryBranch = gitInformation.RepositoryBranch
	}

	if monitoringInformation != nil {
		applicationObj.HasOpenClient = monitoringInformation.HasOpenClient
		applicationObj.HasOpenApi = monitoringInformation.HasOpenApi
		applicationObj.OpenApiPath = monitoringInformation.OpenApiPath
		applicationObj.OpenClientPath = monitoringInformation.OpenClientPath
	}

	err := s.storageService.InsertApplication(ctx, applicationObj)
	if err != nil {
		if errorUtils.Is(err, storage.ErrAlreadyExists) {
			return errors.NewConflictError("application with name " + name + " already exists")
		}
		return errors.NewInternalServerError("failed to insert application: " + err.Error())
	}

	err = s.storageService.CheckPendingDependenciesForApplication(ctx, name)
	if err != nil {
		s.logger.Errorf("Failed to check pending dependencies for application %s: %v", name, err)
		// Not returning error to avoid failing the whole operation
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

func (s *applicationService) GetApplicationsByTeam(ctx context.Context, team string) ([]*model.Application, error) {
	applications, err := s.storageService.GetApplicationsByTeam(ctx, team)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return nil, errors.NewNotFoundError("team not found")
		}
		return nil, errors.NewInternalServerError("failed to retrieve applications by team: " + err.Error())
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

func (s *applicationService) UpdateApplication(ctx context.Context, name string, updateData *model.ApplicationUpdate) (*model.Application, error) {
	existingApp, err := s.storageService.GetApplicationWithName(ctx, name)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return nil, errors.NewNotFoundError("application not found")
		}
		return nil, errors.NewInternalServerError("failed to retrieve application: " + err.Error())
	}

	updateObj := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID:        existingApp.ID,
			CreatedAt: existingApp.CreatedAt,
		},
		Name:                existingApp.Name,
		Description:         existingApp.Description,
		TeamID:              existingApp.TeamID,
		GitProvider:         existingApp.GitProvider,
		GitRepositoryOwner:  existingApp.GitRepositoryOwner,
		GitRepositoryName:   existingApp.GitRepositoryName,
		GitRepositoryBranch: existingApp.GitRepositoryBranch,
		DependenciesSha:     existingApp.DependenciesSha,
		OpenAPISha:          existingApp.OpenAPISha,
		HasOpenApi:          existingApp.HasOpenApi,
		OpenApiPath:         existingApp.OpenApiPath,
		HasOpenClient:       existingApp.HasOpenClient,
		OpenClientPath:      existingApp.OpenClientPath,
	}

	if updateData.Name != nil {
		updateObj.Name = *updateData.Name
	}
	if updateData.Description != nil {
		updateObj.Description = *updateData.Description
	}

	if updateData.Team != nil {
		if *updateData.Team == "" {
			updateObj.TeamID = nil
		} else {
			teamObj, err := s.storageService.GetTeamWithName(ctx, *updateData.Team)
			if err != nil {
				if errorUtils.Is(err, storage.ErrNotFound) {
					return nil, errors.NewNotFoundError("team not found")
				}
				return nil, errors.NewInternalServerError("failed to retrieve team: " + err.Error())
			}
			teamIDInt := int(teamObj.ID)
			updateObj.TeamID = &teamIDInt
		}
	}

	if updateData.GitInformation != nil {
		updateObj.GitProvider = updateData.GitInformation.Provider
		updateObj.GitRepositoryOwner = updateData.GitInformation.RepositoryOwner
		updateObj.GitRepositoryName = updateData.GitInformation.RepositoryName
		updateObj.GitRepositoryBranch = updateData.GitInformation.RepositoryBranch
		if updateData.MonitoringInformation != nil {
			updateObj.HasOpenApi = updateData.MonitoringInformation.HasOpenApi
			updateObj.OpenApiPath = updateData.MonitoringInformation.OpenApiPath
			updateObj.HasOpenClient = updateData.MonitoringInformation.HasOpenClient
			updateObj.OpenClientPath = updateData.MonitoringInformation.OpenClientPath
		}
	}

	err = s.storageService.UpdateApplication(ctx, updateObj)
	if err != nil {
		if errorUtils.Is(err, storage.ErrAlreadyExists) {
			if updateData.Name != nil {
				return nil, errors.NewConflictError("another application with name " + updateObj.Name + " already exists")
			}
		}
		return nil, errors.NewInternalServerError("failed to update application: " + err.Error())
	}

	updatedApp, err := s.storageService.GetApplicationWithName(ctx, updateObj.Name)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to retrieve updated application: " + err.Error())
	}

	s.logger.Infof("Application %s updated successfully", updateObj.Name)
	return s.translator.ToApplicationModel(updatedApp), nil
}

func (s *applicationService) GetApplicationsToMonitor(ctx context.Context) ([]*model.Application, error) {
	applications, err := s.storageService.GetApplicationsToMonitor(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError("failed to retrieve applications to monitor: " + err.Error())
	}

	return s.translator.ToApplicationModels(applications), nil
}
