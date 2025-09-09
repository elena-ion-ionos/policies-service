CREATE TABLE if not exists counter_key_id (
                                             user_id uuid PRIMARY KEY,
                                             key_id BIGINT NOT NULL,
                                             contract_number BIGINT NOT NULL
);
