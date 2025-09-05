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

type Service interface {
	UpdateApplicationInformation(ctx context.Context, application *model.Application) error
	GetApplicationInteractions(ctx context.Context, applicationName string) (*model.ApplicationInteractions, error)
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
	applicationsToProvide := make([]*model.Application, 0)
	applicationsToConsume := make([]*model.Application, 0)

	objDependenciesAsConsumer, err := s.storageService.GetApplicationDependenciesByConsumer(ctx, applicationName)
	if err != nil {
		return nil, err
	}

	for _, objDependency := range objDependenciesAsConsumer {
		modelDependency := s.translator.ToApplicationDependencyModel(objDependency)
		applicationsToConsume = append(applicationsToConsume, modelDependency.Consumer)
		interactions = append(interactions, modelDependency)
	}

	objDependenciesAsProvider, err := s.storageService.GetApplicationDependenciesByProvider(ctx, applicationName)
	if err != nil {
		return nil, err
	}

	for _, objDependency := range objDependenciesAsProvider {
		modelDependency := s.translator.ToApplicationDependencyModel(objDependency)
		applicationsToProvide = append(applicationsToProvide, modelDependency.Provider)
		interactions = append(interactions, modelDependency)
	}

	return &model.ApplicationInteractions{
		MainApplication:       nil, // We fill this in the calling function
		ApplicationsToProvide: applicationsToProvide,
		ApplicationsToConsume: applicationsToConsume,
		Interactions:          interactions,
	}, nil
}
