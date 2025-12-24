-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(32) UNIQUE NOT NULL,
    role       VARCHAR(32)        NOT NULL CHECK (role IN ('admin', 'manager', 'viewer')),
    created_at TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_name ON users (name);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);

-- +goose Down
DROP TABLE IF EXISTS users;

