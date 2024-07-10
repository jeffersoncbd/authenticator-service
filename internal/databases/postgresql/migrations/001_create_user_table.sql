CREATE TABLE IF NOT EXISTS users (
    email           VARCHAR(255)        NOT NULL            UNIQUE,
    name            VARCHAR(255)        NOT NULL,
    password        VARCHAR(255)        NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;
