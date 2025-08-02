package team

import (
	"cosmos-server/api"
	logMock "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/test"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"

	teamMock "cosmos-server/pkg/services/team/mock"
)

func TestHandleCreateTeam(t *testing.T) {
	t.Run("success - create team", handleCreateTeamSuccess)
	t.Run("failure - team name is required", handleCreateTeamNameRequiredFailure)
	t.Run("success- create team without description", handleCreateTeamWithoutDescription)
}

type mocks struct {
	controller      *gomock.Controller
	teamServiceMock *teamMock.MockService
	loggerMock      *logMock.MockLogger
}

func setUp(t *testing.T) (*gin.Engine, *mocks) {
	ctrl := gomock.NewController(t)
	teamServiceMock := teamMock.NewMockService(ctrl)
	loggerMock := logMock.NewMockLogger(ctrl)

	mocks := &mocks{
		controller:      ctrl,
		teamServiceMock: teamServiceMock,
		loggerMock:      loggerMock,
	}

	router := test.NewRouter(loggerMock)
	AddAdminTeamHandler(router.Group("/"), teamServiceMock, NewTranslator())

	return router, mocks
}

func handleCreateTeamSuccess(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"
	teamDescription := "This is a test team"

	mockedTeamRequest := &api.CreateTeamRequest{
		Name:        teamName,
		Description: teamDescription,
	}

	teamModel := &model.Team{
		Name:        mockedTeamRequest.Name,
		Description: mockedTeamRequest.Description,
	}

	mockedCreateTeamResponse := &api.CreateTeamResponse{
		Team: &api.Team{
			Name:        mockedTeamRequest.Name,
			Description: mockedTeamRequest.Description,
		}}

	mocks.teamServiceMock.EXPECT().
		InsertTeam(gomock.Any(), teamModel).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedTeamRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.CreateTeamResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, mockedCreateTeamResponse.Team.Name, actualResponse.Team.Name)
}

func handleCreateTeamNameRequiredFailure(t *testing.T) {
	router, mocks := setUp(t)

	teamDescription := "This is a test team"

	mockedTeamRequest := &api.CreateTeamRequest{
		Description: teamDescription,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedTeamRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func handleCreateTeamWithoutDescription(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"

	mockedTeamRequest := &api.CreateTeamRequest{
		Name: teamName,
	}

	teamModel := &model.Team{
		Name:        mockedTeamRequest.Name,
		Description: mockedTeamRequest.Description,
	}

	mockedCreateTeamResponse := &api.CreateTeamResponse{
		Team: &api.Team{
			Name:        mockedTeamRequest.Name,
			Description: mockedTeamRequest.Description,
		}}

	mocks.teamServiceMock.EXPECT().
		InsertTeam(gomock.Any(), teamModel).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams"

	request, recorder, err := test.NewHTTPRequest("POST", url, mockedTeamRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.CreateTeamResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusCreated, recorder.Code)
	require.Equal(t, mockedCreateTeamResponse.Team.Name, actualResponse.Team.Name)
}
