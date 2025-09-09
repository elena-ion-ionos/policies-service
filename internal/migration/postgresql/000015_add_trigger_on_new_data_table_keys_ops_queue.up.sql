CREATE OR REPLACE FUNCTION on_new_keys_ops_function()
RETURNS trigger
LANGUAGE plpgsql
AS
$$
DECLARE num_rows INTEGER := 0;
BEGIN
    SELECT COUNT(*) INTO num_rows FROM keys_ops_queue;
    if num_rows = 0 THEN
        PERFORM pg_notify('notify-worker', NULL);
    END IF;
    RETURN NEW;
END;
    $$;

CREATE TRIGGER new_keys_ops BEFORE INSERT
    ON keys_ops_queue
    EXECUTE PROCEDURE on_new_keys_ops_function();