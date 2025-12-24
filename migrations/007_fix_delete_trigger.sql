-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS items_history_trigger ON items;
DROP TRIGGER IF EXISTS items_history_delete_trigger ON items;

-- Триггер для INSERT и UPDATE (AFTER)
CREATE TRIGGER items_history_trigger
    AFTER INSERT OR UPDATE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_changes();

CREATE TRIGGER items_history_delete_trigger
    BEFORE DELETE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_changes();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS items_history_trigger ON items;
DROP TRIGGER IF EXISTS items_history_delete_trigger ON items;

CREATE TRIGGER items_history_trigger
    AFTER INSERT OR UPDATE OR DELETE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_changes();
-- +goose StatementEnd

