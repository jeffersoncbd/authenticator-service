// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: groups.sql

package postgresql

import (
	"context"

	"github.com/google/uuid"
)

const addKeyInGroup = `-- name: AddKeyInGroup :exec
UPDATE groups
SET
    permissions = jsonb_set(permissions, $3, $4, true)
WHERE
    id = $2 AND application_id = $1
`

type AddKeyInGroupParams struct {
	ApplicationID uuid.UUID
	ID            uuid.UUID
	Path          interface{}
	Replacement   []byte
}

func (q *Queries) AddKeyInGroup(ctx context.Context, arg AddKeyInGroupParams) error {
	_, err := q.db.Exec(ctx, addKeyInGroup,
		arg.ApplicationID,
		arg.ID,
		arg.Path,
		arg.Replacement,
	)
	return err
}

const getGroupByName = `-- name: GetGroupByName :one
SELECT id, name, application_id, permissions FROM groups
WHERE
    name = $1 AND application_id = $2
`

type GetGroupByNameParams struct {
	Name          string
	ApplicationID uuid.UUID
}

func (q *Queries) GetGroupByName(ctx context.Context, arg GetGroupByNameParams) (Group, error) {
	row := q.db.QueryRow(ctx, getGroupByName, arg.Name, arg.ApplicationID)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApplicationID,
		&i.Permissions,
	)
	return i, err
}

const getPermissionsGroup = `-- name: GetPermissionsGroup :one
SELECT permissions FROM groups
WHERE
    id = $1
`

func (q *Queries) GetPermissionsGroup(ctx context.Context, id uuid.UUID) ([]byte, error) {
	row := q.db.QueryRow(ctx, getPermissionsGroup, id)
	var permissions []byte
	err := row.Scan(&permissions)
	return permissions, err
}

const insertGroup = `-- name: InsertGroup :one
INSERT INTO groups
    ( "name", "application_id", "permissions" ) VALUES
    ( $1, $2, $3 )
RETURNING "id"
`

type InsertGroupParams struct {
	Name          string
	ApplicationID uuid.UUID
	Permissions   []byte
}

func (q *Queries) InsertGroup(ctx context.Context, arg InsertGroupParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, insertGroup, arg.Name, arg.ApplicationID, arg.Permissions)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const listGrousByApplicationId = `-- name: ListGrousByApplicationId :many
SELECT id, name, application_id, permissions FROM groups
WHERE
    application_id = $1
`

func (q *Queries) ListGrousByApplicationId(ctx context.Context, applicationID uuid.UUID) ([]Group, error) {
	rows, err := q.db.Query(ctx, listGrousByApplicationId, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Group
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ApplicationID,
			&i.Permissions,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
