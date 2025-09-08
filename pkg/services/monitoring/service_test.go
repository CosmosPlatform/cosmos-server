package monitoring

import (
	"context"
	log "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/model"
	"cosmos-server/pkg/services/monitoring/mock"
	storageMock "cosmos-server/pkg/storage/mock"
	"cosmos-server/pkg/storage/obj"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateApplicationInformation(t *testing.T) {
	t.Run("update application information - success", updateApplicationInformationSuccess)
	t.Run("update application information - no git information", updateApplicationInformationNoGitInformation)
	t.Run("update application information - get file error", updateApplicationInformationGetFileError)
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
