CREATE TABLE IF NOT EXISTS dependency_relationships (
    id SERIAL PRIMARY KEY,
    consumer_id INTEGER REFERENCES applications(id) ON DELETE CASCADE,
    provider_id INTEGER REFERENCES applications(id) ON DELETE CASCADE,
    reasons TEXT[],
    endpoints JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (consumer_id, provider_id)
);