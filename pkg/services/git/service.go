package git

import (
	"context"
	"cosmos-server/pkg/model"
)

type Service interface {
	GetFileMetadata(ctx context.Context, owner, repo, branch, path string) (*model.FileMetadata, error)
	GetFileWithContent(ctx context.Context, owner, repo, branch, path string) (*model.FileContent, error)
}
