CREATE TABLE if not exists keys (
    key_id uuid PRIMARY KEY,
    contract_number BIGINT NOT NULL,
    user_id uuid NOT NULL,
    description VARCHAR(1000),
    access_key VARCHAR(255) NOT NULL,
    encrypted_secret_key VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    created_at timestamp  DEFAULT ('now'::text)::timestamp without time zone,
    updated_at timestamp not null,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    is_owner BOOLEAN DEFAULT TRUE NOT NULL,
    canonical_user_id VARCHAR(255) NULL,
    created_by VARCHAR (320) DEFAULT '',
    created_by_user_id VARCHAR (36) DEFAULT '',
    updated_by VARCHAR (320) DEFAULT '',
    updated_by_user_id VARCHAR (36) DEFAULT ''
    );

CREATE INDEX idx_keys_contract_number ON keys(contract_number);
CREATE INDEX idx_keys_user_id ON keys(user_id);
CREATE INDEX idx_keys_access_key ON keys(access_key);