package monitoring

import (
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

type Translator interface {
	ToApplicationModel(objApplication *obj.Application) *model.Application
	ToApplicationDependencyObj(modelDependency *model.ApplicationDependency) *obj.ApplicationDependency
	ToApplicationDependencyModel(objDependency *obj.ApplicationDependency) *model.ApplicationDependency
}

type translator struct{}

func NewTranslator() Translator {
	return &translator{}
}

func (t *translator) ToApplicationModel(applicationObj *obj.Application) *model.Application {
	modelApplication := &model.Application{
		Name:        applicationObj.Name,
		Description: applicationObj.Description,
		Team:        t.ToModelTeam(applicationObj.Team),
	}

	if applicationObj.GitProvider != "" || applicationObj.GitRepositoryName != "" || applicationObj.GitRepositoryOwner != "" || applicationObj.GitRepositoryBranch != "" {
		modelApplication.GitInformation = &model.GitInformation{
			Provider:         applicationObj.GitProvider,
			RepositoryOwner:  applicationObj.GitRepositoryOwner,
			RepositoryName:   applicationObj.GitRepositoryName,
			RepositoryBranch: applicationObj.GitRepositoryBranch,
		}
	}

	return modelApplication
}

func (t *translator) ToModelTeam(teamObj *obj.Team) *model.Team {
	if teamObj == nil {
		return nil
	}
	return &model.Team{
		Name:        teamObj.Name,
		Description: teamObj.Description,
	}
}

func (t *translator) ToApplicationDependencyObj(modelDependency *model.ApplicationDependency) *obj.ApplicationDependency {
	if modelDependency == nil {
		return nil
	}

	return &obj.ApplicationDependency{
		Reasons:   modelDependency.Reasons,
		Endpoints: t.toObjEndpoints(modelDependency.Endpoints),
	}
}

func (t *translator) toObjEndpoints(modelEndpoints model.Endpoints) obj.Endpoints {
	endpoints := make(obj.Endpoints)

	for path, methods := range modelEndpoints {
		if _, exists := endpoints[path]; !exists {
			endpoints[path] = make(obj.EndpointMethods)
		}

		for method, details := range methods {
			endpoints[path][method] = obj.EndpointDetails{
				Reasons: details.Reasons,
			}
		}
	}

	return endpoints
}

func (t *translator) ToApplicationDependencyModel(objDependency *obj.ApplicationDependency) *model.ApplicationDependency {
	if objDependency == nil {
		return nil
	}

	return &model.ApplicationDependency{
		Consumer:  t.ToApplicationModel(objDependency.Consumer),
		Provider:  t.ToApplicationModel(objDependency.Provider),
		Reasons:   objDependency.Reasons,
		Endpoints: t.toModelEndpoints(objDependency.Endpoints),
	}
}

func (t *translator) toModelEndpoints(objEndpoints obj.Endpoints) model.Endpoints {
	endpoints := make(model.Endpoints)

	for path, methods := range objEndpoints {
		if _, exists := endpoints[path]; !exists {
			endpoints[path] = make(model.EndpointMethods)
		}

		for method, details := range methods {
			endpoints[path][method] = model.EndpointDetails(details)
		}
	}

	return endpoints
}
