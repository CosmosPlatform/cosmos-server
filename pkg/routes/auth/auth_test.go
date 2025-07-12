package auth

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/test"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"

	authMock "cosmos-server/pkg/auth/mock"
	logMock "cosmos-server/pkg/log/mock"
	authRouteMock "cosmos-server/pkg/routes/auth/mock"
)

func TestHandleLogin(t *testing.T) {
	t.Run("success - successful login", handleLoginSuccessful)
	t.Run("failure - login error", handleLoginBadCredentials)
	t.Run("failure - email required", handleLoginEmailRequired)
	t.Run("failure - password required", handleLoginPasswordRequired)
}

type mocks struct {
	controller      *gomock.Controller
	authServiceMock *authMock.MockService
	translatorMock  *authRouteMock.MockTranslator
	loggerMock      *logMock.MockLogger
}

func setUp(t *testing.T) (*gin.Engine, *mocks) {
	ctrl := gomock.NewController(t)

	authServiceMock := authMock.NewMockService(ctrl)
	translatorMock := authRouteMock.NewMockTranslator(ctrl)
	loggerMock := logMock.NewMockLogger(ctrl)

	mocks := &mocks{
		controller:      ctrl,
		authServiceMock: authServiceMock,
		translatorMock:  translatorMock,
		loggerMock:      loggerMock,
	}

	router := test.NewRouter(loggerMock)

	AddAuthHandler(router.Group("/"), authServiceMock, loggerMock)

	return router, mocks
}

func handleLoginSuccessful(t *testing.T) {
	router, mocks := setUp(t)

	mockedEmail := "test@example.com"
	mockedUsername := "testuser"
	mockedPassword := "test123"
	mockedRole := "user"

	mockedToken := "mockedToken123"

	mockedAuthenticateRequest := &api.AuthenticateRequest{
		Email:    mockedEmail,
		Password: mockedPassword,
	}

	mockedUser := &model.User{
		Username: mockedUsername,
		Email:    mockedEmail,
		Role:     mockedRole,
	}

	mocks.authServiceMock.EXPECT().
		Authenticate(gomock.Any(), mockedEmail, mockedPassword).
		Return(mockedUser, mockedToken, nil)

	mockedAuthenticateResponse := &api.AuthenticateResponse{
		User: api.User{
			Username: mockedUsername,
			Email:    mockedEmail,
			Role:     mockedRole,
		},
		Token: mockedToken,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/auth/login"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedAuthenticateRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.AuthenticateResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, 200, recorder.Code, "Expected status code 200")
	require.Equal(t, mockedAuthenticateResponse.User, actualResponse.User)
}

func handleLoginBadCredentials(t *testing.T) {
	router, mocks := setUp(t)

	mockedEmail := "testErroneo@example.com"
	mockedPassword := "test123"

	mockedAuthenticateRequest := &api.AuthenticateRequest{
		Email:    mockedEmail,
		Password: mockedPassword,
	}

	mockedError := errors.NewUnauthorizedError("user not found")

	mockedErrorResponse := api.ErrorResponse{
		Error:      "user not found",
		StatusCode: 401,
	}

	mocks.authServiceMock.EXPECT().
		Authenticate(gomock.Any(), mockedEmail, mockedPassword).
		Return(nil, "", mockedError)

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/auth/login"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedAuthenticateRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, 401, recorder.Code, "Expected status code 401")
	require.Equal(t, mockedErrorResponse.Error, actualResponse.Error)
}

func handleLoginEmailRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedPassword := "test123"

	mockedAuthenticateRequest := &api.AuthenticateRequest{
		Password: mockedPassword,
	}

	mockedError := api.ErrorResponse{
		Error:      "email: cannot be blank.",
		StatusCode: 400,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	url := "/auth/login"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedAuthenticateRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, 400, recorder.Code, "Expected status code 400")
	require.Equal(t, mockedError.Error, actualResponse.Error)
}

func handleLoginPasswordRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedEmail := "test@example.com"

	mockedAuthenticateRequest := &api.AuthenticateRequest{
		Email: mockedEmail,
	}

	mockedError := api.ErrorResponse{
		Error:      "password: cannot be blank.",
		StatusCode: 400,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	url := "/auth/login"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedAuthenticateRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, 400, recorder.Code, "Expected status code 400")
	require.Equal(t, mockedError.Error, actualResponse.Error)
}
