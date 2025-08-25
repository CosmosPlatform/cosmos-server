package user

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	logMock "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/services/auth"
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

func TestHandleGetUsers(t *testing.T) {
	t.Run("success - get users", handleGetUsersSuccessful)
	t.Run("failure - get users internal error", handleGetUsersInternalError)
}

func TestHandleDeleteUser(t *testing.T) {
	t.Run("success - delete user", handleDeleteUserSuccessful)
	t.Run("failure - email required", handleDeleteUserEmailRequired)
	t.Run("failure - delete user internal error", handleDeleteUserInternalError)
}

func TestHandleGetCurrentUser(t *testing.T) {
	t.Run("success - get current user", handleGetCurrentUserSuccessful)
	t.Run("failure - user email not found in context", handleGetCurrentUserEmailNotInContext)
	t.Run("failure - get current user internal error", handleGetCurrentUserInternalError)
}

type mocks struct {
	controller      *gomock.Controller
	userServiceMock *userMock.MockService
	translatorMock  *userRouteMock.MockTranslator
	loggerMock      *logMock.MockLogger
}

func setUp(t *testing.T, mockedEmailForCtx string) (*gin.Engine, *mocks) {
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

	if mockedEmailForCtx != "" {
		// Add user email to context by creating a middleware
		router.Use(func(c *gin.Context) {
			c.Set(auth.UserEmailContextKey, mockedEmailForCtx)
			c.Next()
		})
	}

	AddAdminUserHandler(router.Group("/"), userServiceMock, translatorMock, loggerMock)

	return router, mocks
}

func handleRegisterUserSuccessful(t *testing.T) {
	router, mocks := setUp(t, "")

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
	router, mocks := setUp(t, "")

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
	router, mocks := setUp(t, "")

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
	router, mocks := setUp(t, "")

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
	router, mocks := setUp(t, "")

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
	router, mocks := setUp(t, "")

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
	router, mocks := setUp(t, "")

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

func handleGetUsersSuccessful(t *testing.T) {
	router, mocks := setUp(t, "")

	mockedUsers := []*model.User{
		{
			Username: "user1",
			Email:    "user1@example.com",
			Role:     "user",
		},
		{
			Username: "user2",
			Email:    "user2@example.com",
			Role:     "admin",
		},
	}

	mocks.userServiceMock.EXPECT().
		GetUsers(gomock.Any()).
		Return(mockedUsers, nil)

	mockedGetUsersResponse := &api.GetUsersResponse{
		Users: []*api.User{
			{
				Username: "user1",
				Email:    "user1@example.com",
				Role:     "user",
			},
			{
				Username: "user2",
				Email:    "user2@example.com",
				Role:     "admin",
			},
		},
	}

	mocks.translatorMock.EXPECT().
		ToGetUsersResponse(mockedUsers).
		Return(mockedGetUsersResponse)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.GetUsersResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, mockedGetUsersResponse, &actualResponse, "Expected response to match mocked response")
}

func handleGetUsersInternalError(t *testing.T) {
	router, mocks := setUp(t, "")

	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.userServiceMock.EXPECT().
		GetUsers(gomock.Any()).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	url := "/users"

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
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
	require.Equal(t, "Internal test error", actualResponse.Error)
}

func handleDeleteUserSuccessful(t *testing.T) {
	router, mocks := setUp(t, "")

	mockedEmail := "test@example.com"

	mocks.userServiceMock.EXPECT().
		DeleteUser(gomock.Any(), mockedEmail).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users?email=" + mockedEmail

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code, "Expected status code 204")
}

func handleDeleteUserEmailRequired(t *testing.T) {
	router, mocks := setUp(t, "")

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users"

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
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
	require.Equal(t, "Email query parameter is required", actualResponse.Error)
}

func handleDeleteUserInternalError(t *testing.T) {
	router, mocks := setUp(t, "")

	mockedEmail := "test@example.com"
	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.userServiceMock.EXPECT().
		DeleteUser(gomock.Any(), mockedEmail).
		Return(mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users?email=" + mockedEmail

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
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
	require.Equal(t, "Internal test error", actualResponse.Error)
}

func handleGetCurrentUserSuccessful(t *testing.T) {
	mockedUserEmail := "test@example.com"

	router, mocks := setUp(t, mockedUserEmail)

	mockedUser := &model.User{
		Username: "testUser",
		Email:    mockedUserEmail,
		Role:     "user",
	}

	mocks.userServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), mockedUserEmail).
		Return(mockedUser, nil)

	mockedGetCurrentUserResponse := &api.GetCurrentUserResponse{
		User: api.User{
			Username: "testUser",
			Email:    mockedUserEmail,
			Role:     "user",
		},
	}

	mocks.translatorMock.EXPECT().
		ToGetCurrentUserResponse(mockedUser).
		Return(mockedGetCurrentUserResponse)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users/me"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.GetCurrentUserResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, mockedGetCurrentUserResponse, &actualResponse, "Expected response to match mocked response")
}

func handleGetCurrentUserEmailNotInContext(t *testing.T) {
	router, mocks := setUp(t, "")

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users/me"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusUnauthorized, recorder.Code, "Expected status code 401")
	require.Equal(t, "user email not found in context", actualResponse.Error)
}

func handleGetCurrentUserInternalError(t *testing.T) {
	mockedUserEmail := "test@example.com"

	router, mocks := setUp(t, mockedUserEmail)

	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.userServiceMock.EXPECT().
		GetUserWithEmail(gomock.Any(), mockedUserEmail).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/users/me"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
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
	require.Equal(t, "Internal test error", actualResponse.Error)
}
