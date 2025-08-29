package model

type FileMetadata struct {
	Name       string
	Path       string
	Size       int
	SHA        string
	Branch     string
	Repository string
	Owner      string
}

type FileContent struct {
	Metadata FileMetadata
	Content  []byte
}
