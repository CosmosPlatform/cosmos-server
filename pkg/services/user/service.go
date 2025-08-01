package user

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	errorUtils "errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const (
	RegularUserRole = "user"
	AdminUserRole   = "admin"
)

//go:generate mockgen -destination=./mock/service_mock.go -package=mock cosmos-server/pkg/services/user Service

type Service interface {
	GetUserWithEmail(ctx context.Context, email string) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	RegisterUser(ctx context.Context, username, email, password, role string) error
	DeleteUser(ctx context.Context, email string) error
	AdminUserPresent(ctx context.Context) (bool, error)
}

type userService struct {
	storageService storage.Service
	translator     Translator
	logger         log.Logger
}

func NewUserService(storageService storage.Service, translator Translator, logger log.Logger) Service {
	return &userService{
		storageService: storageService,
		translator:     translator,
		logger:         logger,
	}
}

func (s *userService) GetUserWithEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.storageService.GetUserWithEmail(ctx, email)
	if err != nil {
		s.logger.Errorf("failed to get user with email %s: %v", email, err)
		if errorUtils.Is(err, storage.ErrNotFound) {
			return nil, errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
		}
		return nil, errors.NewInternalServerError(fmt.Sprintf("failed to retrieve user with email %s: %v", email, err))
	}

	return s.translator.ToUserModel(user), nil
}

func (s *userService) RegisterUser(ctx context.Context, username, email, password, role string) error {
	if role != RegularUserRole && role != AdminUserRole {
		return errors.NewBadRequestError(fmt.Sprintf("invalid role: %s, must be either '%s' or '%s'", role, RegularUserRole, AdminUserRole))
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Role:     role,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("failed to hash password: %v", err))
	}

	if existingUser, err := s.storageService.GetUserWithEmail(ctx, email); err != nil {
		if !errorUtils.Is(err, storage.ErrNotFound) {
			s.logger.Errorf("failed to check for existing user with email %s: %v", email, err)
			return errors.NewInternalServerError(fmt.Sprintf("failed to check for existing user: %v", err))
		}
	} else if existingUser != nil {
		s.logger.Errorf("user with email %s already exists", email)
		return errors.NewBadRequestError(fmt.Sprintf("user with email %s already exists", email))
	}

	err = s.storageService.InsertUser(ctx, s.translator.ToUserObj(user, string(hashedPassword)))
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("failed to insert user into storage: %v", err))
	}

	return nil
}

func (s *userService) AdminUserPresent(ctx context.Context) (bool, error) {
	adminUser, err := s.storageService.GetUserWithRole(ctx, AdminUserRole)
	if err != nil {
		if errorUtils.Is(err, storage.ErrNotFound) {
			return false, nil
		}
		return false, errors.NewInternalServerError(fmt.Sprintf("failed to check for admin user: %v", err))
	}

	return adminUser != nil, nil
}

func (s *userService) GetUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.storageService.GetUsersWithFilter(ctx, "")
	if err != nil {
		s.logger.Errorf("failed to get users: %v", err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("failed to retrieve users: %v", err))
	}

	userModels := s.translator.ToUserModels(users)

	return userModels, nil
}

func (s *userService) DeleteUser(ctx context.Context, email string) error {
	if email == "" {
		return errors.NewBadRequestError("email cannot be empty")
	}

	err := s.storageService.DeleteUser(ctx, email)
	if err != nil {
		s.logger.Errorf("failed to delete user with email %s: %v", email, err)
		if errorUtils.Is(err, storage.ErrNotFound) {
			return errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("failed to delete user with email %s: %v", email, err))
	}

	return nil
}
