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
    ( "name" ) VALUES
    ( $1)
RETURNING "id";
