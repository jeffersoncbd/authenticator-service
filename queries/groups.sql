-- name: GetGroupByName :one
SELECT * FROM groups
WHERE
    name = $1 AND application_id = $2;

-- name: GetPermissionsGroup :one
SELECT permissions FROM groups
WHERE
    id = $1;

-- name: ListGrousByApplicationId :many
SELECT * FROM groups
WHERE
    application_id = $1;

-- name: InsertGroup :one
INSERT INTO groups
    ( "name", "application_id", "permissions" ) VALUES
    ( $1, $2, $3 )
RETURNING "id";

-- name: AddKeyInGroup :exec
UPDATE groups
SET
    permissions = jsonb_set(permissions, $3, $4, true)
WHERE
    id = $2 AND application_id = $1;
