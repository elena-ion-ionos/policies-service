CREATE TABLE if not exists users (
     id uuid PRIMARY KEY,
     contract_number BIGINT NOT NULL,
     phone VARCHAR(100),
     email VARCHAR(255) NOT NULL,
     created_at   timestamp  DEFAULT ('now'::text)::timestamp without time zone,
     updated_at   timestamp null
);
