ALTER TABLE applications
ADD COLUMN dependencies_sha VARCHAR(64),
ADD COLUMN open_api_sha VARCHAR(64);

CREATE TABLE IF NOT EXISTS application_open_apis (
    id SERIAL PRIMARY KEY,
    application_id INTEGER REFERENCES applications(id) ON DELETE CASCADE,
    open_api JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (application_id)
);

CREATE INDEX application_open_apis_application_id_idx ON application_open_apis(application_id);