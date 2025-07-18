package user

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/log"
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
	GetUserWithEmail(ctx context.Context, email string) (*model.User, error)
	RegisterRegularUser(ctx context.Context, username, email, password string) error
	RegisterAdminUser(ctx context.Context, username, email, password string) error
	AdminUserPresent(ctx context.Context) (bool, error)
}

type userService struct {
	storageService storage.Service
	translator     Translator
	logger         log.Logger
}

func NewUserService(storageService storage.Service, logger log.Logger) Service {
	return &userService{
		storageService: storageService,
		translator:     NewTranslator(),
		logger:         logger,
	}
}

func (s *userService) GetUserWithEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.storageService.GetUserWithEmail(ctx, email)
	if err != nil {
		s.logger.Errorf("failed to get user with email %s: %v", email, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("failed to retrieve user with email %s: %v", email, err))
	}

	return s.translator.ToUserModel(user), nil
}

func (s *userService) RegisterRegularUser(ctx context.Context, username, email, password string) error {
	user := &model.User{
		Username: username,
		Email:    email,
		Role:     RegularUserRole,
	}

	err := s.registerUser(ctx, user, password)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (s *userService) registerUser(ctx context.Context, user *model.User, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("failed to hash password: %v", err))
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
		return false, errors.NewInternalServerError(fmt.Sprintf("failed to check for admin user: %v", err))
	}

	return adminUser != nil, nil
}
