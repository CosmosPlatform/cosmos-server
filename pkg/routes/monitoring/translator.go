package monitoring

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToGetApplicationsInteractionsResponse(interactions *model.ApplicationsInteractions) *api.GetApplicationsInteractionsResponse
	ToGetApplicationsInteractionsFilters(teams []string, includeNeighbors bool) model.ApplicationDependencyFilter
	ToGetOpenAPiSpecificationResponse(openAPISpec *model.ApplicationOpenAPISpecification) (*api.GetApplicationOpenAPISpecificationResponse, error)
	ToGetCompleteApplicationMonitoringResponse(application *model.Application, interactions *model.ApplicationsInteractions, openAPISpec *model.ApplicationOpenAPISpecification) (*api.GetCompleteApplicationMonitoringResponse, error)
	ToSentinelSettingsUpdateModel(updateSettingsApi *api.UpdateSentinelSettingsRequest) *model.SentinelSettingsUpdate
	ToGetSentinelSettingsResponse(sentinelSettingsModel *model.SentinelSettings) *api.GetSentinelSettingsResponse
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToGetApplicationsInteractionsResponse(interactions *model.ApplicationsInteractions) *api.GetApplicationsInteractionsResponse {
	if interactions == nil {
		return nil
	}

	return &api.GetApplicationsInteractionsResponse{
		ApplicationsInvolved: t.toApplicationInformationSlice(interactions.ApplicationsInvolved),
		Dependencies:         t.toApplicationDependencySlice(interactions.Interactions),
	}
}

func (t *translator) toApplicationInformation(app *model.Application) api.ApplicationInformation {
	if app == nil {
		return api.ApplicationInformation{}
	}

	teamName := ""
	if app.Team != nil {
		teamName = app.Team.Name
	}
	return api.ApplicationInformation{
		Name: app.Name,
		Team: teamName,
	}
}

func (t *translator) toApplicationInformationSlice(apps map[string]*model.Application) map[string]api.ApplicationInformation {
	if apps == nil {
		return nil
	}

	result := make(map[string]api.ApplicationInformation, len(apps))
	for key, app := range apps {
		result[key] = t.toApplicationInformation(app)
	}
	return result
}

func (t *translator) toApplicationDependencySlice(deps []*model.ApplicationDependency) []api.ApplicationDependency {
	if deps == nil {
		return nil
	}

	result := make([]api.ApplicationDependency, 0, len(deps))
	for _, dep := range deps {
		result = append(result, t.toApplicationDependency(dep))
	}
	return result
}

func (t *translator) toApplicationDependency(dep *model.ApplicationDependency) api.ApplicationDependency {
	if dep == nil {
		return api.ApplicationDependency{}
	}

	return api.ApplicationDependency{
		Consumer:  dep.Consumer.Name,
		Provider:  dep.Provider.Name,
		Reasons:   dep.Reasons,
		Endpoints: t.toDependencyEndpointsMap(dep.Endpoints),
	}
}

func (t *translator) toDependencyEndpointsMap(endpoints model.Endpoints) api.Endpoints {
	if endpoints == nil {
		return nil
	}

	result := make(api.Endpoints)
	for path, methods := range endpoints {
		endpointMethods := make(api.EndpointMethods)
		for method, details := range methods {
			endpointMethods[method] = api.EndpointDetails(details)
		}
		result[path] = endpointMethods
	}
	return result
}

func (t *translator) ToGetApplicationsInteractionsFilters(teams []string, includeNeighbors bool) model.ApplicationDependencyFilter {
	var modelFilters model.ApplicationDependencyFilter

	if teams != nil {
		modelFilters.Teams = teams
	}

	modelFilters.IncludeNeighbors = includeNeighbors

	return modelFilters
}

func (t *translator) ToGetOpenAPiSpecificationResponse(modelOpenAPISpec *model.ApplicationOpenAPISpecification) (*api.GetApplicationOpenAPISpecificationResponse, error) {
	if modelOpenAPISpec == nil {
		return nil, nil
	}

	marshalledOpenAPISpec, err := t.toMarshalledOpenAPISpec(modelOpenAPISpec)
	if err != nil {
		return nil, err
	}

	applicationName := ""
	if modelOpenAPISpec.Application != nil {
		applicationName = modelOpenAPISpec.Application.Name
	}

	return &api.GetApplicationOpenAPISpecificationResponse{
		ApplicationName: applicationName,
		OpenAPISpec:     string(marshalledOpenAPISpec),
	}, nil
}

func (t *translator) toMarshalledOpenAPISpec(openAPISpec *model.ApplicationOpenAPISpecification) (string, error) {
	if openAPISpec == nil {
		return "", nil
	}

	marshalledOpenAPISpec, err := openAPISpec.OpenAPISpec.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(marshalledOpenAPISpec), nil
}

func (t *translator) ToGetCompleteApplicationMonitoringResponse(application *model.Application, interactions *model.ApplicationsInteractions, openAPISpec *model.ApplicationOpenAPISpecification) (*api.GetCompleteApplicationMonitoringResponse, error) {
	if application == nil {
		return nil, nil
	}

	applicationApi := t.ToApplicationApi(application)

	marshalledOpenAPISpec, err := t.toMarshalledOpenAPISpec(openAPISpec)
	if err != nil {
		return nil, err
	}

	dependencies, consumedEndpoints := t.toDependenciesAndConsumedEndpoints(interactions, application.Name)

	return &api.GetCompleteApplicationMonitoringResponse{
		Application:       *applicationApi,
		OpenAPISpec:       marshalledOpenAPISpec,
		Dependencies:      dependencies,
		ConsumedEndpoints: consumedEndpoints,
	}, nil
}

func (t *translator) toDependenciesAndConsumedEndpoints(interactions *model.ApplicationsInteractions, applicationName string) ([]api.ApplicationDependency, api.ConsumedEndpoints) {
	if interactions == nil {
		return nil, nil
	}

	dependencies := make([]api.ApplicationDependency, 0)
	consumedEndpoints := make(api.ConsumedEndpoints)

	for _, interaction := range interactions.Interactions {
		if interaction.Consumer.Name == applicationName {
			dependencies = append(dependencies, t.toApplicationDependency(interaction))
		} else if interaction.Provider.Name == applicationName {
			for path, methods := range interaction.Endpoints {
				if _, exists := consumedEndpoints[path]; !exists {
					consumedEndpoints[path] = make(api.ConsumedEndpointMethods)
				}
				for method := range methods {
					if _, methodExists := consumedEndpoints[path][method]; !methodExists {
						consumedEndpoints[path][method] = api.ConsumedEndpointDetails{
							Consumers: []string{},
						}
					}
					currentConsumedEndpoints := consumedEndpoints[path][method]
					currentConsumedEndpoints.Consumers = append(currentConsumedEndpoints.Consumers, interaction.Consumer.Name)
					consumedEndpoints[path][method] = currentConsumedEndpoints
				}
			}
		}
	}

	return dependencies, consumedEndpoints
}

func (t *translator) ToApplicationApi(applicationModel *model.Application) *api.Application {
	return &api.Application{
		Name:                  applicationModel.Name,
		Description:           applicationModel.Description,
		Team:                  t.ToApiTeam(applicationModel.Team),
		GitInformation:        t.ToApiGitInformation(applicationModel.GitInformation),
		MonitoringInformation: t.ToMonitoringInformationApi(applicationModel.MonitoringInformation),
		Token:                 t.ToTokenApi(applicationModel.Token),
	}
}

func (t *translator) ToTokenApi(tokenModel *model.Token) *api.Token {
	if tokenModel == nil {
		return nil
	}

	team := ""
	if tokenModel.Team != nil {
		team = tokenModel.Team.Name
	}

	return &api.Token{
		Name: tokenModel.Name,
		Team: team,
	}
}

func (t *translator) ToApiTeam(teamModel *model.Team) *api.Team {
	if teamModel == nil {
		return nil
	}
	return &api.Team{
		Name:        teamModel.Name,
		Description: teamModel.Description,
	}
}

func (t *translator) ToApiGitInformation(gitInfo *model.GitInformation) *api.GitInformation {
	if gitInfo == nil {
		return nil
	}
	return &api.GitInformation{
		Provider:         gitInfo.Provider,
		RepositoryOwner:  gitInfo.RepositoryOwner,
		RepositoryName:   gitInfo.RepositoryName,
		RepositoryBranch: gitInfo.RepositoryBranch,
	}
}

func (t *translator) ToMonitoringInformationApi(monitoringInfo *model.MonitoringInformation) *api.MonitoringInformation {
	if monitoringInfo == nil {
		return nil
	}

	return &api.MonitoringInformation{
		HasOpenAPI:     monitoringInfo.HasOpenApi,
		OpenAPIPath:    monitoringInfo.OpenApiPath,
		HasOpenClient:  monitoringInfo.HasOpenClient,
		OpenClientPath: monitoringInfo.OpenClientPath,
	}
}

func (t *translator) ToSentinelSettingsUpdateModel(updateSettingsApi *api.UpdateSentinelSettingsRequest) *model.SentinelSettingsUpdate {
	if updateSettingsApi == nil {
		return nil
	}

	return &model.SentinelSettingsUpdate{
		Enabled:  updateSettingsApi.Enabled,
		Interval: updateSettingsApi.Interval,
	}
}

func (t *translator) ToGetSentinelSettingsResponse(sentinelSettingsModel *model.SentinelSettings) *api.GetSentinelSettingsResponse {
	if sentinelSettingsModel == nil {
		return nil
	}

	return &api.GetSentinelSettingsResponse{
		Enabled:  sentinelSettingsModel.Enabled,
		Interval: sentinelSettingsModel.Interval,
	}
}
