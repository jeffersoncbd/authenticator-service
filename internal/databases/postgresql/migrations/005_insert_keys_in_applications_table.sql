ALTER TABLE applications ADD COLUMN
    keys      VARCHAR(255)[]    NOT NULL;

---- create above / drop below ----

ALTER TABLE applications DROP COLUMN secret;
