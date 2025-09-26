package monitoring

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	applicationMock "cosmos-server/pkg/services/application/mock"
	monitoringMock "cosmos-server/pkg/services/monitoring/mock"
	"cosmos-server/pkg/test"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleUpdateApplicationMonitoring(t *testing.T) {
	t.Run("success - update application monitoring", handleUpdateApplicationMonitoringSuccess)
	t.Run("failure - application not found", handleUpdateApplicationMonitoringApplicationNotFound)
	t.Run("failure - internal server error", handleUpdateApplicationMonitoringInternalServerError)
}

func TestHandleGetApplicationInteractions(t *testing.T) {
	t.Run("success - get application interactions", handleGetApplicationInteractionsSuccess)
	t.Run("failure - application not found", handleGetApplicationInteractionsApplicationNotFound)
	t.Run("failure - monitoring service error", handleGetApplicationInteractionsMonitoringServiceError)
}

func TestHandleGetApplicationsInteractions(t *testing.T) {
	t.Run("success - get applications interactions", handleGetApplicationsInteractionsSuccess)
	t.Run("failure - monitoring service not found", handleGetApplicationsInteractionsNotFound)
	t.Run("failure - monitoring service internal error", handleGetApplicationsInteractionsInternalServerError)
}

type mocks struct {
	controller             *gomock.Controller
	monitoringServiceMock  *monitoringMock.MockService
	applicationServiceMock *applicationMock.MockService
	loggerMock             *log.MockLogger
}

func setUp(t *testing.T) (*gin.Engine, *mocks) {
	ctrl := gomock.NewController(t)

	monitoringServiceMock := monitoringMock.NewMockService(ctrl)
	applicationServiceMock := applicationMock.NewMockService(ctrl)
	loggerMock := log.NewMockLogger(ctrl)

	mocks := &mocks{
		controller:             ctrl,
		monitoringServiceMock:  monitoringServiceMock,
		applicationServiceMock: applicationServiceMock,
		loggerMock:             loggerMock,
	}

	router := test.NewRouter(loggerMock)

	AddAuthenticatedMonitoringHandler(router.Group("/"), monitoringServiceMock, applicationServiceMock, NewTranslator(), loggerMock)

	return router, mocks
}

func handleUpdateApplicationMonitoringSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedApplicationName := "test-application"
	mockedDescription := "Test application description"
	gitProvider := "github"
	gitRepositoryOwner := "test-owner"
	gitRepositoryName := "test-repo"
	gitRepositoryBranch := "main"

	modelApplication := &model.Application{
		Name:        mockedApplicationName,
		Description: mockedDescription,
		Team:        nil,
		GitInformation: &model.GitInformation{
			Provider:         gitProvider,
			RepositoryOwner:  gitRepositoryOwner,
			RepositoryName:   gitRepositoryName,
			RepositoryBranch: gitRepositoryBranch,
		},
	}

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedApplicationName).
		Return(modelApplication, nil)

	mocks.monitoringServiceMock.EXPECT().
		UpdateApplicationInformation(gomock.Any(), modelApplication).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := fmt.Sprintf("/monitoring/update/%s", mockedApplicationName)

	request, recorder, err := test.NewHTTPRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code, "Expected status code 204 No Content")
}

func handleUpdateApplicationMonitoringApplicationNotFound(t *testing.T) {
	router, mocks := setUp(t)

	mockedApplicationName := "non-existent-application"

	mockedError := errors.NewNotFoundError("application not found")

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedApplicationName).
		Return(nil, mockedError)

	url := fmt.Sprintf("/monitoring/update/%s", mockedApplicationName)

	request, recorder, err := test.NewHTTPRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Equal(t, "application not found", actualResponse.Error)
}

func handleUpdateApplicationMonitoringInternalServerError(t *testing.T) {
	router, mocks := setUp(t)

	mockedApplicationName := "test-application"
	mockedDescription := "Test application description"
	gitProvider := "github"
	gitRepositoryOwner := "test-owner"
	gitRepositoryName := "test-repo"
	gitRepositoryBranch := "main"

	modelApplication := &model.Application{
		Name:        mockedApplicationName,
		Description: mockedDescription,
		Team:        nil,
		GitInformation: &model.GitInformation{
			Provider:         gitProvider,
			RepositoryOwner:  gitRepositoryOwner,
			RepositoryName:   gitRepositoryName,
			RepositoryBranch: gitRepositoryBranch,
		},
	}

	mockedError := errors.NewInternalServerError("internal server error")

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedApplicationName).
		Return(modelApplication, nil)

	mocks.monitoringServiceMock.EXPECT().
		UpdateApplicationInformation(gomock.Any(), modelApplication).
		Return(mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := fmt.Sprintf("/monitoring/update/%s", mockedApplicationName)

	request, recorder, err := test.NewHTTPRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, "internal server error", actualResponse.Error)
}

func handleGetApplicationInteractionsSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedApplicationName := "test-application"
	mockedDescription := "Test application description"
	gitProvider := "github"
	gitRepositoryOwner := "test-owner"
	gitRepositoryName := "test-repo"
	gitRepositoryBranch := "main"

	modelApplication := &model.Application{
		Name:        mockedApplicationName,
		Description: mockedDescription,
		Team:        nil,
		GitInformation: &model.GitInformation{
			Provider:         gitProvider,
			RepositoryOwner:  gitRepositoryOwner,
			RepositoryName:   gitRepositoryName,
			RepositoryBranch: gitRepositoryBranch,
		},
	}

	mockedInteractions := &model.ApplicationsInteractions{
		ApplicationsInvolved: map[string]*model.Application{
			mockedApplicationName: modelApplication,
		},
	}

	expectedResponse := api.GetApplicationsInteractionsResponse{
		ApplicationsInvolved: map[string]api.ApplicationInformation{
			mockedApplicationName: {
				Name: mockedApplicationName,
				Team: "",
			},
		},
	}

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedApplicationName).
		Return(modelApplication, nil)

	mocks.monitoringServiceMock.EXPECT().
		GetApplicationInteractions(gomock.Any(), mockedApplicationName).
		Return(mockedInteractions, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := fmt.Sprintf("/monitoring/interactions/%s", mockedApplicationName)

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.GetApplicationsInteractionsResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, expectedResponse, actualResponse)
}

func handleGetApplicationInteractionsApplicationNotFound(t *testing.T) {
	router, mocks := setUp(t)

	mockedApplicationName := "non-existent-application"

	mockedError := errors.NewNotFoundError("application not found")

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedApplicationName).
		Return(nil, mockedError)

	url := fmt.Sprintf("/monitoring/interactions/%s", mockedApplicationName)

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Equal(t, "application not found", actualResponse.Error)
}

func handleGetApplicationInteractionsMonitoringServiceError(t *testing.T) {
	router, mocks := setUp(t)

	mockedApplicationName := "test-application"
	mockedDescription := "Test application description"
	gitProvider := "github"
	gitRepositoryOwner := "test-owner"
	gitRepositoryName := "test-repo"
	gitRepositoryBranch := "main"

	modelApplication := &model.Application{
		Name:        mockedApplicationName,
		Description: mockedDescription,
		Team:        nil,
		GitInformation: &model.GitInformation{
			Provider:         gitProvider,
			RepositoryOwner:  gitRepositoryOwner,
			RepositoryName:   gitRepositoryName,
			RepositoryBranch: gitRepositoryBranch,
		},
	}

	mockedError := errors.NewInternalServerError("failed to retrieve interactions")

	mocks.applicationServiceMock.EXPECT().
		GetApplication(gomock.Any(), mockedApplicationName).
		Return(modelApplication, nil)

	mocks.monitoringServiceMock.EXPECT().
		GetApplicationInteractions(gomock.Any(), mockedApplicationName).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	url := fmt.Sprintf("/monitoring/interactions/%s", mockedApplicationName)

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, "failed to retrieve interactions", actualResponse.Error)
}

func handleGetApplicationsInteractionsSuccess(t *testing.T) {
	router, mocks := setUp(t)

	teams := []string{"team1", "team2"}
	includeNeighbors := true

	mockedFilters := model.ApplicationDependencyFilter{
		Teams:            teams,
		IncludeNeighbors: includeNeighbors,
	}

	mockedInteractions := &model.ApplicationsInteractions{
		ApplicationsInvolved: map[string]*model.Application{
			"app1": {Name: "app1"},
		},
	}

	expectedResponse := api.GetApplicationsInteractionsResponse{
		ApplicationsInvolved: map[string]api.ApplicationInformation{
			"app1": {Name: "app1", Team: ""},
		},
	}

	mocks.monitoringServiceMock.EXPECT().
		GetApplicationsInteractions(gomock.Any(), mockedFilters).
		Return(mockedInteractions, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/monitoring/interactions?teams=team1,team2&includeNeighbors=true"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	actualResponse := api.GetApplicationsInteractionsResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, expectedResponse, actualResponse)
}

func handleGetApplicationsInteractionsNotFound(t *testing.T) {
	router, mocks := setUp(t)

	mockedError := errors.NewNotFoundError("no interactions found")

	mocks.monitoringServiceMock.EXPECT().
		GetApplicationsInteractions(gomock.Any(), gomock.Any()).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	url := "/monitoring/interactions?teams=team1"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	require.NoError(t, err)

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Equal(t, "no interactions found", actualResponse.Error)
}

func handleGetApplicationsInteractionsInternalServerError(t *testing.T) {
	router, mocks := setUp(t)

	mockedError := errors.NewInternalServerError("internal error")

	mocks.monitoringServiceMock.EXPECT().
		GetApplicationsInteractions(gomock.Any(), gomock.Any()).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	mocks.loggerMock.EXPECT().
		Errorf(gomock.Any(), gomock.Any())

	url := "/monitoring/interactions?teams=team1"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	require.NoError(t, err)

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, "internal error", actualResponse.Error)
}
