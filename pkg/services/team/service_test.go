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
