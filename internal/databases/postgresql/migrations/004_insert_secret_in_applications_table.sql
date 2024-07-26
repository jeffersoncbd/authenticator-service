ALTER TABLE applications ADD COLUMN
    "secret"      uuid      NOT NULL        DEFAULT gen_random_uuid();

---- create above / drop below ----

ALTER TABLE applications DROP COLUMN secret;
