CREATE TABLE if not exists keys_ops_queue (
    key_id uuid NOT NULL ,
    operation VARCHAR(255) NOT NULL,
    created_at timestamp  DEFAULT ('now'::text)::timestamp without time zone
    );

CREATE INDEX idx_keys_ops_queue_key_id ON keys_ops_queue(key_id);