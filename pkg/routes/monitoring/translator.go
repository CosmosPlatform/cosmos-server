package monitoring

import (
	"cosmos-server/api"
	"cosmos-server/pkg/model"
)

type Translator interface {
	ToGetApplicationsInteractionsResponse(interactions *model.ApplicationsInteractions) *api.GetApplicationsInteractionsResponse
	ToGetApplicationsInteractionsFilters(teams []string, includeNeighbors bool) model.ApplicationDependencyFilter
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
