package monitoring

import (
	"context"
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/services/monitoring/mock"
	"cosmos-server/pkg/storage"
	"strings"

	//"cosmos-server/pkg/storage"
	storageMock "cosmos-server/pkg/storage/mock"
	"cosmos-server/pkg/storage/obj"
	"encoding/json"
	"errors"
	//"strings"
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
	t.Run("update application information - provider not found", updateApplicationInformationProviderNotFound)
	t.Run("update application information - upsert dependency error", updateApplicationInformationUpsertDependencyError)
}

func TestGetApplicationInteractions(t *testing.T) {
	t.Run("get application interactions - success", getApplicationInteractionsSuccess)
	t.Run("get application interactions - storage error", getApplicationInteractionsStorageError)
	t.Run("get application interactions - empty dependencies", getApplicationInteractionsEmptyDependencies)
}

func TestGetApplicationsInteractions(t *testing.T) {
	t.Run("get applications interactions - success", getApplicationsInteractionsSuccess)
	t.Run("get applications interactions - storage error", getApplicationsInteractionsStorageError)
	t.Run("get applications interactions - empty dependencies", getApplicationsInteractionsEmptyDependencies)
	t.Run("get applications interactions - with filter", getApplicationsInteractionsWithFilter)
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

	service := NewMonitoringService(mocks.storageServiceMock, mocks.gitServiceMock, NewOpenApiService(), NewTranslator(), mocks.loggerMocks)

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

	sha := "abc123"

	metadata := &model.FileMetadata{
		Name:       "openclient.json",
		Path:       "docs/openclient.json",
		Size:       len(jsonContent),
		SHA:        sha,
		Branch:     gitInformation.RepositoryBranch,
		Repository: gitInformation.RepositoryName,
		Owner:      gitInformation.RepositoryOwner,
	}

	fileContent := &model.FileContent{
		Metadata: *metadata,
		Content:  string(jsonContent),
	}

	objApplicationDependency := getObjOpenClientSpecification()

	providerApp := &obj.Application{
		CosmosObj:   obj.CosmosObj{ID: 1},
		Name:        "service-a",
		Description: "Provider service",
	}

	mockedDependenciesToUpsert := map[string]*obj.ApplicationDependency{
		"service-a": objApplicationDependency,
	}

	mockedPendingDependencies := make(map[string]*obj.PendingApplicationDependency)

	mockedDependenciesToDelete := make([]*obj.ApplicationDependency, 0)

	mocks.gitServiceMock.EXPECT().
		GetFileMetadata(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(metadata, nil)

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), providerApp.Name).
		Return(providerApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplicationDependencies(gomock.Any(), modelApplication.Name, mockedDependenciesToUpsert, mockedPendingDependencies, mockedDependenciesToDelete, sha).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesByConsumer(gomock.Any(), modelApplication.Name).
		Return([]*obj.ApplicationDependency{}, nil)

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

func getObjPendingDependency(name string) *obj.PendingApplicationDependency {
	return &obj.PendingApplicationDependency{
		ProviderName: name,
		Reasons:      []string{"reason1", "reason2"},
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
		GetFileMetadata(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(nil, gitError)

	err := service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.Error(t, err)
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

	sha := "abc123"

	metadata := &model.FileMetadata{
		Name:       "openclient.json",
		Path:       "docs/openclient.json",
		Size:       len(invalidJSONContent),
		SHA:        sha,
		Branch:     gitInformation.RepositoryBranch,
		Repository: gitInformation.RepositoryName,
		Owner:      gitInformation.RepositoryOwner,
	}

	fileContent := &model.FileContent{
		Metadata: *metadata,
		Content:  invalidJSONContent,
	}

	mocks.gitServiceMock.EXPECT().
		GetFileMetadata(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(metadata, nil)

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

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

	sha := "abc123"

	metadata := &model.FileMetadata{
		Name:       "openclient.json",
		Path:       "docs/openclient.json",
		Size:       len(jsonContent),
		SHA:        sha,
		Branch:     gitInformation.RepositoryBranch,
		Repository: gitInformation.RepositoryName,
		Owner:      gitInformation.RepositoryOwner,
	}

	fileContent := &model.FileContent{
		Metadata: *metadata,
		Content:  string(jsonContent),
	}

	mocks.gitServiceMock.EXPECT().
		GetFileMetadata(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(metadata, nil)

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

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

func updateApplicationInformationProviderNotFound(t *testing.T) {
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

	sha := "abc123"

	metadata := &model.FileMetadata{
		Name:       "openclient.json",
		Path:       "docs/openclient.json",
		Size:       len(jsonContent),
		SHA:        sha,
		Branch:     gitInformation.RepositoryBranch,
		Repository: gitInformation.RepositoryName,
		Owner:      gitInformation.RepositoryOwner,
	}

	fileContent := &model.FileContent{
		Metadata: *metadata,
		Content:  string(jsonContent),
	}

	objPendingDependency := getObjPendingDependency("service-a")

	mockedDependenciesToUpsert := make(map[string]*obj.ApplicationDependency)
	mockedPendingDependencies := map[string]*obj.PendingApplicationDependency{
		objPendingDependency.ProviderName: objPendingDependency,
	}

	mockedDependenciesToDelete := make([]*obj.ApplicationDependency, 0)

	mocks.gitServiceMock.EXPECT().
		GetFileMetadata(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(metadata, nil)

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), "service-a").
		Return(nil, storage.ErrNotFound)

	mocks.storageServiceMock.EXPECT().
		UpdateApplicationDependencies(gomock.Any(), modelApplication.Name, mockedDependenciesToUpsert, mockedPendingDependencies, mockedDependenciesToDelete, sha).
		Return(nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesByConsumer(gomock.Any(), modelApplication.Name).
		Return([]*obj.ApplicationDependency{}, nil)

	mocks.loggerMocks.EXPECT().
		Warnf(gomock.Any(), gomock.Any(), gomock.Any())

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
	jsonContent, err := json.Marshal(openClientSpecification)
	if err != nil {
		t.Fatalf("Failed to marshal open client specification: %v", err)
	}

	sha := "abc123"

	metadata := &model.FileMetadata{
		Name:       "openclient.json",
		Path:       "docs/openclient.json",
		Size:       len(jsonContent),
		SHA:        sha,
		Branch:     gitInformation.RepositoryBranch,
		Repository: gitInformation.RepositoryName,
		Owner:      gitInformation.RepositoryOwner,
	}

	fileContent := &model.FileContent{
		Metadata: *metadata,
		Content:  string(jsonContent),
	}

	objApplicationDependency := getObjOpenClientSpecification()

	providerApp := &obj.Application{
		CosmosObj:   obj.CosmosObj{ID: 1},
		Name:        "service-a",
		Description: "Provider service",
	}

	mockedDependenciesToUpsert := map[string]*obj.ApplicationDependency{
		"service-a": objApplicationDependency,
	}

	mockedPendingDependencies := make(map[string]*obj.PendingApplicationDependency)

	mockedDependenciesToDelete := make([]*obj.ApplicationDependency, 0)

	upsertError := errors.New("database connection failed")

	mocks.gitServiceMock.EXPECT().
		GetFileMetadata(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(metadata, nil)

	mocks.gitServiceMock.EXPECT().
		GetFileWithContent(gomock.Any(), gitInformation.RepositoryOwner, gitInformation.RepositoryName, gitInformation.RepositoryBranch, "docs/openclient.json").
		Return(fileContent, nil)

	mocks.storageServiceMock.EXPECT().
		GetApplicationWithName(gomock.Any(), providerApp.Name).
		Return(providerApp, nil)

	mocks.storageServiceMock.EXPECT().
		UpdateApplicationDependencies(gomock.Any(), modelApplication.Name, mockedDependenciesToUpsert, mockedPendingDependencies, mockedDependenciesToDelete, sha).
		Return(upsertError)

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesByConsumer(gomock.Any(), modelApplication.Name).
		Return([]*obj.ApplicationDependency{}, nil)

	err = service.UpdateApplicationInformation(context.TODO(), modelApplication)
	require.Error(t, err)
}

func getApplicationInteractionsSuccess(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"

	objDependencies := []*obj.ApplicationDependency{
		{
			Consumer: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 1},
				Name:        "test-application",
				Description: "Consumer app",
			},
			Provider: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 2},
				Name:        "service-a",
				Description: "Provider app",
			},
			Reasons: []string{"data processing"},
			Endpoints: obj.Endpoints{
				"/api/users": obj.EndpointMethods{
					"GET": obj.EndpointDetails{
						Reasons: []string{"fetch users"},
					},
				},
			},
		},
		{
			Consumer: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 3},
				Name:        "service-b",
				Description: "Another consumer",
			},
			Provider: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 1},
				Name:        "test-application",
				Description: "Provider app",
			},
			Reasons: []string{"authentication"},
			Endpoints: obj.Endpoints{
				"/api/auth": obj.EndpointMethods{
					"POST": obj.EndpointDetails{
						Reasons: []string{"authenticate user"},
					},
				},
			},
		},
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesWithApplicationInvolved(gomock.Any(), applicationName).
		Return(objDependencies, nil)

	result, err := service.GetApplicationInteractions(context.TODO(), applicationName)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Interactions, 2)
	require.Len(t, result.ApplicationsInvolved, 3)

	// Verify applications involved
	require.Contains(t, result.ApplicationsInvolved, "test-application")
	require.Contains(t, result.ApplicationsInvolved, "service-a")
	require.Contains(t, result.ApplicationsInvolved, "service-b")

	// Verify interactions
	require.Equal(t, "test-application", result.Interactions[0].Consumer.Name)
	require.Equal(t, "service-a", result.Interactions[0].Provider.Name)
	require.Equal(t, []string{"data processing"}, result.Interactions[0].Reasons)

	require.Equal(t, "service-b", result.Interactions[1].Consumer.Name)
	require.Equal(t, "test-application", result.Interactions[1].Provider.Name)
	require.Equal(t, []string{"authentication"}, result.Interactions[1].Reasons)
}

func getApplicationInteractionsStorageError(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"
	storageError := errors.New("database connection failed")

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesWithApplicationInvolved(gomock.Any(), applicationName).
		Return(nil, storageError)

	result, err := service.GetApplicationInteractions(context.TODO(), applicationName)
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, storageError, err)
}

func getApplicationInteractionsEmptyDependencies(t *testing.T) {
	service, mocks := setUp(t)

	applicationName := "test-application"

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesWithApplicationInvolved(gomock.Any(), applicationName).
		Return([]*obj.ApplicationDependency{}, nil)

	result, err := service.GetApplicationInteractions(context.TODO(), applicationName)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Interactions, 0)
	require.Len(t, result.ApplicationsInvolved, 0)
}

func getApplicationsInteractionsSuccess(t *testing.T) {
	service, mocks := setUp(t)

	filter := model.ApplicationDependencyFilter{
		Teams:            []string{"team-a", "team-b"},
		IncludeNeighbors: true,
	}

	objDependencies := []*obj.ApplicationDependency{
		{
			Consumer: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 1},
				Name:        "app-1",
				Description: "Consumer app 1",
			},
			Provider: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 2},
				Name:        "app-2",
				Description: "Provider app 2",
			},
			Reasons: []string{"data processing"},
			Endpoints: obj.Endpoints{
				"/api/data": obj.EndpointMethods{
					"GET": obj.EndpointDetails{
						Reasons: []string{"fetch data"},
					},
				},
			},
		},
		{
			Consumer: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 3},
				Name:        "app-3",
				Description: "Consumer app 3",
			},
			Provider: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 1},
				Name:        "app-1",
				Description: "Provider app 1",
			},
			Reasons: []string{"authentication"},
			Endpoints: obj.Endpoints{
				"/api/auth": obj.EndpointMethods{
					"POST": obj.EndpointDetails{
						Reasons: []string{"authenticate"},
					},
				},
			},
		},
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesWithFilter(gomock.Any(), filter).
		Return(objDependencies, nil)

	result, err := service.GetApplicationsInteractions(context.TODO(), filter)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Interactions, 2)
	require.Len(t, result.ApplicationsInvolved, 3)

	// Verify applications involved
	require.Contains(t, result.ApplicationsInvolved, "app-1")
	require.Contains(t, result.ApplicationsInvolved, "app-2")
	require.Contains(t, result.ApplicationsInvolved, "app-3")

	// Verify interactions
	require.Equal(t, "app-1", result.Interactions[0].Consumer.Name)
	require.Equal(t, "app-2", result.Interactions[0].Provider.Name)
	require.Equal(t, []string{"data processing"}, result.Interactions[0].Reasons)

	require.Equal(t, "app-3", result.Interactions[1].Consumer.Name)
	require.Equal(t, "app-1", result.Interactions[1].Provider.Name)
	require.Equal(t, []string{"authentication"}, result.Interactions[1].Reasons)
}

func getApplicationsInteractionsStorageError(t *testing.T) {
	service, mocks := setUp(t)

	filter := model.ApplicationDependencyFilter{
		Teams:            []string{"team-a"},
		IncludeNeighbors: false,
	}

	storageError := errors.New("database connection failed")

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesWithFilter(gomock.Any(), filter).
		Return(nil, storageError)

	result, err := service.GetApplicationsInteractions(context.TODO(), filter)
	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, storageError, err)
}

func getApplicationsInteractionsEmptyDependencies(t *testing.T) {
	service, mocks := setUp(t)

	filter := model.ApplicationDependencyFilter{
		Teams:            []string{"team-a"},
		IncludeNeighbors: true,
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesWithFilter(gomock.Any(), filter).
		Return([]*obj.ApplicationDependency{}, nil)

	result, err := service.GetApplicationsInteractions(context.TODO(), filter)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Interactions, 0)
	require.Len(t, result.ApplicationsInvolved, 0)
}

func getApplicationsInteractionsWithFilter(t *testing.T) {
	service, mocks := setUp(t)

	filter := model.ApplicationDependencyFilter{
		Teams:            []string{"frontend-team"},
		IncludeNeighbors: false,
	}

	objDependencies := []*obj.ApplicationDependency{
		{
			Consumer: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 1},
				Name:        "frontend-app",
				Description: "Frontend application",
			},
			Provider: &obj.Application{
				CosmosObj:   obj.CosmosObj{ID: 2},
				Name:        "backend-service",
				Description: "Backend service",
			},
			Reasons: []string{"API calls"},
			Endpoints: obj.Endpoints{
				"/api/users": obj.EndpointMethods{
					"GET": obj.EndpointDetails{
						Reasons: []string{"fetch user data"},
					},
					"POST": obj.EndpointDetails{
						Reasons: []string{"create user"},
					},
				},
			},
		},
	}

	mocks.storageServiceMock.EXPECT().
		GetApplicationDependenciesWithFilter(gomock.Any(), filter).
		Return(objDependencies, nil)

	result, err := service.GetApplicationsInteractions(context.TODO(), filter)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Interactions, 1)
	require.Len(t, result.ApplicationsInvolved, 2)

	// Verify applications involved
	require.Contains(t, result.ApplicationsInvolved, "frontend-app")
	require.Contains(t, result.ApplicationsInvolved, "backend-service")

	// Verify interaction details
	interaction := result.Interactions[0]
	require.Equal(t, "frontend-app", interaction.Consumer.Name)
	require.Equal(t, "backend-service", interaction.Provider.Name)
	require.Equal(t, []string{"API calls"}, interaction.Reasons)

	// Verify endpoints
	require.Contains(t, interaction.Endpoints, "/api/users")
	require.Contains(t, interaction.Endpoints["/api/users"], "GET")
	require.Contains(t, interaction.Endpoints["/api/users"], "POST")
	require.Equal(t, []string{"fetch user data"}, interaction.Endpoints["/api/users"]["GET"].Reasons)
	require.Equal(t, []string{"create user"}, interaction.Endpoints["/api/users"]["POST"].Reasons)
}
