CREATE TABLE if not exists policies
(
    id         uuid PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    prefix     VARCHAR(255) NOT NULL,
    action     VARCHAR(255) NOT NULL,
    time       VARCHAR(255) NOT NULL,
    created_at timestamp DEFAULT ('now'::text)::timestamp without time zone
);
