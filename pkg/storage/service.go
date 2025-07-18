package storage

import (
	"context"
	"cosmos-server/pkg/storage/obj"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/storage Service

type Service interface {
	GetUserWithEmail(ctx context.Context, email string) (*obj.User, error)
	InsertUser(ctx context.Context, user *obj.User) error
	GetUserWithRole(ctx context.Context, role string) (*obj.User, error)
}
