ALTER TABLE users ADD COLUMN
    "groups"    JSONB    NULL;

---- create above / drop below ----

ALTER TABLE users DROP COLUMN secret;
