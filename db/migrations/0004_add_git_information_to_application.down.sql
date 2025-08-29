ALTER TABLE applications
DROP COLUMN IF EXISTS git_provider,
DROP COLUMN IF EXISTS git_repository_owner,
DROP COLUMN IF EXISTS git_repository_name,
DROP COLUMN IF EXISTS git_repository_branch;