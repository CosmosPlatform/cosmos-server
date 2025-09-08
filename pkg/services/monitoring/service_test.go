package monitoring

import (
	"context"
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/services/monitoring/mock"
	"cosmos-server/pkg/storage"
	storageMock "cosmos-server/pkg/storage/mock"
	"cosmos-server/pkg/storage/obj"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateApplicationInformation(t *testing.T) {
	t.Run("update application information - success", updateApplicationInformationSuccess)
	t.Run("update application information - no git information", updateApplicationInformationNoGitInformation)
	t.Run("update application information - get file error", updateApplicationInformationGetFileError)
	t.Run("update application information - invalid json", updateApplicationInformationInvalidJSON)
	t.Run("update application information - invalid specification", updateApplicationInformationInvalidSpecification)
	t.Run("update application information - provider not found error", updateApplicationInformationProviderNotFoundError)
	t.Run("update application information - upsert dependency error", updateApplicationInformationUpsertDependencyError)
	t.Run("update application information - delete obsolete dependencies error", updateApplicationInformationDeleteObsoleteDependenciesError)
	t.Run("update application information - with obsolete dependencies", updateApplicationInformationWithObsoleteDependencies)
}

type mocks struct {
	controller         *gomock.Controller
	gitServiceMock     *mock.MockGitService
	storageServiceMock *storageMock.MockService
	loggerMocks        *log.MockLogger
}

func setUp(t *testing.T) (Service, *mocks) {
	controller := gomock.NewController(t)

	mocks := &mocks{
		controller:         controller,
		gitServiceMock:     mock.NewMockGitService(controller),
		storageServiceMock: storageMock.NewMockService(controller),
		loggerMocks:        log.NewMockLogger(controller),
	}

	service := NewMonitoringService(mocks.storageServiceMock, mocks.gitServiceMock, NewTranslator(), mocks.loggerMocks)

	return service, mocks
}

func updateApplicationInformationSuccess(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: gitInformation,
	}

	openClientSpecification := getMockedOpenClientSpecification()
	jsonContent, err := json.Marshal(openClientSpecification)
	if err != nil {
		t.Fatalf("Failed to marshal open client specification: %v", err)
	}

	fileContent := &model.FileContent{
		Metadata: model.FileMetadata{
			Name:       "openclient.json",
			Path:       "docs/openclient.json",
			Size:       len(jsonContent),
			SHA:        "abc123",
			Branch:     gitInformation.RepositoryBranch,
			Repository: gitInformation.RepositoryName,
			Owner:      gitInformation.RepositoryOwner,
		},
		Content: string(jsonContent),
	}

	objApplicationDependency := getObjOpenClientSpecification()

	providerApp := &obj.Application{
		CosmosObj:   obj.CosmosObj{ID: 1},
		Name:        "service-a",
		Description: "Provider service",
	}

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), providerApp.Name).
		Return(providerApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpsertApplicationDependency(gomock.Any(), modelApplication.Name, providerApp.Name, objApplicationDependency).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesByConsumer(gomock.Any(), modelApplication.Name).
		Return([]*obj.ApplicationDependency{}, nil)

	mocks.loggerMocks.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any())

	err = service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.NoError(t, err)
}

func getMockedOpenClientSpecification() *model.OpenClientSpecification {
	return &model.OpenClientSpecification{
		Dependencies: map[string]model.DependencySpecification{
			"service-a": {
				Reasons: []string{"reason1", "reason2"},
				Endpoints: map[string]model.EndpointMethodsSpecification{
					"/users": map[string]model.EndpointSpecification{
						"GET": {
							Reasons: []string{"fetch users"},
						},
						"POST": {
							Reasons: []string{"create user"},
						},
					},
				},
			},
		},
	}
}

func getObjOpenClientSpecification() *obj.ApplicationDependency {
	return &obj.ApplicationDependency{
		Reasons: []string{"reason1", "reason2"},
		Endpoints: obj.Endpoints{
			"/users": obj.EndpointMethods{
				"GET": obj.EndpointDetails{
					Reasons: []string{"fetch users"},
				},
				"POST": obj.EndpointDetails{
					Reasons: []string{"create user"},
				},
			},
		},
	}
}

func updateApplicationInformationNoGitInformation(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: nil,
	}

	mocks.loggerMocks.EXPECT().
		Infof("No git information for application %s, skipping monitoring update", applicationName)

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.NoError(t, err)
}

func updateApplicationInformationGetFileError(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: gitInformation,
	}

	gitError := errors.New("repository not found")

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(nil, gitError)

	mocks.loggerMocks.EXPECT().
		Errorf("Failed to get openclient.json for application %s: %v", applicationName, gitError)

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.Error(t, err)
	require.Equal(t, gitError, err)
}

func updateApplicationInformationInvalidJSON(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: gitInformation,
	}

	invalidJSONContent := "{ invalid json content"

	fileContent := &model.FileContent{
		Metadata: model.FileMetadata{
			Name:       "openclient.json",
			Path:       "docs/openclient.json",
			Size:       len(invalidJSONContent),
			SHA:        "abc123",
			Branch:     gitInformation.RepositoryBranch,
			Repository: gitInformation.RepositoryName,
			Owner:      gitInformation.RepositoryOwner,
		},
		Content: invalidJSONContent,
	}

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.loggerMocks.EXPECT().
		Errorf("Failed to unmarshal openclient.json for application %s: %v", applicationName, gomock.Any())

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "invalid character"))
}

func updateApplicationInformationInvalidSpecification(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: gitInformation,
	}

	invalidSpecification := getMockedInvalidOpenClientSpecification()

	jsonContent, _ := json.Marshal(invalidSpecification)

	fileContent := &model.FileContent{
		Metadata: model.FileMetadata{
			Name:       "openclient.json",
			Path:       "docs/openclient.json",
			Size:       len(jsonContent),
			SHA:        "abc123",
			Branch:     gitInformation.RepositoryBranch,
			Repository: gitInformation.RepositoryName,
			Owner:      gitInformation.RepositoryOwner,
		},
		Content: string(jsonContent),
	}

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.loggerMocks.EXPECT().
		Errorf("Invalid openclient.json for application %s: %v", applicationName, gomock.Any())

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.Error(t, err)
}

func getMockedInvalidOpenClientSpecification() *model.OpenClientSpecification {
	return &model.OpenClientSpecification{
		Dependencies: map[string]model.DependencySpecification{
			"service-a": {
				Reasons: []string{"reason1", "reason2"},
				Endpoints: map[string]model.EndpointMethodsSpecification{
					"/users that I want to bring": map[string]model.EndpointSpecification{
						"GET": {
							Reasons: []string{"fetch users"},
						},
						"POST": {
							Reasons: []string{"create user"},
						},
					},
				},
			},
		},
	}
}

func updateApplicationInformationProviderNotFoundError(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: gitInformation,
	}

	openClientSpecification := getMockedOpenClientSpecification()
	jsonContent, _ := json.Marshal(openClientSpecification)

	fileContent := &model.FileContent{
		Metadata: model.FileMetadata{
			Name:       "openclient.json",
			Path:       "docs/openclient.json",
			Size:       len(jsonContent),
			SHA:        "abc123",
			Branch:     gitInformation.RepositoryBranch,
			Repository: gitInformation.RepositoryName,
			Owner:      gitInformation.RepositoryOwner,
		},
		Content: string(jsonContent),
	}

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), "service-a").
		Return(nil, storage.ErrNotFound)

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesByConsumer(gomock.Any(), modelApplication.Name).
		Return([]*obj.ApplicationDependency{}, nil)

	mocks.loggerMocks.EXPECT().
		Errorf("Failed to transform dependency for application %s: %v", applicationName, gomock.Any())

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.NoError(t, err) // Should not fail overall, just log the error
}

func updateApplicationInformationUpsertDependencyError(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: gitInformation,
	}

	openClientSpecification := getMockedOpenClientSpecification()
	jsonContent, _ := json.Marshal(openClientSpecification)

	fileContent := &model.FileContent{
		Metadata: model.FileMetadata{
			Name:       "openclient.json",
			Path:       "docs/openclient.json",
			Size:       len(jsonContent),
			SHA:        "abc123",
			Branch:     gitInformation.RepositoryBranch,
			Repository: gitInformation.RepositoryName,
			Owner:      gitInformation.RepositoryOwner,
		},
		Content: string(jsonContent),
	}

	objApplicationDependency := getObjOpenClientSpecification()

	providerApp := &obj.Application{
		CosmosObj:   obj.CosmosObj{ID: 1},
		Name:        "service-a",
		Description: "Provider service",
	}

	upsertError := errors.New("database connection failed")

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), providerApp.Name).
		Return(providerApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpsertApplicationDependency(gomock.Any(), modelApplication.Name, providerApp.Name, objApplicationDependency).
		Return(upsertError)

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesByConsumer(gomock.Any(), modelApplication.Name).
		Return([]*obj.ApplicationDependency{}, nil)

	mocks.loggerMocks.EXPECT().
		Errorf("Failed to upsert dependency for application %s: %v", applicationName, upsertError)

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.NoError(t, err) // Should not fail overall, just log the error
}

func updateApplicationInformationDeleteObsoleteDependenciesError(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: gitInformation,
	}

	openClientSpecification := getMockedOpenClientSpecification()
	jsonContent, _ := json.Marshal(openClientSpecification)

	fileContent := &model.FileContent{
		Metadata: model.FileMetadata{
			Name:       "openclient.json",
			Path:       "docs/openclient.json",
			Size:       len(jsonContent),
			SHA:        "abc123",
			Branch:     gitInformation.RepositoryBranch,
			Repository: gitInformation.RepositoryName,
			Owner:      gitInformation.RepositoryOwner,
		},
		Content: string(jsonContent),
	}

	objApplicationDependency := getObjOpenClientSpecification()

	providerApp := &obj.Application{
		CosmosObj:   obj.CosmosObj{ID: 1},
		Name:        "service-a",
		Description: "Provider service",
	}

	getDependenciesError := errors.New("failed to get dependencies")

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), providerApp.Name).
		Return(providerApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpsertApplicationDependency(gomock.Any(), modelApplication.Name, providerApp.Name, objApplicationDependency).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesByConsumer(gomock.Any(), modelApplication.Name).
		Return(nil, getDependenciesError)

	mocks.loggerMocks.EXPECT().
		Infof("Successfully upserted dependency from %s to %s", applicationName, providerApp.Name)

	mocks.loggerMocks.EXPECT().
		Errorf("Failed to delete obsolete dependencies for application %s: %v", applicationName, getDependenciesError)

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.Error(t, err)
	require.Equal(t, getDependenciesError, err)
}

func updateApplicationInformationWithObsoleteDependencies(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	applicationDescription := "test-description"
	gitInformation := &model.GitInformation{
		Provider:         "github",
		RepositoryOwner:  "test-owner",
		RepositoryName:   "test-repo",
		RepositoryBranch: "main",
	}

	modelApplication := &model.Application{
		Name:           applicationName,
		Description:    applicationDescription,
		GitInformation: gitInformation,
	}

	openClientSpecification := getMockedOpenClientSpecification()
	jsonContent, _ := json.Marshal(openClientSpecification)

	fileContent := &model.FileContent{
		Metadata: model.FileMetadata{
			Name:       "openclient.json",
			Path:       "docs/openclient.json",
			Size:       len(jsonContent),
			SHA:        "abc123",
			Branch:     gitInformation.RepositoryBranch,
			Repository: gitInformation.RepositoryName,
			Owner:      gitInformation.RepositoryOwner,
		},
		Content: string(jsonContent),
	}

	objApplicationDependency := getObjOpenClientSpecification()

	providerApp := &obj.Application{
		CosmosObj:   obj.CosmosObj{ID: 1},
		Name:        "service-a",
		Description: "Provider service",
	}

	obsoleteProviderApp := &obj.Application{
		CosmosObj:   obj.CosmosObj{ID: 2},
		Name:        "service-b",
		Description: "Obsolete service",
	}

	existingDependencies := []*obj.ApplicationDependency{
		{
			Consumer: &obj.Application{Name: applicationName},
			Provider: obsoleteProviderApp,
		},
	}

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), providerApp.Name).
		Return(providerApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpsertApplicationDependency(gomock.Any(), modelApplication.Name, providerApp.Name, objApplicationDependency).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesByConsumer(gomock.Any(), modelApplication.Name).
		Return(existingDependencies, nil)

	mocks.storageServiceMock.EXPECT().
		DeleteApplicationDependency(gomock.Any(), applicationName, "service-b").
		Return(nil)

	mocks.loggerMocks.EXPECT().
		Infof("Successfully upserted dependency from %s to %s", applicationName, providerApp.Name)

	mocks.loggerMocks.EXPECT().
		Infof("Successfully deleted obsolete dependency from %s to %s", applicationName, "service-b")

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.NoError(t, err)
}
