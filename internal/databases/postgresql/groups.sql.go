// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: groups.sql

package postgresql

import (
	"context"

	"github.com/google/uuid"
)

const getGroup = `-- name: GetGroup :one
SELECT id, name, application_id, permissions FROM groups
WHERE
    id = $1 AND application_id = $2
`

type GetGroupParams struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
}

func (q *Queries) GetGroup(ctx context.Context, arg GetGroupParams) (Group, error) {
	row := q.db.QueryRow(ctx, getGroup, arg.ID, arg.ApplicationID)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ApplicationID,
		&i.Permissions,
	)
	return i, err
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
