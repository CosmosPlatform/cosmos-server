CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    encrypted_password VARCHAR(255) NOT NULL,
    role VARCHAR(50) CHECK (role IN ('admin', 'user'))
);

CREATE INDEX users_email_idx ON users(email);
CREATE INDEX users_role_idx ON users(role);