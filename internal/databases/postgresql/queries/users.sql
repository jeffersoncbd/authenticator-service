-- name: ListUsers :many
SELECT u.name, u.email, u.status, g.id AS "group_id", g.name AS "group_name"
FROM
    users u
JOIN
    groups g
ON
    u.group_id = g.id
WHERE
    u.application_id = $1
ORDER BY u.name ASC;

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
