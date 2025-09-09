CREATE TABLE if not exists keys_backends_mapping (
    key_id uuid NOT NULL ,
    backend VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    constraint pk_keys_backends_mapping primary key (key_id, backend)
    );

CREATE INDEX idx_keys_backends_mapping_key_id ON keys_backends_mapping(key_id);