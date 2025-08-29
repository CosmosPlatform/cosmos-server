package git

import (
	"context"
	"cosmos-server/pkg/model"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/v74/github"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type githubService struct {
	client *github.Client
}

func NewGithubService() Service {
	var httpClient *http.Client

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient = oauth2.NewClient(context.Background(), src)
	}

	client := github.NewClient(httpClient)
	return &githubService{client: client}
}

func (g *githubService) GetFileMetadata(ctx context.Context, owner, repo, branch, path string) (*model.FileMetadata, error) {
	tree, _, err := g.client.Git.GetTree(ctx, owner, repo, branch, true)
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

	return nil, fmt.Errorf("file %s not found in repo %s/%s on branch %s", path, owner, repo, branch)
}

func (g *githubService) GetFileWithContent(ctx context.Context, owner, repo, branch, path string) (*model.FileContent, error) {
	file, _, _, err := g.client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: branch})
	if err != nil {
		return nil, err
	}

	rawContent, err := file.GetContent()
	if err != nil {
		return nil, err
	}

	content, err := base64.StdEncoding.DecodeString(rawContent)
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
