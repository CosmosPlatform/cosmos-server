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
	GetApplicationInteractions(ctx context.Context, applicationName string) (*model.ApplicationsInteractions, error)
	GetApplicationsInteractions(ctx context.Context, filter model.ApplicationDependencyFilter) (*model.ApplicationsInteractions, error)
	UpdateApplicationOpenAPISpecification(ctx context.Context, application *model.Application) error
}

type monitoringService struct {
	storageService storage.Service
	gitService     GitService
	openApiService OpenApiService
	translator     Translator
	logger         log.Logger
}

func NewMonitoringService(storageService storage.Service, gitService GitService, openApiService OpenApiService, translator Translator, logger log.Logger) Service {
	return &monitoringService{
		storageService: storageService,
		gitService:     gitService,
		openApiService: openApiService,
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
			s.logger.Warnf("Failed to transform dependency for application %s: %v", application.Name, err)
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

func (s *monitoringService) GetApplicationInteractions(ctx context.Context, applicationName string) (*model.ApplicationsInteractions, error) {
	objDependencies, err := s.storageService.GetApplicationDependenciesWithApplicationInvolved(ctx, applicationName)
	if err != nil {
		return nil, err
	}

	return s.translator.ToApplicationsInteractionsModel(objDependencies), nil
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

func (s *monitoringService) UpdateApplicationOpenAPISpecification(ctx context.Context, application *model.Application) error {
	if application.GitInformation == nil {
		s.logger.Infof("No git information for application %s, skipping OpenAPI spec update", application.Name)
		return nil // Could be an error because there is nothing to update.
	}

	openApiSpecMetadata, err := s.gitService.GetFileMetadata(ctx, application.GitInformation.RepositoryOwner, application.GitInformation.RepositoryName, application.GitInformation.RepositoryBranch, "docs/swagger.json")
	if err != nil {
		s.logger.Errorf("Failed to get swagger.json metadata for application %s: %v", application.Name, err)
		return err
	}

	if application.MonitoringInformation != nil && application.MonitoringInformation.OpenAPISha == openApiSpecMetadata.SHA {
		s.logger.Infof("OpenAPI specification for application %s is up to date, skipping update", application.Name)
		return nil
	}

	openAPISpecRaw, err := s.gitService.GetFileWithContent(ctx, application.GitInformation.RepositoryOwner, application.GitInformation.RepositoryName, application.GitInformation.RepositoryBranch, "docs/swagger.json")
	if err != nil {
		s.logger.Errorf("Failed to get swagger.json for application %s: %v", application.Name, err)
		return err
	}

	if openApiSpecMetadata.SHA != openAPISpecRaw.Metadata.SHA {
		return fmt.Errorf("SHA mismatch for swagger.json of application %s", application.Name)
	}

	openApiSpec, err := s.openApiService.ParseOpenApiSpec(openAPISpecRaw.Content)
	if err != nil {
		s.logger.Errorf("Failed to parse OpenAPI spec for application %s: %v", application.Name, err)

	}

	applicationOpenApiObj, err := s.translator.ToApplicationOpenApiObj(openApiSpec)
	if err != nil {
		return fmt.Errorf("failed to transform OpenAPI spec for application %s: %v", application.Name, err)
	}

	err = s.storageService.UpsertOpenAPISpecification(ctx, application.Name, applicationOpenApiObj, openAPISpecRaw.Metadata.SHA)
	if err != nil {
		return fmt.Errorf("failed to upsert OpenAPI spec for application %s: %v", application.Name, err)
	}

	return nil
}
