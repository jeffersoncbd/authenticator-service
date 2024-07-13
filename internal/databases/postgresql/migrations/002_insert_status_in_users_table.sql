ALTER TABLE users ADD COLUMN
    status      VARCHAR(20)     DEFAULT 'active';

---- create above / drop below ----

ALTER TABLE users DROP COLUMN status;
