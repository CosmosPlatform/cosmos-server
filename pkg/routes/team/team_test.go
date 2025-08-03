package team

import (
	"cosmos-server/api"
	"cosmos-server/pkg/errors"
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
	t.Run("success- create team without description", handleCreateTeamWithoutDescription)
	t.Run("failure - team name is required", handleCreateTeamNameRequiredFailure)
}

func TestHandleGetTeams(t *testing.T) {
	t.Run("success - get teams", handleGetTeamsSuccess)
	t.Run("success - get empty teams list", handleGetTeamsEmptySuccess)
	t.Run("failure - get teams with service error", handleGetTeamsServiceError)
}

func TestHandleDeleteTeam(t *testing.T) {
	t.Run("success - delete team", handleDeleteTeamSuccess)
	t.Run("failure - team name is required", handleDeleteTeamNameRequiredFailure)
	t.Run("failure - delete team with service error", handleDeleteTeamServiceError)
}

func TestHandleAddUserToTeam(t *testing.T) {
	t.Run("success - add user to team", handleAddUserToTeamSuccess)
	t.Run("failure - team name is required", handleAddUserToTeamNameRequiredFailure)
	t.Run("failure - invalid request format", handleAddUserToTeamInvalidRequestFailure)
	t.Run("failure - email validation error", handleAddUserToTeamEmailValidationFailure)
	t.Run("failure - add user with service error", handleAddUserToTeamServiceError)
}

func TestHandleRemoveUserFromTeam(t *testing.T) {
	t.Run("success - remove user from team", handleRemoveUserFromTeamSuccess)
	t.Run("failure - team name is required", handleRemoveUserFromTeamNameRequiredFailure)
	t.Run("failure - email query parameter is required", handleRemoveUserFromTeamEmailRequiredFailure)
	t.Run("failure - team not found", handleRemoveUserFromTeamTeamNotFoundFailure)
	t.Run("failure - remove user with service error", handleRemoveUserFromTeamServiceError)
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

func handleGetTeamsSuccess(t *testing.T) {
	router, mocks := setUp(t)

	teams := []*model.Team{
		{
			Name:        "Team 1",
			Description: "First team",
		},
		{
			Name:        "Team 2",
			Description: "Second team",
		},
	}

	expectedResponse := &api.GetTeamsResponse{
		Teams: []*api.Team{
			{
				Name:        "Team 1",
				Description: "First team",
			},
			{
				Name:        "Team 2",
				Description: "Second team",
			},
		},
	}

	mocks.teamServiceMock.EXPECT().
		GetAllTeams(gomock.Any()).
		Return(teams, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.GetTeamsResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, len(expectedResponse.Teams), len(actualResponse.Teams))
	require.Equal(t, expectedResponse.Teams[0].Name, actualResponse.Teams[0].Name)
	require.Equal(t, expectedResponse.Teams[1].Name, actualResponse.Teams[1].Name)
}

func handleGetTeamsEmptySuccess(t *testing.T) {
	router, mocks := setUp(t)

	var teams []*model.Team

	expectedResponse := &api.GetTeamsResponse{
		Teams: []*api.Team{},
	}

	mocks.teamServiceMock.EXPECT().
		GetAllTeams(gomock.Any()).
		Return(teams, nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.GetTeamsResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, len(expectedResponse.Teams), len(actualResponse.Teams))
}

func handleGetTeamsServiceError(t *testing.T) {
	router, mocks := setUp(t)

	expectedError := errors.NewInternalServerError("test error")

	mocks.teamServiceMock.EXPECT().
		GetAllTeams(gomock.Any()).
		Return(nil, expectedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams"

	request, recorder, err := test.NewHTTPRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, expectedError.Error(), actualResponse.Error)
}

func handleDeleteTeamSuccess(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"

	mocks.teamServiceMock.EXPECT().
		DeleteTeam(gomock.Any(), teamName).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams?name=" + teamName

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func handleDeleteTeamNameRequiredFailure(t *testing.T) {
	router, mocks := setUp(t)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams"

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func handleDeleteTeamServiceError(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"
	expectedError := errors.NewInternalServerError("test error")

	mocks.teamServiceMock.EXPECT().
		DeleteTeam(gomock.Any(), teamName).
		Return(expectedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams?name=" + teamName

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, expectedError.Error(), actualResponse.Error)
}

func handleAddUserToTeamSuccess(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"
	userEmail := "test@example.com"

	addUserRequest := &api.AddUserToTeamRequest{
		Email: userEmail,
	}

	mocks.teamServiceMock.EXPECT().
		AddUserToTeam(gomock.Any(), userEmail, teamName).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams/" + teamName + "/members"

	request, recorder, err := test.NewHTTPRequest("POST", url, addUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func handleAddUserToTeamNameRequiredFailure(t *testing.T) {
	router, mocks := setUp(t)

	userEmail := "test@example.com"

	addUserRequest := &api.AddUserToTeamRequest{
		Email: userEmail,
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams//members"

	request, recorder, err := test.NewHTTPRequest("POST", url, addUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func handleAddUserToTeamInvalidRequestFailure(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams/" + teamName + "/members"

	request, recorder, err := test.NewHTTPRequest("POST", url, "invalid json")
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func handleAddUserToTeamEmailValidationFailure(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"

	addUserRequest := &api.AddUserToTeamRequest{
		Email: "invalid-email",
	}

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams/" + teamName + "/members"

	request, recorder, err := test.NewHTTPRequest("POST", url, addUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func handleAddUserToTeamServiceError(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"
	userEmail := "test@example.com"
	expectedError := errors.NewInternalServerError("test error")

	addUserRequest := &api.AddUserToTeamRequest{
		Email: userEmail,
	}

	mocks.teamServiceMock.EXPECT().
		AddUserToTeam(gomock.Any(), userEmail, teamName).
		Return(expectedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams/" + teamName + "/members"

	request, recorder, err := test.NewHTTPRequest("POST", url, addUserRequest)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, expectedError.Error(), actualResponse.Error)
}

func handleRemoveUserFromTeamSuccess(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"
	userEmail := "test@example.com"

	teamModel := &model.Team{
		Name:        teamName,
		Description: "Test description",
	}

	mocks.teamServiceMock.EXPECT().
		GetTeamByName(gomock.Any(), teamName).
		Return(teamModel, nil)

	mocks.teamServiceMock.EXPECT().
		RemoveUserFromTeam(gomock.Any(), userEmail).
		Return(nil)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams/" + teamName + "/members?email=" + userEmail

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusNoContent, recorder.Code)
}

func handleRemoveUserFromTeamNameRequiredFailure(t *testing.T) {
	router, mocks := setUp(t)

	userEmail := "test@example.com"

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams//members?email=" + userEmail

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func handleRemoveUserFromTeamEmailRequiredFailure(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams/" + teamName + "/members"

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func handleRemoveUserFromTeamTeamNotFoundFailure(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Nonexistent Team"
	userEmail := "test@example.com"
	expectedError := errors.NewNotFoundError("team not found")

	mocks.teamServiceMock.EXPECT().
		GetTeamByName(gomock.Any(), teamName).
		Return(nil, expectedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams/" + teamName + "/members?email=" + userEmail

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Equal(t, expectedError.Error(), actualResponse.Error)
}

func handleRemoveUserFromTeamServiceError(t *testing.T) {
	router, mocks := setUp(t)

	teamName := "Test Team"
	userEmail := "test@example.com"
	expectedError := errors.NewInternalServerError("test error")

	teamModel := &model.Team{
		Name:        teamName,
		Description: "Test description",
	}

	mocks.teamServiceMock.EXPECT().
		GetTeamByName(gomock.Any(), teamName).
		Return(teamModel, nil)

	mocks.teamServiceMock.EXPECT().
		RemoveUserFromTeam(gomock.Any(), userEmail).
		Return(expectedError)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

	url := "/teams/" + teamName + "/members?email=" + userEmail

	request, recorder, err := test.NewHTTPRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	actualResponse := &api.ErrorResponse{}
	err = json.NewDecoder(recorder.Body).Decode(actualResponse)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
	require.Equal(t, expectedError.Error(), actualResponse.Error)
}
