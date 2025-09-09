CREATE TABLE if not exists backend_keys (
    id uuid PRIMARY KEY,
    contract_number BIGINT NOT NULL,
    user_id uuid NOT NULL,
    key_id uuid NOT NULL ,
    backend VARCHAR(255) NOT NULL,
    access_key VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL
);

CREATE INDEX idx_contract_number_user_id_backend ON backend_keys (contract_number, user_id);