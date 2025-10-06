package storage

import (
	"context"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/storage Service

type Service interface {
	GetUserWithEmail(ctx context.Context, email string) (*obj.User, error)
	GetUserWithRole(ctx context.Context, role string) (*obj.User, error)
	GetUsersWithFilter(ctx context.Context, filter string) ([]*obj.User, error)
	InsertUser(ctx context.Context, user *obj.User) error
	DeleteUser(ctx context.Context, email string) error

	InsertTeam(ctx context.Context, team *obj.Team) error
	GetTeamsWithFilter(ctx context.Context, filter string) ([]*obj.Team, error)
	DeleteTeam(ctx context.Context, name string) error
	GetTeamWithName(ctx context.Context, name string) (*obj.Team, error)
	AddUserToTeam(ctx context.Context, teamName, username string) error
	RemoveUserFromTeam(ctx context.Context, username string) error

	InsertApplication(ctx context.Context, application *obj.Application) error
	GetApplicationWithName(ctx context.Context, name string) (*obj.Application, error)
	GetApplicationsByTeam(ctx context.Context, team string) ([]*obj.Application, error)
	GetApplicationsWithFilter(ctx context.Context, filter string) ([]*obj.Application, error)
	DeleteApplicationWithName(ctx context.Context, name string) error
	UpdateApplication(ctx context.Context, application *obj.Application) error

	GetApplicationDependency(ctx context.Context, consumerID, providerID int) (*obj.ApplicationDependency, error)
	GetApplicationDependenciesWithApplicationInvolved(ctx context.Context, applicationName string) ([]*obj.ApplicationDependency, error)
	GetApplicationDependenciesWithFilter(ctx context.Context, filters model.ApplicationDependencyFilter) ([]*obj.ApplicationDependency, error)
	GetApplicationDependenciesByConsumer(ctx context.Context, consumerName string) ([]*obj.ApplicationDependency, error)

	UpsertOpenAPISpecification(ctx context.Context, applicationName string, openAPISpec *obj.ApplicationOpenAPI, applicationOpenApiSHA string) error
	UpdateApplicationDependencies(ctx context.Context, applicationName string, dependenciesToUpsert map[string]*obj.ApplicationDependency, pendingDependencies map[string]*obj.PendingApplicationDependency, dependenciesToDelete []*obj.ApplicationDependency, applicationDependenciesSHA string) error
	CheckPendingDependenciesForApplication(ctx context.Context, applicationName string) error
}
