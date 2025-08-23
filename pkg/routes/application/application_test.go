package application

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	applicationMock "cosmos-server/pkg/services/application/mock"
	"cosmos-server/pkg/test"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestHandleCreateApplication(t *testing.T) {
	t.Run("success - create application", handleCreateApplicationSuccess)
	t.Run("failure - name required", handleCreateApplicationNameRequired)
	t.Run("failure - insert application error", handleCreateApplicationInsertApplicationError)
}

func TestHandleGetApplication(t *testing.T) {
	t.Run("success - get application", handleGetApplicationSuccess)
	t.Run("failure - get application error", handleGetApplicationError)
	t.Run("failure - get application error does not exist", handleGetApplicationErrorDoesNotExist)
}

func TestHandleGetApplications(t *testing.T) {
	t.Run("success - get applications", handleGetApplicationsSuccess)
	t.Run("success - get applications with name filter", handleGetApplicationsWithNameFilter)
	t.Run("failure - get applications error", handleGetApplicationsError)
}

type mocks struct {
	controller             *gomock.Controller
	applicationServiceMock *applicationMock.MockService
	loggerMock             *log.MockLogger
}

func setUp(t *testing.T) (*gin.Engine, *mocks) {
	ctrl := gomock.NewController(t)

	applicationServiceMock := applicationMock.NewMockService(ctrl)
	loggerMock := log.NewMockLogger(ctrl)

	mocks := &mocks{
		controller:             ctrl,
		applicationServiceMock: applicationServiceMock,
		loggerMock:             loggerMock,
	}

	router := test.NewRouter(loggerMock)

	AddApplicationHandler(router.Group("/"), applicationServiceMock, NewTranslator(), loggerMock)

	return router, mocks
}

func handleCreateApplicationSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"

	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
		Name:        mockedName,
		Description: mockedDescription,
		Team:        mockedTeam,
	}

	expectedResponse := &api.CreateApplicationResponse{
		Application: &api.Application{
			Name:        mockedName,
			Description: mockedDescription,
			Team:        &api.Team{Name: mockedTeam},
		},
	}

	mocks.applicationServiceMock.EXPECT().
		AddApplication(gomock.Any(), mockedName, mockedDescription, mockedTeam).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedCreateApplicationRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.CreateApplicationResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusCreated, recorder.Code, "Expected status code 201")
	require.Equal(t, expectedResponse, &actualResponse, "Response body mismatch")
}

func handleCreateApplicationNameRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedDescription := "Test application description"
	mockedTeam := "test-team"

	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
		Description: mockedDescription,
		Team:        mockedTeam,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedCreateApplicationRequest)
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

func handleCreateApplicationInsertApplicationError(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"

	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
		Name:        mockedName,
		Description: mockedDescription,
		Team:        mockedTeam,
	}

	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.applicationServiceMock.EXPECT().
		AddApplication(gomock.Any(), mockedName, mockedDescription, mockedTeam).
		Return(mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedCreateApplicationRequest)
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

func handleGetApplicationSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"

	mockedApplication := &api.Application{
		Name:        mockedName,
		Description: mockedDescription,
		Team:        &api.Team{Name: mockedTeam},
	}

	mockedApplicationModel := &model.Application{
		Name:        mockedName,
		Description: mockedDescription,
		Team:        &model.Team{Name: mockedTeam},
	}

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedName).
		Return(mockedApplicationModel, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.Application{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, mockedApplication, &actualResponse, "Response body mismatch")
}

func handleGetApplicationError(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedName).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

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

func handleGetApplicationErrorDoesNotExist(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "nonexistent-app"
	mockedError := errors.NewNotFoundError("Application does not exist")

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedName).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

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

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Equal(t, "Application does not exist", actualResponse.Error)
}

func handleGetApplicationsSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedApplication1 := &model.Application{
		Name:        "test-app-1",
		Description: "Test application 1 description",
		Team:        &model.Team{Name: "test-team-1"},
	}

	mockedApplication2 := &model.Application{
		Name:        "test-app-2",
		Description: "Test application 2 description",
		Team:        &model.Team{Name: "test-team-2"},
	}

	mockedApplications := []*model.Application{mockedApplication1, mockedApplication2}

	expectedResponse := &api.GetApplicationsResponse{
		Applications: []*api.Application{
			{
				Name:        "test-app-1",
				Description: "Test application 1 description",
				Team:        &api.Team{Name: "test-team-1"},
			},
			{
				Name:        "test-app-2",
				Description: "Test application 2 description",
				Team:        &api.Team{Name: "test-team-2"},
			},
		},
	}

	mocks.applicationServiceMock.EXPECT().
		GetApplicationsWithFilter(gomock.Any(), "").
		Return(mockedApplications, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.GetApplicationsResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, expectedResponse, &actualResponse, "Response body mismatch")
}

func handleGetApplicationsWithNameFilter(t *testing.T) {
	router, mocks := setUp(t)

	mockedNameFilter := "test-app"
	mockedApplication := &model.Application{
		Name:        "test-app-1",
		Description: "Test application description",
		Team:        &model.Team{Name: "test-team"},
	}

	mockedApplications := []*model.Application{mockedApplication}

	expectedResponse := &api.GetApplicationsResponse{
		Applications: []*api.Application{
			{
				Name:        "test-app-1",
				Description: "Test application description",
				Team:        &api.Team{Name: "test-team"},
			},
		},
	}

	mocks.applicationServiceMock.EXPECT().
		GetApplicationsWithFilter(gomock.Any(), mockedNameFilter).
		Return(mockedApplications, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications?name=" + mockedNameFilter

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.GetApplicationsResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, expectedResponse, &actualResponse, "Response body mismatch")
}

func handleGetApplicationsError(t *testing.T) {
	router, mocks := setUp(t)

	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.applicationServiceMock.EXPECT().
		GetApplicationsWithFilter(gomock.Any(), "").
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications"

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
