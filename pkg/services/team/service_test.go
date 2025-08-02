package team

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/model"
	teamMock "cosmos-server/pkg/services/team/mock"
	"cosmos-server/pkg/storage"
	storageMock "cosmos-server/pkg/storage/mock"
	"cosmos-server/pkg/storage/obj"
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetAllTeams(t *testing.T) {
	t.Run("get all teams - success", getAllTeamsSuccess)
	t.Run("get all teams - storage error", getAllTeamsStorageError)
}

func TestInsertTeam(t *testing.T) {
	t.Run("insert team - success", insertTeamSuccess)
	t.Run("insert team - already exists error", insertTeamAlreadyExistsError)
	t.Run("insert team - storage error", insertTeamStorageError)
}

func TestDeleteTeam(t *testing.T) {
	t.Run("delete team - success", deleteTeamSuccess)
	t.Run("delete team - not found error", deleteTeamNotFoundError)
	t.Run("delete team - storage error", deleteTeamStorageError)
}

func TestAddUserToTeam(t *testing.T) {
	t.Run("add user to team - success", addUserToTeamSuccess)
	t.Run("add user to team - storage error", addUserToTeamStorageError)
}

func TestRemoveUserFromTeam(t *testing.T) {
	t.Run("remove user from team - success", removeUserFromTeamSuccess)
	t.Run("remove user from team - not found error", removeUserFromTeamNotFoundError)
	t.Run("remove user from team - storage error", removeUserFromTeamStorageError)
}

func TestGetTeamByName(t *testing.T) {
	t.Run("get team by name - success", getTeamByNameSuccess)
	t.Run("get team by name - not found error", getTeamByNameNotFoundError)
	t.Run("get team by name - storage error", getTeamByNameStorageError)
}

type mocks struct {
	controller         *gomock.Controller
	storageServiceMock *storageMock.MockService
	translatorMock     *teamMock.MockTranslator
}

func setUp(t *testing.T) (Service, *mocks) {
	ctrl := gomock.NewController(t)
	storageServiceMock := storageMock.NewMockService(ctrl)
	translatorMock := teamMock.NewMockTranslator(ctrl)

	mocks := &mocks{
		controller:         ctrl,
		storageServiceMock: storageServiceMock,
		translatorMock:     translatorMock,
	}

	service := NewTeamService(storageServiceMock, translatorMock)

	return service, mocks
}

func getAllTeamsSuccess(t *testing.T) {
	service, mocks := setUp(t)

	objTeams := []*obj.Team{
		{Name: "Team A", Description: "Description A"},
		{Name: "Team B", Description: "Description B"},
	}

	modelTeams := []*model.Team{
		{Name: "Team A", Description: "Description A"},
		{Name: "Team B", Description: "Description B"},
	}

	mocks.storageServiceMock.EXPECT().
		GetTeamsWithFilter(gomock.Any(), "").
		Return(objTeams, nil)

	mocks.translatorMock.EXPECT().
		ToModelTeams(objTeams).
		Return(modelTeams)

	teams, err := service.GetAllTeams(context.Background())

	require.NoError(t, err)
	require.Len(t, teams, len(modelTeams))
	require.Equal(t, modelTeams[0], teams[0])
	require.Equal(t, modelTeams[1], teams[1])
}

func getAllTeamsStorageError(t *testing.T) {
	service, mocks := setUp(t)

	mockedError := fmt.Errorf("storage error")

	mocks.storageServiceMock.EXPECT().
		GetTeamsWithFilter(gomock.Any(), "").
		Return(nil, mockedError)

	teams, err := service.GetAllTeams(context.Background())

	require.Error(t, err)
	require.Nil(t, teams)
	require.Equal(t, mockedError, err)
	mocks.controller.Finish()
}

func insertTeamSuccess(t *testing.T) {
	service, mocks := setUp(t)

	modelTeam := &model.Team{
		Name:        "Team A",
		Description: "Description A",
	}

	objTeam := &obj.Team{
		Name:        "Team A",
		Description: "Description A",
	}

	mocks.translatorMock.EXPECT().
		ToObjTeam(modelTeam).
		Return(objTeam)

	mocks.storageServiceMock.EXPECT().
		InsertTeam(gomock.Any(), objTeam).
		Return(nil)

	err := service.InsertTeam(context.Background(), modelTeam)

	require.NoError(t, err)
}

func insertTeamAlreadyExistsError(t *testing.T) {
	service, mocks := setUp(t)

	modelTeam := &model.Team{
		Name:        "Team A",
		Description: "Description A",
	}

	objTeam := &obj.Team{
		Name:        "Team A",
		Description: "Description A",
	}

	mocks.translatorMock.EXPECT().
		ToObjTeam(modelTeam).
		Return(objTeam)

	mocks.storageServiceMock.EXPECT().
		InsertTeam(gomock.Any(), objTeam).
		Return(storage.ErrAlreadyExists)

	expectedError := errors.NewConflictError(fmt.Sprint("team with name ", modelTeam.Name, " already exists"))

	err := service.InsertTeam(context.Background(), modelTeam)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func insertTeamStorageError(t *testing.T) {
	service, mocks := setUp(t)

	modelTeam := &model.Team{
		Name:        "Team A",
		Description: "Description A",
	}

	objTeam := &obj.Team{
		Name:        "Team A",
		Description: "Description A",
	}

	mockedError := fmt.Errorf("database connection failed")

	mocks.translatorMock.EXPECT().
		ToObjTeam(modelTeam).
		Return(objTeam)

	mocks.storageServiceMock.EXPECT().
		InsertTeam(gomock.Any(), objTeam).
		Return(mockedError)

	err := service.InsertTeam(context.Background(), modelTeam)

	require.Error(t, err)
	require.Equal(t, mockedError, err)
}

func deleteTeamSuccess(t *testing.T) {
	service, mocks := setUp(t)

	teamName := "Team A"

	mocks.storageServiceMock.EXPECT().
		DeleteTeam(gomock.Any(), teamName).
		Return(nil)

	err := service.DeleteTeam(context.Background(), teamName)

	require.NoError(t, err)
}

func deleteTeamNotFoundError(t *testing.T) {
	service, mocks := setUp(t)

	teamName := "NonExistent Team"

	mocks.storageServiceMock.EXPECT().
		DeleteTeam(gomock.Any(), teamName).
		Return(storage.ErrNotFound)

	expectedError := errors.NewNotFoundError(fmt.Sprintf("team with name %s not found", teamName))

	err := service.DeleteTeam(context.Background(), teamName)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func deleteTeamStorageError(t *testing.T) {
	service, mocks := setUp(t)

	teamName := "Team A"
	mockedError := fmt.Errorf("database connection failed")

	mocks.storageServiceMock.EXPECT().
		DeleteTeam(gomock.Any(), teamName).
		Return(mockedError)

	expectedError := errors.NewInternalServerError(fmt.Sprintf("failed to delete team with name %s: %v", teamName, mockedError))

	err := service.DeleteTeam(context.Background(), teamName)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func addUserToTeamSuccess(t *testing.T) {
	service, mocks := setUp(t)

	userEmail := "user@example.com"
	teamName := "Team A"

	mocks.storageServiceMock.EXPECT().
		AddUserToTeam(gomock.Any(), userEmail, teamName).
		Return(nil)

	err := service.AddUserToTeam(context.Background(), userEmail, teamName)

	require.NoError(t, err)
}

func addUserToTeamStorageError(t *testing.T) {
	service, mocks := setUp(t)

	userEmail := "user@example.com"
	teamName := "Team A"
	mockedError := fmt.Errorf("database connection failed")

	mocks.storageServiceMock.EXPECT().
		AddUserToTeam(gomock.Any(), userEmail, teamName).
		Return(mockedError)

	err := service.AddUserToTeam(context.Background(), userEmail, teamName)

	require.Error(t, err)
	require.Equal(t, mockedError, err)
}

func removeUserFromTeamSuccess(t *testing.T) {
	service, mocks := setUp(t)

	userEmail := "user@example.com"

	mocks.storageServiceMock.EXPECT().
		RemoveUserFromTeam(gomock.Any(), userEmail).
		Return(nil)

	err := service.RemoveUserFromTeam(context.Background(), userEmail)

	require.NoError(t, err)
}

func removeUserFromTeamNotFoundError(t *testing.T) {
	service, mocks := setUp(t)

	userEmail := "nonexistent@example.com"

	mocks.storageServiceMock.EXPECT().
		RemoveUserFromTeam(gomock.Any(), userEmail).
		Return(storage.ErrNotFound)

	expectedError := errors.NewNotFoundError(fmt.Sprintf("user with email %s not found", userEmail))

	err := service.RemoveUserFromTeam(context.Background(), userEmail)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func removeUserFromTeamStorageError(t *testing.T) {
	service, mocks := setUp(t)

	userEmail := "user@example.com"
	mockedError := fmt.Errorf("database connection failed")

	mocks.storageServiceMock.EXPECT().
		RemoveUserFromTeam(gomock.Any(), userEmail).
		Return(mockedError)

	expectedError := errors.NewInternalServerError(fmt.Sprintf("failed to remove user %s from team: %v", userEmail, mockedError))

	err := service.RemoveUserFromTeam(context.Background(), userEmail)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func getTeamByNameSuccess(t *testing.T) {
	service, mocks := setUp(t)

	teamName := "Team A"
	objTeam := &obj.Team{
		Name:        "Team A",
		Description: "Description A",
	}
	modelTeam := &model.Team{
		Name:        "Team A",
		Description: "Description A",
	}

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), teamName).
		Return(objTeam, nil)

	mocks.translatorMock.EXPECT().
		ToModelTeam(objTeam).
		Return(modelTeam)

	team, err := service.GetTeamByName(context.Background(), teamName)

	require.NoError(t, err)
	require.Equal(t, modelTeam, team)
}

func getTeamByNameNotFoundError(t *testing.T) {
	service, mocks := setUp(t)

	teamName := "NonExistent Team"

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), teamName).
		Return(nil, storage.ErrNotFound)

	expectedError := errors.NewNotFoundError(fmt.Sprintf("team with name %s not found", teamName))

	team, err := service.GetTeamByName(context.Background(), teamName)

	require.Error(t, err)
	require.Nil(t, team)
	require.Equal(t, expectedError, err)
}

func getTeamByNameStorageError(t *testing.T) {
	service, mocks := setUp(t)

	teamName := "Team A"
	mockedError := fmt.Errorf("database connection failed")

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), teamName).
		Return(nil, mockedError)

	expectedError := errors.NewInternalServerError(fmt.Sprintf("failed to get team with name %s: %v", teamName, mockedError))

	team, err := service.GetTeamByName(context.Background(), teamName)

	require.Error(t, err)
	require.Nil(t, team)
	require.Equal(t, expectedError, err)
}
