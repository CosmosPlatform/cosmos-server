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
	RegisterRegularUser(ctx context.Context, username, email, password string) error
	RegisterAdminUser(ctx context.Context, username, email, password string) error
	AdminUserPresent(ctx context.Context) (bool, error)
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

func (s *userService) RegisterRegularUser(ctx context.Context, username, email, password string) error {
	user := &model.User{
		Username: username,
		Email:    email,
		Role:     RegularUserRole,
	}

	err := s.registerUser(ctx, user, password)
	if err != nil {
		return fmt.Errorf("failed to register regular user: %w", err)
	}

	return nil
}

func (s *userService) RegisterAdminUser(ctx context.Context, username, email, password string) error {
	user := &model.User{
		Username: username,
		Email:    email,
		Role:     AdminUserRole,
	}

	err := s.registerUser(ctx, user, password)
	if err != nil {
		return fmt.Errorf("failed to register admin user: %w", err)
	}

	return nil
}

func (s *userService) registerUser(ctx context.Context, user *model.User, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.storageService.InsertUser(ctx, s.translator.ToUserObj(user, string(hashedPassword)))
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (s *userService) AdminUserPresent(ctx context.Context) (bool, error) {
	adminUser, err := s.storageService.GetUserWithRole(ctx, AdminUserRole)
	if err != nil {
		return false, fmt.Errorf("failed to check for admin user: %w", err)
	}

	return adminUser != nil, nil
}
