package auth

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type StorageAuthService struct {
	storageService storage.Service
	translator     Translator
	config         config.AuthConfig
}

func NewStorageAuthService(config config.AuthConfig, storageService storage.Service) *StorageAuthService {
	return &StorageAuthService{
		storageService: storageService,
		translator:     NewTranslator(),
		config:         config,
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

	userModel := s.translator.ToUserModel(user)

	token, err := s.GenerateToken(userModel)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return userModel, token, nil
}

func (s *StorageAuthService) IsAuthenticated(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (s *StorageAuthService) GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // TODO: Move to config.
		"role":    user.Role,
	})

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
