package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	logMock "cosmos-server/pkg/log/mock"
	userMock "cosmos-server/pkg/services/user/mock"
	"cosmos-server/pkg/test"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"

	userRouteMock "cosmos-server/pkg/routes/user/mock"
)

func TestHandleRegisterUser(t *testing.T) {
	t.Run("success - register user", handleRegisterUserSuccessful)
	t.Run("failure - username required", handleRegisterUserUsernameRequired)
	t.Run("failure - email required", handleRegisterUserEmailRequired)
	t.Run("failure - password required", handleRegisterUserPasswordRequired)
	t.Run("failure - role invalid", handleRegisterUserRoleInvalid)
	t.Run("failure - role required", handleRegisterUserRoleRequired)
	t.Run("failure - register user internal error", handleRegisterUserInternalError)
}

type mocks struct {
	controller      *gomock.Controller
	userServiceMock *userMock.MockService
	translatorMock  *userRouteMock.MockTranslator
	loggerMock      *logMock.MockLogger
}

func setUp(t *testing.T) (*gin.Engine, *mocks) {
	ctrl := gomock.NewController(t)

	userServiceMock := userMock.NewMockService(ctrl)
	translatorMock := userRouteMock.NewMockTranslator(ctrl)
	loggerMock := logMock.NewMockLogger(ctrl)

	mocks := &mocks{
		controller:      ctrl,
		userServiceMock: userServiceMock,
		translatorMock:  translatorMock,
		loggerMock:      loggerMock,
	}

	router := test.NewRouter(loggerMock)

	AddAdminUserHandler(router.Group("/"), userServiceMock, translatorMock, loggerMock)

	return router, mocks
}

func handleRegisterUserSuccessful(t *testing.T) {
	router, mocks := setUp(t)

	mockedUsername := "testUser"
	mockedEmail := "test@example.com"
	mockedPassword := "test12345678"
	mockedRole := "user"

	mockedRegisterUserRequest := &api.RegisterUserRequest{
		Username: mockedUsername,
		Email:    mockedEmail,
		Password: mockedPassword,
		Role:     mockedRole,
	}

	mocks.userServiceMock.EXPECT().
		RegisterUser(gomock.Any(), mockedUsername, mockedEmail, mockedPassword, mockedRole).
		Return(nil)

	mockedRegisterUserResponse := &api.RegisterUserResponse{
		User: api.User{
			Username: mockedUsername,
			Email:    mockedEmail,
			Role:     mockedRole,
		},
	}

	mocks.translatorMock.EXPECT().
		ToRegisterUserResponse(mockedUsername, mockedEmail, mockedRole).
		Return(mockedRegisterUserResponse)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedRegisterUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.RegisterUserResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusCreated, recorder.Code, "Expected status code 201")
	require.Equal(t, mockedRegisterUserResponse, &actualResponse, "Expected response to match mocked response")
}

func handleRegisterUserUsernameRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedEmail := "test@example.com"
	mockedPassword := "test12345678"
	mockedRole := "user"

	mockedRegisterUserRequest := &api.RegisterUserRequest{
		Email:    mockedEmail,
		Password: mockedPassword,
		Role:     mockedRole,
	}

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedRegisterUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
}

func handleRegisterUserEmailRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedUsername := "testUser"
	mockedPassword := "test12345678"
	mockedRole := "user"

	mockedRegisterUserRequest := &api.RegisterUserRequest{
		Username: mockedUsername,
		Password: mockedPassword,
		Role:     mockedRole,
	}

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedRegisterUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
}

func handleRegisterUserPasswordRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedUsername := "testUser"
	mockedEmail := "test@example.com"
	mockedRole := "user"

	mockedRegisterUserRequest := &api.RegisterUserRequest{
		Username: mockedUsername,
		Email:    mockedEmail,
		Role:     mockedRole,
	}

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedRegisterUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
}

func handleRegisterUserRoleInvalid(t *testing.T) {
	router, mocks := setUp(t)

	mockedUsername := "testUser"
	mockedEmail := "test@example.com"
	mockedPassword := "test12345678"
	mockedRole := "invalidRole"

	mockedRegisterUserRequest := &api.RegisterUserRequest{
		Username: mockedUsername,
		Email:    mockedEmail,
		Password: mockedPassword,
		Role:     mockedRole,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedRegisterUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
}

func handleRegisterUserRoleRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedUsername := "testUser"
	mockedEmail := "test@example.com"
	mockedPassword := "test12345678"

	mockedRegisterUserRequest := &api.RegisterUserRequest{
		Username: mockedUsername,
		Email:    mockedEmail,
		Password: mockedPassword,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedRegisterUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
}

func handleRegisterUserInternalError(t *testing.T) {
	router, mocks := setUp(t)

	mockedUsername := "testUser"
	mockedEmail := "test@example.com"
	mockedPassword := "test12345678"
	mockedRole := "user"

	mockedRegisterUserRequest := &api.RegisterUserRequest{
		Username: mockedUsername,
		Email:    mockedEmail,
		Password: mockedPassword,
		Role:     mockedRole,
	}

	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.userServiceMock.EXPECT().
		RegisterUser(gomock.Any(), mockedUsername, mockedEmail, mockedPassword, mockedRole).
		Return(mockedError)

	mockedErrorResponse := &api.ErrorResponse{
		Error: "Internal test error",
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedRegisterUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, mockedErrorResponse, &actualResponse)
}
