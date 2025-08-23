package application

import (
	"cosmos-server/api"
	log "cosmos-server/pkg/log/mock"
	applicationMock "cosmos-server/pkg/services/application/mock"
	"cosmos-server/pkg/test"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

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

func TestHandleCreateApplication(t *testing.T) {
	t.Run("success - create application", handleCreateApplicationSuccess)
	//t.Run("failure - name required", handleCreateApplicationNameRequired)
	//t.Run("failure - description required", handleCreateApplicationDescriptionRequired)
	//t.Run("failure - team required", handleCreateApplicationTeamRequired)
	//t.Run("failure - insert application error", handleCreateApplicationInsertApplicationError)
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

//func handleCreateApplicationNameRequired(t *testing.T) {
//	router, mocks := setUp(t)
//
//	mockedDescription := "Test application description"
//	mockedTeam := "test-team"
//
//	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
//		Description: mockedDescription,
//		Team:        mockedTeam,
//	}
//
//	mocks.loggerMock.EXPECT().
//		Errorf(gomock.Any(), gomock.Any())
//
//	url := "/applications"
//
//	request, recorder, err := test.NewHTTPRequest("POST", url, mockedCreateApplicationRequest)
//	if err != nil {
//		t.Fatalf("Failed to create request: %v", err)
//	}
//
//	router.ServeHTTP(recorder, request)
//
//	actualResponse := api.ErrorResponse{}
//	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
//	if err != nil {
//		t.Fatalf("Failed to decode response: %v", err)
//	}
//
//	require.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
//}
//
//func handleCreateApplicationDescriptionRequired(t *testing.T) {
//	router, mocks := setUp(t)
//
//	mockedName := "test-app"
//	mockedTeam := "test-team"
//
//	// Create request with description exceeding 500 characters to trigger validation error
//	longDescription := make([]byte, 501)
//	for i := range longDescription {
//		longDescription[i] = 'a'
//	}
//
//	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
//		Name:        mockedName,
//		Description: string(longDescription),
//		Team:        mockedTeam,
//	}
//
//	url := "/applications"
//
//	request, recorder, err := test.NewHTTPRequest("POST", url, mockedCreateApplicationRequest)
//	if err != nil {
//		t.Fatalf("Failed to create request: %v", err)
//	}
//
//	router.ServeHTTP(recorder, request)
//
//	actualResponse := api.ErrorResponse{}
//	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
//	if err != nil {
//		t.Fatalf("Failed to decode response: %v", err)
//	}
//
//	require.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
//}
//
//func handleCreateApplicationTeamRequired(t *testing.T) {
//	router, mocks := setUp(t)
//
//	mockedName := "test-app"
//	mockedDescription := "Test application description"
//
//	// Create request with team exceeding 100 characters to trigger validation error
//	longTeam := make([]byte, 101)
//	for i := range longTeam {
//		longTeam[i] = 'a'
//	}
//
//	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
//		Name:        mockedName,
//		Description: mockedDescription,
//		Team:        string(longTeam),
//	}
//
//	url := "/applications"
//
//	request, recorder, err := test.NewHTTPRequest("POST", url, mockedCreateApplicationRequest)
//	if err != nil {
//		t.Fatalf("Failed to create request: %v", err)
//	}
//
//	router.ServeHTTP(recorder, request)
//
//	actualResponse := api.ErrorResponse{}
//	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
//	if err != nil {
//		t.Fatalf("Failed to decode response: %v", err)
//	}
//
//	require.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400")
//}
//
//func handleCreateApplicationInsertApplicationError(t *testing.T) {
//	router, mocks := setUp(t)
//
//	mockedName := "test-app"
//	mockedDescription := "Test application description"
//	mockedTeam := "test-team"
//
//	mockedCreateApplicationRequest := &api.CreateApplicationRequest{
//		Name:        mockedName,
//		Description: mockedDescription,
//		Team:        mockedTeam,
//	}
//
//	mockedError := errors.NewInternalServerError("Internal test error")
//
//	mocks.applicationServiceMock.EXPECT().
//		AddApplication(gomock.Any(), mockedName, mockedDescription, mockedTeam).
//		Return(mockedError)
//
//	url := "/applications"
//
//	request, recorder, err := test.NewHTTPRequest("POST", url, mockedCreateApplicationRequest)
//	if err != nil {
//		t.Fatalf("Failed to create request: %v", err)
//	}
//
//	router.ServeHTTP(recorder, request)
//
//	actualResponse := api.ErrorResponse{}
//	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
//	if err != nil {
//		t.Fatalf("Failed to decode response: %v", err)
//	}
//
//	require.Equal(t, http.StatusInternalServerError, recorder.Code)
//	require.Equal(t, "Internal test error", actualResponse.Error)
//}
