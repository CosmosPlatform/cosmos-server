ALTER TABLE applications
ADD COLUMN git_provider VARCHAR(100),
ADD COLUMN git_repository_owner VARCHAR(255),
ADD COLUMN git_repository_name VARCHAR(500),
ADD COLUMN git_repository_branch VARCHAR(255);