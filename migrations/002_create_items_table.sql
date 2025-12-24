-- +goose Up
CREATE TABLE IF NOT EXISTS items
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(64) NOT NULL,
    description TEXT,
    quantity    INT         NOT NULL CHECK (quantity >= 0),
    price       INT         NOT NULL CHECK (price >= 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_items_name ON items (name);
CREATE INDEX IF NOT EXISTS idx_items_created_at ON items (created_at);

-- +goose Down
DROP TABLE IF EXISTS items;
