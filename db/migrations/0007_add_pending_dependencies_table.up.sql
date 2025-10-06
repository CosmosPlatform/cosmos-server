CREATE TABLE IF NOT EXISTS pending_application_dependencies (
    id SERIAL PRIMARY KEY,
    consumer_id INTEGER REFERENCES applications(id) ON DELETE CASCADE,
    provider_name VARCHAR(255) NOT NULL,
    reasons TEXT[],
    endpoints JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (consumer_id, provider_name)
);

CREATE INDEX IF NOT EXISTS pending_application_dependencies_consumer_id_idx ON pending_application_dependencies(consumer_id);
CREATE INDEX IF NOT EXISTS pending_application_dependencies_provider_name_idx ON pending_application_dependencies(provider_name);

CREATE INDEX IF NOT EXISTS idx_application_dependencies_consumer_id ON application_dependencies(consumer_id);
CREATE INDEX IF NOT EXISTS idx_application_dependencies_provider_id ON application_dependencies(provider_id);
