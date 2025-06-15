package storage

import (
	"context"
	"cosmos-server/pkg/storage/obj"
)

type Service interface {
	GetUserWithEmail(ctx context.Context, email string) (*obj.User, error)
	InsertUser(ctx context.Context, user *obj.User) (string, error)
}
