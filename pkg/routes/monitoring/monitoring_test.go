package monitoring

import (
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	applicationMock "cosmos-server/pkg/services/application/mock"
	monitoringMock "cosmos-server/pkg/services/monitoring/mock"
	"cosmos-server/pkg/test"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleUpdateApplicationMonitoring(t *testing.T) {
	t.Run("success - update application monitoring", handleUpdateApplicationMonitoringSuccess)
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
