-- name: ListApplicaions :many
SELECT id, name FROM applications
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
    ( "name" ) VALUES ( $1 )
RETURNING "id";

-- name: RenameApplication :exec
UPDATE applications
SET
    "name" = $2
WHERE
    id = $1;
