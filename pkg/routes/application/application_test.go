package application

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	applicationMock "cosmos-server/pkg/services/application/mock"
	monitoringMock "cosmos-server/pkg/services/monitoring/mock"
	"cosmos-server/pkg/test"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleCreateApplication(t *testing.T) {
	t.Run("success - create application", handleCreateApplicationSuccess)
	t.Run("failure - name required", handleCreateApplicationNameRequired)
	t.Run("failure - insert application error", handleCreateApplicationInsertApplicationError)
}

func TestHandleCreateApplicationWithGitInformation(t *testing.T) {
	t.Run("success - create application with git information", handleCreateApplicationWithGitInformationSuccess)
	t.Run("failure - git information provider required", handleCreateApplicationGitInformationProviderRequired)
	t.Run("failure - git information repository owner required", handleCreateApplicationGitInformationRepositoryOwnerRequired)
	t.Run("failure - git information repository name required", handleCreateApplicationGitInformationRepositoryNameRequired)
	t.Run("failure - git information repository branch required", handleCreateApplicationGitInformationRepositoryBranchRequired)
}

func TestHandleGetApplication(t *testing.T) {
	t.Run("success - get application", handleGetApplicationSuccess)
	t.Run("success - get application with git information", handleGetApplicationWithGitInformationSuccess)
	t.Run("failure - get application error", handleGetApplicationError)
	t.Run("failure - get application error does not exist", handleGetApplicationErrorDoesNotExist)
}

func TestHandleGetApplications(t *testing.T) {
	t.Run("success - get applications", handleGetApplicationsSuccess)
	t.Run("success - get applications with name filter", handleGetApplicationsWithNameFilter)
	t.Run("failure - get applications error", handleGetApplicationsError)
}

func TestHandleGetApplicationByTeam(t *testing.T) {
	t.Run("success - get applications by team", handleGetApplicationByTeamSuccess)
	t.Run("success - get applications by team no applications", handleGetApplicationByTeamNoApplications)
	t.Run("failure - get applications by team error", handleGetApplicationByTeamError)
}

func TestHandleDeleteApplication(t *testing.T) {
	t.Run("success - delete application", handleDeleteApplicationSuccess)
	t.Run("failure - delete application error", handleDeleteApplicationError)
	t.Run("failure - delete application not found", handleDeleteApplicationNotFound)
}

func TestHandleUpdateApplication(t *testing.T) {
	t.Run("success - update application", handleUpdateApplicationSuccess)
	t.Run("success - update application with git information", handleUpdateApplicationWithGitInformationSuccess)
	t.Run("failure - update application with partial git information", handleUpdateApplicationWithPartialGitInformationError)
	t.Run("failure - update application error", handleUpdateApplicationError)
}

type mocks struct {
	controller             *gomock.Controller
	applicationServiceMock *applicationMock.MockService
	monitoringServiceMock  *monitoringMock.MockService
	loggerMock             *log.MockLogger
}

func setUp(t *testing.T) (*gin.Engine, *mocks) {
	ctrl := gomock.NewController(t)

	applicationServiceMock := applicationMock.NewMockService(ctrl)
	monitoringServiceMock := monitoringMock.NewMockService(ctrl)
	loggerMock := log.NewMockLogger(ctrl)

	mocks := &mocks{
		controller:             ctrl,
		applicationServiceMock: applicationServiceMock,
		monitoringServiceMock:  monitoringServiceMock,
		loggerMock:             loggerMock,
	}

	router := test.NewRouter(loggerMock)

	AddAuthenticatedApplicationHandler(router.Group("/"), applicationServiceMock, monitoringServiceMock, NewTranslator(), loggerMock)

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
		AddApplication(gomock.Any(), mockedName, mockedDescription, mockedTeam, nil).
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
		AddApplication(gomock.Any(), mockedName, mockedDescription, mockedTeam, nil).
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

	expectedResponse := &api.GetApplicationResponse{
		Application: mockedApplication,
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

	actualResponse := api.GetApplicationResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, expectedResponse, &actualResponse, "Response body mismatch")
}

func handleGetApplicationWithGitInformationSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"
	mockedGitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	mockedApplicationModel := &model.Application{
		Name:           mockedName,
		Description:    mockedDescription,
		Team:           &model.Team{Name: mockedTeam},
		GitInformation: mockedGitInformation,
	}

	expectedResponse := &api.GetApplicationResponse{
		Application: &api.Application{
			Name:        mockedName,
			Description: mockedDescription,
			Team:        &api.Team{Name: mockedTeam},
			GitInformation: &api.GitInformation{
				Provider:         "github",
				RepositoryOwner:  "test-owner",
				RepositoryName:   "test-repo",
				RepositoryBranch: "main",
			},
		},
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

	actualResponse := api.GetApplicationResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, expectedResponse, &actualResponse, "Response body mismatch")
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

func handleDeleteApplicationSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"

	mocks.applicationServiceMock.EXPECT().
		DeleteApplication(gomock.Any(), mockedName).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code, "Expected status code 204")
	require.Empty(t, recorder.Body.String(), "Expected empty response body")
}

func handleDeleteApplicationError(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.applicationServiceMock.EXPECT().
		DeleteApplication(gomock.Any(), mockedName).
		Return(mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

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

func handleDeleteApplicationNotFound(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "nonexistent-app"
	mockedError := errors.NewNotFoundError("Application not found")

	mocks.applicationServiceMock.EXPECT().
		DeleteApplication(gomock.Any(), mockedName).
		Return(mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

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

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Equal(t, "Application not found", actualResponse.Error)
}

func handleGetApplicationByTeamSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedTeamName := "test-team"
	mockedApplication1 := &model.Application{
		Name:        "test-app-1",
		Description: "Test application 1 description",
		Team:        &model.Team{Name: mockedTeamName},
	}

	mockedApplication2 := &model.Application{
		Name:        "test-app-2",
		Description: "Test application 2 description",
		Team:        &model.Team{Name: mockedTeamName},
	}

	mockedApplications := []*model.Application{mockedApplication1, mockedApplication2}

	expectedResponse := &api.GetApplicationsResponse{
		Applications: []*api.Application{
			{
				Name:        "test-app-1",
				Description: "Test application 1 description",
				Team:        &api.Team{Name: mockedTeamName},
			},
			{
				Name:        "test-app-2",
				Description: "Test application 2 description",
				Team:        &api.Team{Name: mockedTeamName},
			},
		},
	}

	mocks.applicationServiceMock.EXPECT().
		GetApplicationsByTeam(gomock.Any(), mockedTeamName).
		Return(mockedApplications, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/team/" + mockedTeamName

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

func handleGetApplicationByTeamError(t *testing.T) {
	router, mocks := setUp(t)

	mockedTeamName := "test-team"
	mockedError := errors.NewInternalServerError("Internal test error")

	mocks.applicationServiceMock.EXPECT().
		GetApplicationsByTeam(gomock.Any(), mockedTeamName).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/team/" + mockedTeamName

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

func handleGetApplicationByTeamNoApplications(t *testing.T) {
	router, mocks := setUp(t)

	mockedTeamName := "empty-team"

	expectedResponse := &api.GetApplicationsResponse{
		Applications: []*api.Application{},
	}

	mocks.applicationServiceMock.EXPECT().
		GetApplicationsByTeam(gomock.Any(), mockedTeamName).
		Return([]*model.Application{}, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/team/" + mockedTeamName

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

func handleCreateApplicationWithGitInformationSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"
	mockedGitInformation := &api.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
		Name:           mockedName,
		Description:    mockedDescription,
		Team:           mockedTeam,
		GitInformation: mockedGitInformation,
	}

	expectedResponse := &api.CreateApplicationResponse{
		Application: &api.Application{
			Name:           mockedName,
			Description:    mockedDescription,
			Team:           &api.Team{Name: mockedTeam},
			GitInformation: mockedGitInformation,
		},
	}

	expectedGitInfo := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	mocks.applicationServiceMock.EXPECT().
		AddApplication(gomock.Any(), mockedName, mockedDescription, mockedTeam, expectedGitInfo).
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

func handleCreateApplicationGitInformationProviderRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"
	mockedGitInformation := &api.GitInformation{
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
		Name:           mockedName,
		Description:    mockedDescription,
		Team:           mockedTeam,
		GitInformation: mockedGitInformation,
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
	require.Contains(t, actualResponse.Error, "provider")
}

func handleCreateApplicationGitInformationRepositoryOwnerRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"
	mockedGitInformation := &api.GitInformation{
		Provider:         "github",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
		Name:           mockedName,
		Description:    mockedDescription,
		Team:           mockedTeam,
		GitInformation: mockedGitInformation,
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
	require.Contains(t, actualResponse.Error, "repositoryOwner")
}

func handleCreateApplicationGitInformationRepositoryNameRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"
	mockedGitInformation := &api.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryBranch: "main",
	}

	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
		Name:           mockedName,
		Description:    mockedDescription,
		Team:           mockedTeam,
		GitInformation: mockedGitInformation,
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
	require.Contains(t, actualResponse.Error, "repositoryName")
}

func handleCreateApplicationGitInformationRepositoryBranchRequired(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedDescription := "Test application description"
	mockedTeam := "test-team"
	mockedGitInformation := &api.GitInformation{
		Provider:        "github",
		RepositoryOwner: "test-owner",
		RepositoryName:  "test-repo",
	}

	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
		Name:           mockedName,
		Description:    mockedDescription,
		Team:           mockedTeam,
		GitInformation: mockedGitInformation,
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
	require.Contains(t, actualResponse.Error, "repositoryBranch")
}

func handleUpdateApplicationSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	newName := "updated-app"
	newDescription := "Updated description"
	newTeam := "updated-team"

	mockedUpdateRequest := &api.UpdateApplicationRequest{
		Name:        &newName,
		Description: &newDescription,
		Team:        &newTeam,
	}

	mockedUpdatedApplication := &model.Application{
		Name:        newName,
		Description: newDescription,
		Team:        &model.Team{Name: newTeam},
	}

	expectedResponse := &api.UpdateApplicationResponse{
		Application: &api.Application{
			Name:        newName,
			Description: newDescription,
			Team:        &api.Team{Name: newTeam},
		},
	}

	expectedUpdateData := &model.ApplicationUpdate{
		Name:        &newName,
		Description: &newDescription,
		Team:        &newTeam,
	}

	mocks.applicationServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), mockedName, expectedUpdateData).
		Return(mockedUpdatedApplication, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

	request, recorder, err := test.NewHTTPRequest("PUT", url, mockedUpdateRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.UpdateApplicationResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, expectedResponse, &actualResponse, "Response body mismatch")
}

func handleUpdateApplicationWithGitInformationSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	newName := "updated-app"
	newDescription := "Updated description"
	newTeam := "updated-team"
	mockedGitInformation := &api.GitInformation{
		Provider:         "gitlab",
		RepositoryOwner:  "updated-owner",
		RepositoryName:   "updated-repo",
		RepositoryBranch: "develop",
	}

	mockedUpdateRequest := &api.UpdateApplicationRequest{
		Name:           &newName,
		Description:    &newDescription,
		Team:           &newTeam,
		GitInformation: mockedGitInformation,
	}

	mockedUpdatedApplication := &model.Application{
		Name:        newName,
		Description: newDescription,
		Team:        &model.Team{Name: newTeam},
		GitInformation: &model.GitInformation{
			Provider:         "gitlab",
			RepositoryOwner:  "updated-owner",
			RepositoryName:   "updated-repo",
			RepositoryBranch: "develop",
		},
	}

	expectedResponse := &api.UpdateApplicationResponse{
		Application: &api.Application{
			Name:        newName,
			Description: newDescription,
			Team:        &api.Team{Name: newTeam},
			GitInformation: &api.GitInformation{
				Provider:         "gitlab",
				RepositoryOwner:  "updated-owner",
				RepositoryName:   "updated-repo",
				RepositoryBranch: "develop",
			},
		},
	}

	expectedUpdateData := &model.ApplicationUpdate{
		Name:        &newName,
		Description: &newDescription,
		Team:        &newTeam,
		GitInformation: &model.GitInformation{
			Provider:         "gitlab",
			RepositoryOwner:  "updated-owner",
			RepositoryName:   "updated-repo",
			RepositoryBranch: "develop",
		},
	}

	mocks.applicationServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), mockedName, expectedUpdateData).
		Return(mockedUpdatedApplication, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

	request, recorder, err := test.NewHTTPRequest("PUT", url, mockedUpdateRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.UpdateApplicationResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	require.Equal(t, expectedResponse, &actualResponse, "Response body mismatch")
}

func handleUpdateApplicationWithPartialGitInformationError(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	mockedGitInformation := &api.GitInformation{
		Provider:        "github",
		RepositoryOwner: "test-owner",
		RepositoryName:  "test-repo",
	}

	mockedUpdateRequest := &api.UpdateApplicationRequest{
		GitInformation: mockedGitInformation,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

	request, recorder, err := test.NewHTTPRequest("PUT", url, mockedUpdateRequest)
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
	require.Contains(t, actualResponse.Error, "repositoryBranch")
}

func handleUpdateApplicationError(t *testing.T) {
	router, mocks := setUp(t)

	mockedName := "test-app"
	newName := "updated-name"

	mockedUpdateRequest := &api.UpdateApplicationRequest{
		Name: &newName,
	}

	mockedError := errors.NewInternalServerError("Internal test error")

	expectedUpdateData := &model.ApplicationUpdate{
		Name: &newName,
	}

	mocks.applicationServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), mockedName, expectedUpdateData).
		Return(nil, mockedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/applications/" + mockedName

	request, recorder, err := test.NewHTTPRequest("PUT", url, mockedUpdateRequest)
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
