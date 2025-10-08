ALTER TABLE applications
DROP COLUMN IF EXISTS has_open_api,
DROP COLUMN IF EXISTS open_api_path,
DROP COLUMN IF EXISTS has_open_client,
DROP COLUMN IF EXISTS open_client_path;