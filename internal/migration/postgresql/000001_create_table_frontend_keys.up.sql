CREATE TABLE if not exists frontend_keys (
     id uuid PRIMARY KEY,
     contract_number BIGINT NOT NULL,
     user_id uuid NOT NULL,
     description VARCHAR(1000),
     access_key VARCHAR(255) NOT NULL,
     encrypted_secret_key VARCHAR(255) NOT NULL,
     created_at   timestamp  DEFAULT ('now'::text)::timestamp without time zone,
     updated_at   timestamp not null
);

CREATE INDEX idx_contract_number_user_id_frontend ON frontend_keys (contract_number, user_id);