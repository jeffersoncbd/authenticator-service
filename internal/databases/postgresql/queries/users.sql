-- name: InsertUser :exec
INSERT INTO users
    ( "email", "name", "password" ) VALUES
    ( $1, $2, $3 );

-- name: GetUser :one
SELECT * FROM users
WHERE
    email = $1;
