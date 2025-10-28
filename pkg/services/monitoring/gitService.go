package monitoring

import (
	"context"
	"cosmos-server/pkg/errors"
	"cosmos-server/pkg/model"
	"net/http"
	"os"
	"sync"

	"github.com/google/go-github/v74/github"
	"golang.org/x/oauth2"
)

//go:generate mockgen -destination=./mock/gitService_mock.go -package=mock cosmos-server/pkg/services/monitoring GitService

type GitService interface {
	GetFileMetadata(ctx context.Context, owner, repo, branch, path, token string) (*model.FileMetadata, error)
	GetFileWithContent(ctx context.Context, owner, repo, branch, path, token string) (*model.FileContent, error)
}

type githubService struct {
	defaultToken  string
	clients       sync.Map // map[string]*github.Client
	defaultClient *github.Client
}

func NewGithubService() GitService {
	var httpClient *http.Client

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient = oauth2.NewClient(context.Background(), src)
	}

	client := github.NewClient(httpClient)
	return &githubService{defaultClient: client}
}

func (g *githubService) getClient(token string) *github.Client {
	if token == "" {
		return g.defaultClient
	}

	if client, ok := g.clients.Load(token); ok {
		return client.(*github.Client)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := github.NewClient(httpClient)
	g.clients.Store(token, client)
	return client
}

func (g *githubService) GetFileMetadata(ctx context.Context, owner, repo, branch, path, token string) (*model.FileMetadata, error) {
	tree, _, err := g.getClient(token).Git.GetTree(ctx, owner, repo, branch, true)
	if err != nil {
		return nil, err
	}

	for _, entry := range tree.Entries {
		if entry.GetPath() == path {
			return &model.FileMetadata{
				Name:       entry.GetPath(),
				Path:       entry.GetPath(),
				Size:       entry.GetSize(),
				SHA:        entry.GetSHA(),
				Branch:     branch,
				Repository: repo,
				Owner:      owner,
			}, nil
		}
	}

	return nil, errors.NewNotFoundError("file %s not found in repo %s/%s on branch %s", path, owner, repo, branch)
}

func (g *githubService) GetFileWithContent(ctx context.Context, owner, repo, branch, path, token string) (*model.FileContent, error) {
	file, _, _, err := g.getClient(token).Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: branch})
	if err != nil {
		return nil, err
	}

	content, err := file.GetContent()
	if err != nil {
		return nil, err
	}

	metadata := model.FileMetadata{
		Name:       file.GetName(),
		Path:       file.GetPath(),
		Size:       file.GetSize(),
		SHA:        file.GetSHA(),
		Branch:     branch,
		Repository: repo,
		Owner:      owner,
	}

	return &model.FileContent{Metadata: metadata, Content: content}, nil
}
