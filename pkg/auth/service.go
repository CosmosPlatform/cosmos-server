package auth

import (
	"context"
	"cosmos-server/pkg/model"
	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	// Authenticate authenticates a user with the provided credentials and returns the user and a token if successful.
	Authenticate(ctx context.Context, email, password string) (*model.User, string, error)
	// IsAuthenticated checks if the provided token is valid for the user and returns the user ID if authenticated.
	IsAuthenticated(tokenString string) (*jwt.Token, error)
}
