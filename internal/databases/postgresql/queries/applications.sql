-- name: ListApplicaions :many
SELECT id, name, keys FROM applications
ORDER BY name ASC;

-- name: GetApplicationByName :one
SELECT * FROM applications
WHERE
    name = $1;

-- name: GetApplication :one
SELECT * FROM applications
WHERE
    id = $1;

-- name: InsertApplication :one
INSERT INTO applications
    ( "name", "keys" ) VALUES
    ( $1, $2 )
RETURNING "id";
