package storage

import (
	"context"
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

	InsertApplicationDependency(ctx context.Context, dependency *obj.ApplicationDependency) error
	GetApplicationDependency(ctx context.Context, consumerID, providerID int) (*obj.ApplicationDependency, error)
	UpsertApplicationDependency(ctx context.Context, consumerName, providerName string, dependency *obj.ApplicationDependency) error
}
