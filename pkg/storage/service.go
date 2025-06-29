package storage

import (
	"context"
	"cosmos-server/pkg/storage/obj"
)

type Service interface {
	GetUserWithEmail(ctx context.Context, email string) (*obj.User, error)
	InsertUser(ctx context.Context, user *obj.User) error
	GetUserWithRole(ctx context.Context, role string) (*obj.User, error)
}
