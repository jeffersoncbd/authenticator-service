-- name: ListUsers :many
SELECT name, email, status FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE
    email = $1;

-- name: InsertUser :exec
INSERT INTO users
    ( "email", "name", "password" ) VALUES
    ( $1, $2, $3 );
