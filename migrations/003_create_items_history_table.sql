-- +goose Up

CREATE TYPE item_status AS ENUM ('create', 'update', 'delete');

CREATE TABLE IF NOT EXISTS items_history
(
    id         UUID PRIMARY KEY,
    item_id    UUID        NOT NULL REFERENCES items (id) ON DELETE CASCADE,
    action     item_status NOT NULL,
    user_id    UUID        REFERENCES users (id) ON DELETE SET NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    old_data   JSONB,
    new_data   JSONB
);

CREATE INDEX IF NOT EXISTS idx_history_item_id ON items_history (item_id);
CREATE INDEX IF NOT EXISTS idx_history_user_id ON items_history (user_id);
CREATE INDEX IF NOT EXISTS idx_history_action ON items_history (action);
CREATE INDEX IF NOT EXISTS idx_history_changed_at ON items_history (changed_at);

-- +goose Down
DROP TABLE IF EXISTS items_history;
DROP TYPE IF EXISTS item_status;

