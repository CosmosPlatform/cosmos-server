package user

import (
	"context"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	RegularUserRole = "user"
	AdminUserRole   = "admin"
)

type Service interface {
	RegisterRegularUser(ctx context.Context, username, email, password string) (string, error)
	RegisterAdminUser(ctx context.Context, username, email, password string) (string, error)
}

type userService struct {
	storageService storage.Service
	translator     Translator
}

func NewUserService(storageService storage.Service) Service {
	return &userService{
		storageService: storageService,
		translator:     NewTranslator(),
	}
}

func (s *userService) RegisterRegularUser(ctx context.Context, username, email, password string) (string, error) {
	user := &model.User{
		Username: username,
		Email:    email,
		Role:     RegularUserRole,
	}

	userID, err := s.registerUser(ctx, user, password)
	if err != nil {
		return "", fmt.Errorf("failed to register regular user: %w", err)
	}

	return userID, nil
}

func (s *userService) RegisterAdminUser(ctx context.Context, username, email, password string) (string, error) {
	user := &model.User{
		Username: username,
		Email:    email,
		Role:     AdminUserRole,
	}

	userID, err := s.registerUser(ctx, user, password)
	if err != nil {
		return "", fmt.Errorf("failed to register admin user: %w", err)
	}

	return userID, nil
}

func (s *userService) registerUser(ctx context.Context, user *model.User, password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := s.storageService.InsertUser(ctx, s.translator.ToUserObj(user, string(hashedPassword)))
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, nil
}
