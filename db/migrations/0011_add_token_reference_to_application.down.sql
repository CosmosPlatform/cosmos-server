DROP INDEX IF EXISTS applications_token_id_idx;

ALTER TABLE applications
DROP COLUMN IF EXISTS token_id;
