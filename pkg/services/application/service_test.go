package application

import (
	"context"
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/storage"
	storageMock "cosmos-server/pkg/storage/mock"
	"cosmos-server/pkg/storage/obj"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"strings"
	"testing"
)

func TestAddApplication(t *testing.T) {
	t.Run("add application - success", addApplicationSuccess)
	t.Run("add application - no team success", addApplicationNoTeamSuccess)
	t.Run("add application - invalid team error", addApplicationInvalidTeamError)
	t.Run("add application - insert application error", addApplicationInsertApplicationError)
}

type mocks struct {
	controller         *gomock.Controller
	storageServiceMock *storageMock.MockService
	loggerMocks        *log.MockLogger
}

func setUp(t *testing.T) (Service, *mocks) {
	ctrl := gomock.NewController(t)

	mocks := &mocks{
		controller:         ctrl,
		storageServiceMock: storageMock.NewMockService(ctrl),
		loggerMocks:        log.NewMockLogger(ctrl),
	}

	applicationService := NewApplicationService(mocks.storageServiceMock, NewTranslator(), mocks.loggerMocks)
	return applicationService, mocks
}

func addApplicationSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	applicationTeam := "test-team"

	objTeam := &obj.Team{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: applicationDescription,
	}

	teamID := int(objTeam.CosmosObj.ID)

	applicationObj := &obj.Application{
		Name:        applicationName,
		Description: applicationDescription,
		TeamID:      &teamID,
	}

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), applicationTeam).
		Return(objTeam, nil)

	mocks.storageServiceMock.EXPECT().
		InsertApplication(gomock.Any(), applicationObj).
		Return(nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam)
	require.NoError(t, err)
}

func addApplicationNoTeamSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	applicationTeam := ""

	applicationObj := &obj.Application{
		Name:        applicationName,
		Description: applicationDescription,
	}

	mocks.storageServiceMock.EXPECT().
		InsertApplication(gomock.Any(), applicationObj).
		Return(nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam)
	require.NoError(t, err)
}

func addApplicationInvalidTeamError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	applicationTeam := "invalid-team"

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), applicationTeam).
		Return(nil, storage.ErrNotFound)

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "team not found"))
}

func addApplicationInsertApplicationError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	applicationTeam := "test-team"

	objTeam := &obj.Team{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: applicationDescription,
	}

	teamID := int(objTeam.CosmosObj.ID)

	applicationObj := &obj.Application{
		Name:        applicationName,
		Description: applicationDescription,
		TeamID:      &teamID,
	}

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), applicationTeam).
		Return(objTeam, nil)

	mocks.storageServiceMock.EXPECT().
		InsertApplication(gomock.Any(), applicationObj).
		Return(storage.ErrAlreadyExists)

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "application with name "+applicationName+" already exists"))
}
