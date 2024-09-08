-- name: ListUsers :many
SELECT name, email, status FROM users
WHERE
    application_id = $1
ORDER BY name ASC;

-- name: GetUser :one
SELECT * FROM users
WHERE
    email = $2 AND application_id = $1;

-- name: InsertUser :exec
INSERT INTO users
    ( "email", "name", "password", "application_id", "group_id" ) VALUES
    ( $1, $2, $3, $4, $5 );

-- name: UpdateUserStatus :exec
UPDATE users
SET status = $3
WHERE email = $2 AND application_id = $1;
