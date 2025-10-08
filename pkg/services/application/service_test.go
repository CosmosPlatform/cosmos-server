package application

import (
	"context"
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/storage"
	storageMock "cosmos-server/pkg/storage/mock"
	"cosmos-server/pkg/storage/obj"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddApplication(t *testing.T) {
	t.Run("add application - success", addApplicationSuccess)
	t.Run("add application with git information - success", addApplicationWithGitInformationSuccess)
	t.Run("add application with git information and team - success", addApplicationWithGitInformationAndTeamSuccess)
	t.Run("add application - no team success", addApplicationNoTeamSuccess)
	t.Run("add application - invalid team error", addApplicationInvalidTeamError)
	t.Run("add application - insert application error", addApplicationInsertApplicationError)
}

func TestGetApplication(t *testing.T) {
	t.Run("get application - success", getApplicationSuccess)
	t.Run("get application - not found error", getApplicationNotFoundError)
	t.Run("get application - storage error", getApplicationStorageError)
}

func TestGetApplicationsWithFilter(t *testing.T) {
	t.Run("get applications with filter - success", getApplicationsWithFilterSuccess)
	t.Run("get applications with filter - empty result", getApplicationsWithFilterEmptyResult)
	t.Run("get applications with filter - storage error", getApplicationsWithFilterStorageError)
}

func TestDeleteApplication(t *testing.T) {
	t.Run("delete application - success", deleteApplicationSuccess)
	t.Run("delete application - not found error", deleteApplicationNotFoundError)
	t.Run("delete application - storage error", deleteApplicationStorageError)
}

func TestGetApplicationsByTeam(t *testing.T) {
	t.Run("get applications by team - success", getApplicationsByTeamSuccess)
	t.Run("get applications by team - empty result", getApplicationsByTeamEmptyResult)
	t.Run("get applications by team - team not found error", getApplicationsByTeamNotFoundError)
	t.Run("get applications by team - storage error", getApplicationsByTeamStorageError)
}

func TestUpdateApplication(t *testing.T) {
	t.Run("update application - success", updateApplicationSuccess)
	t.Run("update application with git information - success", updateApplicationWithGitInformationSuccess)
	t.Run("update application with team - success", updateApplicationWithTeamSuccess)
	t.Run("update application remove team - success", updateApplicationRemoveTeamSuccess)
	t.Run("update application - not found error", updateApplicationNotFoundError)
	t.Run("update application - invalid team error", updateApplicationInvalidTeamError)
	t.Run("update application - name conflict error", updateApplicationNameConflictError)
	t.Run("update application - storage error", updateApplicationStorageError)
	t.Run("update application - get updated application error", updateApplicationGetUpdatedError)
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

	mocks.storageServiceMock.EXPECT().
		CheckPendingDependenciesForApplication(gomock.Any(), applicationName).
		Return(nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam, nil)
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

	mocks.storageServiceMock.EXPECT().
		CheckPendingDependenciesForApplication(gomock.Any(), applicationName).
		Return(nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam, nil)
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

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam, nil)
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

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam, nil)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "application with name "+applicationName+" already exists"))
}

func getApplicationSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"

	objApplication := &obj.Application{
		Name:        applicationName,
		Description: applicationDescription,
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(objApplication, nil)

	result, err := applicationService.GetApplication(context.Background(), applicationName)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, applicationName, result.Name)
	require.Equal(t, applicationDescription, result.Description)
}

func getApplicationNotFoundError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "non-existent-application"

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(nil, storage.ErrNotFound)

	result, err := applicationService.GetApplication(context.Background(), applicationName)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "application not found"))
}

func getApplicationStorageError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(nil, storage.ErrInternal)

	result, err := applicationService.GetApplication(context.Background(), applicationName)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "failed to retrieve application"))
}

func getApplicationsWithFilterSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	filter := "test"

	objApplications := []*obj.Application{
		{
			Name:        "test-application-1",
			Description: "first test description",
		},
		{
			Name:        "test-application-2",
			Description: "second test description",
		},
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationsWithFilter(gomock.Any(), filter).
		Return(objApplications, nil)

	result, err := applicationService.GetApplicationsWithFilter(context.Background(), filter)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 2)
	require.Equal(t, "test-application-1", result[0].Name)
	require.Equal(t, "first test description", result[0].Description)
	require.Equal(t, "test-application-2", result[1].Name)
	require.Equal(t, "second test description", result[1].Description)
}

func getApplicationsWithFilterEmptyResult(t *testing.T) {
	applicationService, mocks := setUp(t)

	filter := "nonexistent"

	mocks.storageServiceMock.EXPECT().
		GetApplicationsWithFilter(gomock.Any(), filter).
		Return([]*obj.Application{}, nil)

	result, err := applicationService.GetApplicationsWithFilter(context.Background(), filter)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 0)
}

func getApplicationsWithFilterStorageError(t *testing.T) {
	applicationService, mocks := setUp(t)

	filter := "test"

	mocks.storageServiceMock.EXPECT().
		GetApplicationsWithFilter(gomock.Any(), filter).
		Return(nil, storage.ErrInternal)

	result, err := applicationService.GetApplicationsWithFilter(context.Background(), filter)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "failed to retrieve applications"))
}

func deleteApplicationSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"

	mocks.storageServiceMock.EXPECT().
		DeleteApplicationWithName(gomock.Any(), applicationName).
		Return(nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	err := applicationService.DeleteApplication(context.Background(), applicationName)
	require.NoError(t, err)
}

func deleteApplicationNotFoundError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "non-existent-application"

	mocks.storageServiceMock.EXPECT().
		DeleteApplicationWithName(gomock.Any(), applicationName).
		Return(storage.ErrNotFound)

	err := applicationService.DeleteApplication(context.Background(), applicationName)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "application not found"))
}

func deleteApplicationStorageError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"

	mocks.storageServiceMock.EXPECT().
		DeleteApplicationWithName(gomock.Any(), applicationName).
		Return(storage.ErrInternal)

	err := applicationService.DeleteApplication(context.Background(), applicationName)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "failed to delete application"))
}

func getApplicationsByTeamSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	teamName := "test-team"

	objApplications := []*obj.Application{
		{
			Name:        "test-application-1",
			Description: "first team application",
		},
		{
			Name:        "test-application-2",
			Description: "second team application",
		},
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationsByTeam(gomock.Any(), teamName).
		Return(objApplications, nil)

	result, err := applicationService.GetApplicationsByTeam(context.Background(), teamName)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 2)
	require.Equal(t, "test-application-1", result[0].Name)
	require.Equal(t, "first team application", result[0].Description)
	require.Equal(t, "test-application-2", result[1].Name)
	require.Equal(t, "second team application", result[1].Description)
}

func getApplicationsByTeamEmptyResult(t *testing.T) {
	applicationService, mocks := setUp(t)

	teamName := "empty-team"

	mocks.storageServiceMock.EXPECT().
		GetApplicationsByTeam(gomock.Any(), teamName).
		Return([]*obj.Application{}, nil)

	result, err := applicationService.GetApplicationsByTeam(context.Background(), teamName)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 0)
}

func getApplicationsByTeamNotFoundError(t *testing.T) {
	applicationService, mocks := setUp(t)

	teamName := "non-existent-team"

	mocks.storageServiceMock.EXPECT().
		GetApplicationsByTeam(gomock.Any(), teamName).
		Return(nil, storage.ErrNotFound)

	result, err := applicationService.GetApplicationsByTeam(context.Background(), teamName)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "team not found"))
}

func getApplicationsByTeamStorageError(t *testing.T) {
	applicationService, mocks := setUp(t)

	teamName := "test-team"

	mocks.storageServiceMock.EXPECT().
		GetApplicationsByTeam(gomock.Any(), teamName).
		Return(nil, storage.ErrInternal)

	result, err := applicationService.GetApplicationsByTeam(context.Background(), teamName)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "failed to retrieve applications by team"))
}

func addApplicationWithGitInformationSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	applicationTeam := ""
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	expectedApplicationObj := &obj.Application{
		Name:                applicationName,
		Description:         applicationDescription,
		GitProvider:         "github",
		GitRepositoryOwner:  "test-owner",
		GitRepositoryName:   "test-repo",
		GitRepositoryBranch: "main",
	}

	mocks.storageServiceMock.EXPECT().
		InsertApplication(gomock.Any(), expectedApplicationObj).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		CheckPendingDependenciesForApplication(gomock.Any(), applicationName).
		Return(nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam, gitInformation)
	require.NoError(t, err)
}

func addApplicationWithGitInformationAndTeamSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	applicationTeam := "test-team"
	gitInformation := &model.GitInformation{
		Provider:         "gitlab",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "develop",
	}

	objTeam := &obj.Team{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationTeam,
		Description: "team description",
	}

	teamID := int(objTeam.CosmosObj.ID)

	expectedApplicationObj := &obj.Application{
		Name:                applicationName,
		Description:         applicationDescription,
		TeamID:              &teamID,
		GitProvider:         "gitlab",
		GitRepositoryOwner:  "test-owner",
		GitRepositoryName:   "test-repo",
		GitRepositoryBranch: "develop",
	}

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), applicationTeam).
		Return(objTeam, nil)

	mocks.storageServiceMock.EXPECT().
		InsertApplication(gomock.Any(), expectedApplicationObj).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		CheckPendingDependenciesForApplication(gomock.Any(), applicationName).
		Return(nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	err := applicationService.AddApplication(context.Background(), applicationName, applicationDescription, applicationTeam, gitInformation)
	require.NoError(t, err)
}

func updateApplicationSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	newName := "updated-application"
	newDescription := "updated description"

	updateData := &model.ApplicationUpdate{
		Name:        &newName,
		Description: &newDescription,
	}

	existingApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "original description",
	}

	expectedUpdateObj := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        newName,
		Description: newDescription,
	}

	updatedApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        newName,
		Description: newDescription,
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(existingApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), expectedUpdateObj).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), newName).
		Return(updatedApp, nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, newName, result.Name)
	require.Equal(t, newDescription, result.Description)
}

func updateApplicationWithGitInformationSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	updateData := &model.ApplicationUpdate{
		GitInformation: gitInformation,
	}

	existingApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
	}

	expectedUpdateObj := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:                applicationName,
		Description:         "test description",
		GitProvider:         "github",
		GitRepositoryOwner:  "test-owner",
		GitRepositoryName:   "test-repo",
		GitRepositoryBranch: "main",
	}

	updatedApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:                applicationName,
		Description:         "test description",
		GitProvider:         "github",
		GitRepositoryOwner:  "test-owner",
		GitRepositoryName:   "test-repo",
		GitRepositoryBranch: "main",
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(existingApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), expectedUpdateObj).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(updatedApp, nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, applicationName, result.Name)
}

func updateApplicationWithTeamSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	teamName := "new-team"

	updateData := &model.ApplicationUpdate{
		Team: &teamName,
	}

	existingApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
	}

	teamObj := &obj.Team{
		CosmosObj: obj.CosmosObj{
			ID: 2,
		},
		Name: teamName,
	}

	teamID := int(teamObj.CosmosObj.ID)

	expectedUpdateObj := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
		TeamID:      &teamID,
	}

	updatedApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
		TeamID:      &teamID,
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(existingApp, nil)

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), teamName).
		Return(teamObj, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), expectedUpdateObj).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(updatedApp, nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, applicationName, result.Name)
}

func updateApplicationRemoveTeamSuccess(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	emptyTeam := ""

	updateData := &model.ApplicationUpdate{
		Team: &emptyTeam,
	}

	teamID := 2
	existingApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
		TeamID:      &teamID,
	}

	expectedUpdateObj := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
		TeamID:      nil,
	}

	updatedApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
		TeamID:      nil,
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(existingApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), expectedUpdateObj).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(updatedApp, nil)

	mocks.loggerMocks.EXPECT().
		Infof(gomock.Any(), gomock.Any())

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, applicationName, result.Name)
}

func updateApplicationNotFoundError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "non-existent-application"
	newName := "updated-name"

	updateData := &model.ApplicationUpdate{
		Name: &newName,
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(nil, storage.ErrNotFound)

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "application not found"))
}

func updateApplicationInvalidTeamError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	invalidTeam := "invalid-team"

	updateData := &model.ApplicationUpdate{
		Team: &invalidTeam,
	}

	existingApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(existingApp, nil)

	mocks.storageServiceMock.EXPECT().
		GetTeamWithName(gomock.Any(), invalidTeam).
		Return(nil, storage.ErrNotFound)

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "team not found"))
}

func updateApplicationNameConflictError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	conflictingName := "existing-application"

	updateData := &model.ApplicationUpdate{
		Name: &conflictingName,
	}

	existingApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
	}

	expectedUpdateObj := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        conflictingName,
		Description: "test description",
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(existingApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), expectedUpdateObj).
		Return(storage.ErrAlreadyExists)

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "application with name "+conflictingName+" already exists"))
}

func updateApplicationStorageError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	newName := "updated-name"

	updateData := &model.ApplicationUpdate{
		Name: &newName,
	}

	existingApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
	}

	expectedUpdateObj := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        newName,
		Description: "test description",
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(existingApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), expectedUpdateObj).
		Return(storage.ErrInternal)

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "failed to update application"))
}

func updateApplicationGetUpdatedError(t *testing.T) {
	applicationService, mocks := setUp(t)

	applicationName := "test-application"
	newName := "updated-name"

	updateData := &model.ApplicationUpdate{
		Name: &newName,
	}

	existingApp := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        applicationName,
		Description: "test description",
	}

	expectedUpdateObj := &obj.Application{
		CosmosObj: obj.CosmosObj{
			ID: 1,
		},
		Name:        newName,
		Description: "test description",
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), applicationName).
		Return(existingApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplication(gomock.Any(), expectedUpdateObj).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), newName).
		Return(nil, storage.ErrInternal)

	result, err := applicationService.UpdateApplication(context.Background(), applicationName, updateData)
	require.Error(t, err)
	require.Nil(t, result)
	require.True(t, strings.Contains(err.Error(), "failed to retrieve updated application"))
}
