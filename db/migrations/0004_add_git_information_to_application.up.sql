ALTER TABLE applications
ADD COLUMN git_repository VARCHAR(500),
ADD COLUMN git_branch VARCHAR(255);