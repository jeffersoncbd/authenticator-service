-- name: AddOrUpdateKeyInGroup :exec
UPDATE groups
SET
    permissions = jsonb_set(permissions, $3, $4, true)
WHERE
    id = $2 AND application_id = $1;

-- name: RemoveKeyInGroup :exec
UPDATE groups
SET
    permissions = permissions - $3
WHERE
    id = $2 AND application_id = $1;