package monitoring

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"cosmos-server/pkg/storage/obj"
	"encoding/json"
	errorUtils "errors"
	"fmt"
	"strings"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/services/monitoring Service

type Service interface {
	UpdateApplicationInformation(ctx context.Context, application *model.Application) error
	GetApplicationInteractions(ctx context.Context, applicationName string) (*model.ApplicationsInteractions, error)
	GetApplicationsInteractions(ctx context.Context, filter model.ApplicationDependencyFilter) (*model.ApplicationsInteractions, error)
	UpdateApplicationOpenAPISpecification(ctx context.Context, application *model.Application) error
	GetApplicationOpenAPISpecification(ctx context.Context, application *model.Application) (*model.ApplicationOpenAPISpecification, error)
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

	if !application.MonitoringInformation.HasOpenClient {
		s.logger.Infof("Application %s does not have OpenClient enabled, skipping monitoring update", application.Name)
		return nil
	}

	openClientMetadata, err := s.gitService.GetFileMetadata(ctx, application.GitInformation.RepositoryOwner, application.GitInformation.RepositoryName, application.GitInformation.RepositoryBranch, application.MonitoringInformation.OpenClientPath)
	if err != nil {
		return fmt.Errorf("failed to get openclient.json metadata for application %s: %v", application.Name, err)
	}

	if application.MonitoringInformation != nil && application.MonitoringInformation.DependenciesSha == openClientMetadata.SHA {
		s.logger.Infof("Dependencies for application %s are up to date, skipping update", application.Name)
		return nil
	}

	rawOpenClientDefinition, err := s.gitService.GetFileWithContent(ctx, application.GitInformation.RepositoryOwner, application.GitInformation.RepositoryName, application.GitInformation.RepositoryBranch, application.MonitoringInformation.OpenClientPath)
	if err != nil {
		return fmt.Errorf("failed to get openclient.json for application %s: %v", application.Name, err)
	}

	if openClientMetadata.SHA != rawOpenClientDefinition.Metadata.SHA {
		return fmt.Errorf("SHA mismatch for openclient.json of application %s", application.Name)
	}

	openClientDef, err := s.transformToOpenClientDefinition(rawOpenClientDefinition)
	if err != nil {
		return fmt.Errorf("failed to transform openclient.json for application %s: %v", application.Name, err)
	}

	dependenciesToUpsert, pendingDependencies, err := s.getDependenciesToModify(ctx, application, openClientDef)
	if err != nil {
		return fmt.Errorf("failed to get dependencies to modify for application %s: %v", application.Name, err)
	}

	// We get the obsolete dependencies to delete them in batch
	objDependenciesToDelete, err := s.getObsoleteDependencies(ctx, application, openClientDef)
	if err != nil {
		return fmt.Errorf("failed to get obsolete dependencies for application %s: %v", application.Name, err)
	}

	err = s.storageService.UpdateApplicationDependencies(ctx, application.Name, dependenciesToUpsert, pendingDependencies, objDependenciesToDelete, rawOpenClientDefinition.Metadata.SHA)
	if err != nil {
		return fmt.Errorf("failed to update dependencies for application %s: %v", application.Name, err)
	}

	return nil
}

func (s *monitoringService) getDependenciesToModify(ctx context.Context, application *model.Application, openClientDef *model.OpenClientSpecification) (map[string]*obj.ApplicationDependency, map[string]*obj.PendingApplicationDependency, error) {
	dependenciesToUpsert := make(map[string]*obj.ApplicationDependency)
	pendingDependencies := make(map[string]*obj.PendingApplicationDependency)

	for dependencyName, dependency := range openClientDef.Dependencies {
		dependencyObj, err := s.storageService.GetApplicationWithName(ctx, dependencyName)
		if err != nil {
			if errorUtils.Is(err, storage.ErrNotFound) {
				s.logger.Warnf("Dependency application %s not found for application %s, skipping dependency creation", dependencyName, application.Name)
				modelPendingDependency := s.transformToModelPendingDependency(application, dependencyName, dependency)
				pendingDependencies[dependencyName] = s.translator.ToPendingApplicationDependencyObj(modelPendingDependency)
				continue
			}
			return nil, nil, fmt.Errorf("failed to get dependency application %s for application %s: %v", dependencyName, application.Name, err)
		}

		modelDependency := s.transformToModelDependency(application, s.translator.ToApplicationModel(dependencyObj), dependency)
		dependenciesToUpsert[dependencyName] = s.translator.ToApplicationDependencyObj(modelDependency)
	}

	return dependenciesToUpsert, pendingDependencies, nil
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

func (s *monitoringService) transformToOpenClientDefinition(rawOpenClientDefinition *model.FileContent) (*model.OpenClientSpecification, error) {
	var openClientDef model.OpenClientSpecification
	decoder := json.NewDecoder(strings.NewReader(rawOpenClientDefinition.Content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&openClientDef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal openclient.json: %s", err.Error())
	}

	if err := openClientDef.Validate(); err != nil {
		return nil, fmt.Errorf("invalid openclient.json :%s", err.Error())
	}

	return &openClientDef, nil
}

func (s *monitoringService) transformToModelDependency(consumer *model.Application, providerAppModel *model.Application, dependency model.DependencySpecification) *model.ApplicationDependency {
	endpoints := s.transformToEndpointsModel(dependency)

	modelDependency := &model.ApplicationDependency{
		Consumer:  consumer,
		Provider:  providerAppModel,
		Reasons:   dependency.Reasons,
		Endpoints: endpoints,
	}

	return modelDependency
}

func (s *monitoringService) transformToModelPendingDependency(consumer *model.Application, dependencyName string, dependency model.DependencySpecification) *model.PendingApplicationDependency {
	endpoints := s.transformToEndpointsModel(dependency)

	modelPendingDependency := &model.PendingApplicationDependency{
		Consumer:     consumer,
		ProviderName: dependencyName,
		Reasons:      dependency.Reasons,
		Endpoints:    endpoints,
	}

	return modelPendingDependency
}

func (s *monitoringService) transformToEndpointsModel(dependency model.DependencySpecification) model.Endpoints {
	endpoints := make(model.Endpoints)
	for path, methods := range dependency.Endpoints {
		endpointMethods := make(model.EndpointMethods)
		for method, details := range methods {
			endpointMethods[method] = model.EndpointDetails(details)
		}
		endpoints[path] = endpointMethods
	}
	return endpoints
}

func (s *monitoringService) GetApplicationInteractions(ctx context.Context, applicationName string) (*model.ApplicationsInteractions, error) {
	objDependencies, err := s.storageService.GetApplicationDependenciesWithApplicationInvolved(ctx, applicationName)
	if err != nil {
		return nil, err
	}

	return s.translator.ToApplicationsInteractionsModel(objDependencies), nil
}

func (s *monitoringService) getObsoleteDependencies(ctx context.Context, application *model.Application, openClientSpecification *model.OpenClientSpecification) ([]*obj.ApplicationDependency, error) {
	dependenciesToKeep := make(map[string]bool)
	dependenciesToDelete := make([]*obj.ApplicationDependency, 0)

	for dependencyName := range openClientSpecification.Dependencies {
		dependenciesToKeep[dependencyName] = true
	}

	existingDependencies, err := s.storageService.GetApplicationDependenciesByConsumer(ctx, application.Name)
	if err != nil {
		return nil, err
	}

	for _, existingDependency := range existingDependencies {
		if _, exists := dependenciesToKeep[existingDependency.Provider.Name]; !exists {
			dependenciesToDelete = append(dependenciesToDelete, existingDependency)
		}
	}

	return dependenciesToDelete, nil
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

	if !application.MonitoringInformation.HasOpenApi {
		s.logger.Infof("Application %s does not have OpenAPI specification enabled, skipping OpenAPI spec update", application.Name)
		return nil
	}

	openApiSpecMetadata, err := s.gitService.GetFileMetadata(ctx, application.GitInformation.RepositoryOwner, application.GitInformation.RepositoryName, application.GitInformation.RepositoryBranch, application.MonitoringInformation.OpenApiPath)
	if err != nil {
		s.logger.Errorf("Failed to get swagger.json metadata for application %s: %v", application.Name, err)
		return err
	}

	if application.MonitoringInformation != nil && application.MonitoringInformation.OpenAPISha == openApiSpecMetadata.SHA {
		s.logger.Infof("OpenAPI specification for application %s is up to date, skipping update", application.Name)
		return nil
	}

	openAPISpecRaw, err := s.gitService.GetFileWithContent(ctx, application.GitInformation.RepositoryOwner, application.GitInformation.RepositoryName, application.GitInformation.RepositoryBranch, application.MonitoringInformation.OpenApiPath)
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

func (s *monitoringService) GetApplicationOpenAPISpecification(ctx context.Context, application *model.Application) (*model.ApplicationOpenAPISpecification, error) {
	if application.GitInformation == nil {
		return nil, errors.NewNotFoundError("The application does not have a git repository associated with it")
	}

	if !application.MonitoringInformation.HasOpenApi {
		return nil, errors.NewNotFoundError("The application does not have OpenAPI specification enabled")
	}

	openApiSpecObj, err := s.storageService.GetOpenAPISpecificationByApplicationName(ctx, application.Name)
	if err != nil {
		return nil, err
	}

	applicationOpenApiModel, err := s.translator.ToApplicationOpenApiModel(openApiSpecObj)
	if err != nil {
		return nil, fmt.Errorf("failed to transform OpenAPI spec for application %s: %v", application.Name, err)
	}

	return applicationOpenApiModel, nil
}
