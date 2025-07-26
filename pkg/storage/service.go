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
}
