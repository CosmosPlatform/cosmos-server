CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    encrypted_value TEXT NOT NULL,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX tokens_team_id_idx ON tokens(team_id);
CREATE INDEX tokens_name_idx ON tokens(name);

-- Add unique constraint on name within a team
ALTER TABLE tokens
ADD CONSTRAINT tokens_team_id_name_unique UNIQUE (team_id, name);
