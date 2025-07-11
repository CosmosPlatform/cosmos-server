package auth

import (
	"context"
	"cosmos-server/pkg/config"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage/obj"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"

	authMock "cosmos-server/pkg/auth/mock"
	logMock "cosmos-server/pkg/log/mock"
	storageMock "cosmos-server/pkg/storage/mock"
)

const (
	TestJWTSecret = "my_secret_key"
)

func TestAuthenticate(t *testing.T) {
	t.Run("authenticate - success", authenticateSuccess)
	t.Run("authenticate - wrong email", authenticateWrongEmail)
}

func TestIsAuthenticated(t *testing.T) {
	t.Run("isAuthenticated - success", isAuthenticatedSuccess)
	t.Run("isAuthenticated - failure", isAuthenticatedFailure)
}

type mocks struct {
	controller         *gomock.Controller
	configuration      config.AuthConfig
	storageServiceMock *storageMock.MockService
	translatorMock     *authMock.MockTranslator
	loggerMock         *logMock.MockLogger
}

func setUp(t *testing.T) (Service, *mocks) {
	ctrl := gomock.NewController(t)

	configuration := config.AuthConfig{
		JWTSecret: TestJWTSecret,
	}

	mocks := &mocks{
		controller:         ctrl,
		configuration:      configuration,
		storageServiceMock: storageMock.NewMockService(ctrl),
		translatorMock:     authMock.NewMockTranslator(ctrl),
		loggerMock:         logMock.NewMockLogger(ctrl),
	}

	authService := NewAuthService(configuration, mocks.storageServiceMock, mocks.translatorMock, mocks.loggerMock)

	return authService, mocks
}

func authenticateSuccess(t *testing.T) {
	authService, mocks := setUp(t)

	mockedUsername := "testUser"
	mockedEmail := "test@example.com"
	mockedPassword := "password123"
	mockedRole := "user"

	mockedEncryptedPassword, err := bcrypt.GenerateFromPassword([]byte(mockedPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate encrypted password: %v", err)
	}

	mockedObjUser := &obj.User{
		Email:             mockedEmail,
		Username:          mockedUsername,
		EncryptedPassword: string(mockedEncryptedPassword),
		Role:              mockedRole,
	}

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), mockedEmail).
		Return(mockedObjUser, nil)

	mockedModelUser := &model.User{
		Email:    mockedEmail,
		Username: mockedUsername,
		Role:     mockedRole,
	}

	mocks.translatorMock.EXPECT().
		ToUserModel(mockedObjUser).
		Return(mockedModelUser)

	_, _, err = authService.Authenticate(context.TODO(), mockedEmail, mockedPassword)

	require.NoError(t, err)
}

func authenticateWrongEmail(t *testing.T) {
	authService, mocks := setUp(t)

	mockedEmail := "testError@example.com"
	mockedPassword := "password123"

	mocks.storageServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), mockedEmail).
		Return(nil, errors.NewInternalServerError("User not found"))

	mocks.loggerMock.EXPECT().Errorf(gomock.Any(), gomock.Any())

	_, _, err := authService.Authenticate(context.TODO(), mockedEmail, mockedPassword)

	require.Error(t, err)
}

func isAuthenticatedSuccess(t *testing.T) {
	authService, _ := setUp(t)

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(TestJWTSecret))
	require.NoError(t, err)

	_, err = authService.IsAuthenticated(tokenString)
	require.NoError(t, err)
}

func isAuthenticatedFailure(t *testing.T) {
	authService, _ := setUp(t)

	invalidToken := "invalid_token"

	_, err := authService.IsAuthenticated(invalidToken)
	require.Error(t, err)
}
