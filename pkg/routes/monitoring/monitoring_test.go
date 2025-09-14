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
