CREATE TABLE IF NOT EXISTS users (
    "email"             VARCHAR(255)        NOT NULL            UNIQUE,
    "name"              VARCHAR(255)        NOT NULL,
    "password"          VARCHAR(255)        NOT NULL,
    "status"            VARCHAR(20)         NOT NULL            DEFAULT 'active',
    "groups"            JSONB               NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;
