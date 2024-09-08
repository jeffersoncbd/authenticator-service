CREATE TABLE IF NOT EXISTS users (
    "id"                uuid                NOT NULL            PRIMARY KEY     DEFAULT gen_random_uuid(),
    "email"             VARCHAR(255)        NOT NULL,
    "name"              VARCHAR(255)        NOT NULL,
    "password"          VARCHAR(255)        NOT NULL,
    "status"            VARCHAR(20)         NOT NULL            DEFAULT 'active',
    "application_id"    uuid                NOT NULL            REFERENCES applications(id) ON DELETE CASCADE,
    "group_id"          uuid                NOT NULL            REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT          unique_email_per_application UNIQUE (email, application_id)
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;
