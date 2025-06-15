package auth

import (
	"context"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type StorageAuthService struct {
	storageService storage.Service
	translator     Translator
}

func NewStorageAuthService(storageService storage.Service) *StorageAuthService {
	return &StorageAuthService{
		storageService: storageService,
		translator:     NewTranslator(),
	}
}

func (s *StorageAuthService) RegisterUser(ctx context.Context, user *model.User, password string) (string, error) {

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

func (s *StorageAuthService) Authenticate(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.storageService.GetUserWithEmail(ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, "", fmt.Errorf("invalid credentials")
		}
		return nil, "", fmt.Errorf("error comparing password: %v", err)
	}

	return s.translator.ToUserModel(user), "", nil
}

func (s *StorageAuthService) IsAuthenticated(token string) (string, error) {
	return "", nil // Placeholder return
}
