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

-- name: InsertKey :exec
UPDATE applications
    SET keys = (
        SELECT array_agg(DISTINCT unnested_keys)
        FROM unnest(array_cat(keys, $1)) AS unnested_keys
    )
    WHERE id = $2;

-- name: RemoveKey :exec
UPDATE applications
    SET keys = array_remove(keys, $1)
    WHERE id = $2;
