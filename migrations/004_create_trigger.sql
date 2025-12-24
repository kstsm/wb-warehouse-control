-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_item_changes()
RETURNS TRIGGER AS $func$
DECLARE
    old_json JSONB;
    new_json JSONB;
    user_uuid UUID;
    user_id_str TEXT;
BEGIN
    BEGIN
        user_id_str := current_setting('app.user_id', true);
        IF user_id_str IS NULL OR user_id_str = '' THEN
            user_uuid := NULL;
        ELSE
            user_uuid := user_id_str::UUID;
        END IF;
    EXCEPTION WHEN OTHERS THEN
        user_uuid := NULL;
    END;

    IF TG_OP = 'INSERT' THEN
        new_json := jsonb_build_object(
            'id', NEW.id,
            'name', NEW.name,
            'description', NEW.description,
            'quantity', NEW.quantity,
            'price', NEW.price,
            'created_at', NEW.created_at,
            'updated_at', NEW.updated_at
        );
        
        INSERT INTO items_history (id, item_id, action, user_id, old_data, new_data)
        VALUES (gen_random_uuid(), NEW.id, 'create'::item_status, user_uuid, NULL, new_json);
        
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        old_json := jsonb_build_object(
            'id', OLD.id,
            'name', OLD.name,
            'description', OLD.description,
            'quantity', OLD.quantity,
            'price', OLD.price,
            'created_at', OLD.created_at,
            'updated_at', OLD.updated_at
        );
        
        new_json := jsonb_build_object(
            'id', NEW.id,
            'name', NEW.name,
            'description', NEW.description,
            'quantity', NEW.quantity,
            'price', NEW.price,
            'created_at', NEW.created_at,
            'updated_at', NEW.updated_at
        );
        
        INSERT INTO items_history (id, item_id, action, user_id, old_data, new_data)
        VALUES (gen_random_uuid(), NEW.id, 'update'::item_status, user_uuid, old_json, new_json);
        
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        old_json := jsonb_build_object(
            'id', OLD.id,
            'name', OLD.name,
            'description', OLD.description,
            'quantity', OLD.quantity,
            'price', OLD.price,
            'created_at', OLD.created_at,
            'updated_at', OLD.updated_at
        );
        
        INSERT INTO items_history (id, item_id, action, user_id, old_data, new_data)
        VALUES (gen_random_uuid(), OLD.id, 'delete'::item_status, user_uuid, old_json, NULL);
        
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$func$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER items_history_trigger
    AFTER INSERT OR UPDATE OR DELETE ON items
    FOR EACH ROW
    EXECUTE FUNCTION log_item_changes();

-- +goose Down
DROP TRIGGER IF EXISTS items_history_trigger ON items;
-- +goose StatementBegin
DROP FUNCTION IF EXISTS log_item_changes();
-- +goose StatementEnd

