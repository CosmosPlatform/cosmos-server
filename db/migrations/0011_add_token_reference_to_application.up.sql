ALTER TABLE applications
ADD COLUMN token_id INTEGER REFERENCES tokens(id) ON DELETE SET NULL;

CREATE INDEX applications_token_id_idx ON applications(token_id);
