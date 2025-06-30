package auth

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	UserEmailClaimKey      = "email"
	UserUsernameClaimKey   = "username"
	UserRoleClaimKey       = "role"
	UserExpirationClaimKey = "exp"
)

type Service interface {
	// Authenticate authenticates a user with the provided credentials and returns the user and a token if successful.
	Authenticate(ctx context.Context, email, password string) (*model.User, string, error)
	// IsAuthenticated checks if the provided token is valid for the user and returns the user ID if authenticated.
	IsAuthenticated(tokenString string) (*jwt.Token, error)
}

type authService struct {
	storageService storage.Service
	translator     Translator
	config         config.AuthConfig
	logger         log.Logger
}

func NewAuthService(config config.AuthConfig, storageService storage.Service, logger log.Logger) Service {
	return &authService{
		storageService: storageService,
		translator:     NewTranslator(),
		config:         config,
		logger:         logger,
	}
}

func (s *authService) Authenticate(ctx context.Context, email, password string) (*model.User, string, error) {
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

func (s *authService) IsAuthenticated(tokenString string) (*jwt.Token, error) {
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

func (s *authService) GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		UserEmailClaimKey:      user.Email,
		UserUsernameClaimKey:   user.Username,
		UserRoleClaimKey:       user.Role,
		UserExpirationClaimKey: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // TODO: Move to config.
	})

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
