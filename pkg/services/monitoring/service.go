package monitoring

import (
	"context"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"encoding/json"
	"fmt"
	"strings"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/services/monitoring Service

type Service interface {
	UpdateApplicationInformation(ctx context.Context, application *model.Application) error
	GetApplicationInteractions(ctx context.Context, applicationName string) (*model.ApplicationInteractions, error)
	GetApplicationsInteractions(ctx context.Context, filter model.ApplicationDependencyFilter) (*model.ApplicationsInteractions, error)
}

type monitoringService struct {
	storageService storage.Service
	gitService     GitService
	translator     Translator
	logger         log.Logger
}

func NewMonitoringService(storageService storage.Service, gitService GitService, translator Translator, logger log.Logger) Service {
	return &monitoringService{
		storageService: storageService,
		gitService:     gitService,
		translator:     translator,
		logger:         logger,
	}
}

func (s *monitoringService) UpdateApplicationInformation(ctx context.Context, application *model.Application) error {
	if application.GitInformation == nil {
		s.logger.Infof("No git information for application %s, skipping monitoring update", application.Name)
		return nil // Could be an error because there is nothing to update.
	}

	openClientDef, err := s.getOpenClientDefinition(ctx, application)
	if err != nil {
		return err
	}

	for dependencyName, dependency := range openClientDef.Dependencies {
		modelDependency, err := s.transformToModelDependency(ctx, application, dependencyName, dependency)
		if err != nil {
			s.logger.Errorf("Failed to transform dependency for application %s: %v", application.Name, err)
			continue
		}

		objDependency := s.translator.ToApplicationDependencyObj(modelDependency)

		err = s.storageService.UpsertApplicationDependency(ctx, application.Name, dependencyName, objDependency)
		if err != nil {
			s.logger.Errorf("Failed to upsert dependency for application %s: %v", application.Name, err)
			continue
		}

		s.logger.Infof("Successfully upserted dependency from %s to %s", application.Name, dependencyName)
	}

	err = s.deleteObsoleteDependencies(ctx, application, openClientDef)
	if err != nil {
		s.logger.Errorf("Failed to delete obsolete dependencies for application %s: %v", application.Name, err)
		return err
	}

	return nil
}

func (s *monitoringService) getOpenClientDefinition(ctx context.Context, application *model.Application) (*model.OpenClientSpecification, error) {
	rawOpenClientDefinition, err := s.gitService.GetFileWithContent(ctx, application.GitInformation.RepositoryOwner, application.GitInformation.RepositoryName, application.GitInformation.RepositoryBranch, "docs/openclient.json")
	if err != nil {
		s.logger.Errorf("Failed to get openclient.json for application %s: %v", application.Name, err)
		return nil, err
	}

	var openClientDef model.OpenClientSpecification
	decoder := json.NewDecoder(strings.NewReader(rawOpenClientDefinition.Content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&openClientDef); err != nil {
		s.logger.Errorf("Failed to unmarshal openclient.json for application %s: %v", application.Name, err)
		return nil, err
	}

	if err := openClientDef.Validate(); err != nil {
		s.logger.Errorf("Invalid openclient.json for application %s: %v", application.Name, err)
		return nil, err
	}

	return &openClientDef, nil
}

func (s *monitoringService) transformToModelDependency(ctx context.Context, consumer *model.Application, dependencyName string, dependency model.DependencySpecification) (*model.ApplicationDependency, error) {
	providerApp, err := s.storageService.GetApplicationWithName(ctx, dependencyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider application %s: %v", dependencyName, err)
	}

	providerAppModel := s.translator.ToApplicationModel(providerApp)

	endpoints := make(model.Endpoints)
	for path, methods := range dependency.Endpoints {
		endpointMethods := make(model.EndpointMethods)
		for method, details := range methods {
			endpointMethods[method] = model.EndpointDetails(details)
		}
		endpoints[path] = endpointMethods
	}

	modelDependency := &model.ApplicationDependency{
		Consumer:  consumer,
		Provider:  providerAppModel,
		Reasons:   dependency.Reasons,
		Endpoints: endpoints,
	}

	return modelDependency, nil
}

func (s *monitoringService) GetApplicationInteractions(ctx context.Context, applicationName string) (*model.ApplicationInteractions, error) {
	interactions := make([]*model.ApplicationDependency, 0)
	applicationsInvolved := make(map[string]*model.Application)

	objDependencies, err := s.storageService.GetApplicationDependenciesWithApplicationInvolved(ctx, applicationName)
	if err != nil {
		return nil, err
	}

	for _, objDependency := range objDependencies {
		modelDependency := s.translator.ToApplicationDependencyModel(objDependency)
		interactions = append(interactions, modelDependency)
		if _, exists := applicationsInvolved[modelDependency.Consumer.Name]; !exists {
			applicationsInvolved[modelDependency.Consumer.Name] = modelDependency.Consumer
		}
		if _, exists := applicationsInvolved[modelDependency.Provider.Name]; !exists {
			applicationsInvolved[modelDependency.Provider.Name] = modelDependency.Provider
		}
	}

	return &model.ApplicationInteractions{
		MainApplication:      applicationName,
		ApplicationsInvolved: applicationsInvolved,
		Interactions:         interactions,
	}, nil
}

func (s *monitoringService) deleteObsoleteDependencies(ctx context.Context, application *model.Application, openClientSpecification *model.OpenClientSpecification) error {
	dependencies := make(map[string]bool)

	for dependencyName := range openClientSpecification.Dependencies {
		dependencies[dependencyName] = true
	}

	existingDependencies, err := s.storageService.GetApplicationDependenciesByConsumer(ctx, application.Name)
	if err != nil {
		return err
	}

	for _, existingDependency := range existingDependencies {
		if _, exists := dependencies[existingDependency.Provider.Name]; !exists {
			err := s.storageService.DeleteApplicationDependency(ctx, application.Name, existingDependency.Provider.Name)
			if err != nil {
				s.logger.Errorf("Failed to delete obsolete dependency from %s to %s: %v", application.Name, existingDependency.Provider.Name, err)
				continue
			}
			s.logger.Infof("Successfully deleted obsolete dependency from %s to %s", application.Name, existingDependency.Provider.Name)
		}
	}

	return nil
}

func (s *monitoringService) GetApplicationsInteractions(ctx context.Context, filter model.ApplicationDependencyFilter) (*model.ApplicationsInteractions, error) {
	objDependencies, err := s.storageService.GetApplicationDependenciesWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return s.translator.ToApplicationsInteractionsModel(objDependencies), nil
}
