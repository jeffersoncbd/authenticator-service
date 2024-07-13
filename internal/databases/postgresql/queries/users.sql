-- name: ListUsers :many
SELECT name, email, status FROM users
ORDER BY name ASC;

-- name: GetUser :one
SELECT * FROM users
WHERE
    email = $1;

-- name: InsertUser :exec
INSERT INTO users
    ( "email", "name", "password" ) VALUES
    ( $1, $2, $3 );

-- name: UpdateUserStatus :exec
UPDATE users
SET status = $2
WHERE email = $1;
