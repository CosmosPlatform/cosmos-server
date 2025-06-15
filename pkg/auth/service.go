package auth

import (
	"cosmos-server/pkg/model"
)

type Service interface {
	// Authenticate authenticates a user with the provided credentials and returns the user and a token if successful.
	Authenticate(email, password string) (*model.User, string, error)
	// IsAuthenticated checks if the provided token is valid for the user and returns the user ID if authenticated.
	IsAuthenticated(token string) (string, error)
}
