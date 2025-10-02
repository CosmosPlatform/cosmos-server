ALTER TABLE applications
DROP COLUMN IF EXISTS dependencies_sha,
DROP COLUMN IF EXISTS open_api_sha;

DROP TABLE IF EXISTS application_open_apis;
