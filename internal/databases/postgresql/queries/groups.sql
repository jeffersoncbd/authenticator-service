-- name: GetGroupByName :one
SELECT * FROM groups
WHERE
    name = $1 AND application_id = $2;

-- name: GetPermissionsGroup :one
SELECT permissions FROM groups
WHERE
    id = $1;

-- name: InsertGroup :one
INSERT INTO groups
    ( "name", "application_id", "permissions" ) VALUES
    ( $1, $2, $3 )
RETURNING "id";
