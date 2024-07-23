CREATE TABLE IF NOT EXISTS groups (
    "id"                uuid            PRIMARY KEY         NOT NULL        DEFAULT gen_random_uuid(),
    "name"              VARCHAR(255)                        NOT NULL,
    "application_id"    uuid                                NOT NULL        REFERENCES applications(id) ON DELETE CASCADE,
    "permissions"       JSONB                               NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS groups;
