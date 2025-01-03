// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: applications.sql

package postgresql

import (
	"context"

	"github.com/google/uuid"
)

const getApplication = `-- name: GetApplication :one
SELECT id, name, secret FROM applications
WHERE
    id = $1
`

func (q *Queries) GetApplication(ctx context.Context, id uuid.UUID) (Application, error) {
	row := q.db.QueryRow(ctx, getApplication, id)
	var i Application
	err := row.Scan(&i.ID, &i.Name, &i.Secret)
	return i, err
}

const getApplicationByName = `-- name: GetApplicationByName :one
SELECT id, name, secret FROM applications
WHERE
    name = $1
`

func (q *Queries) GetApplicationByName(ctx context.Context, name string) (Application, error) {
	row := q.db.QueryRow(ctx, getApplicationByName, name)
	var i Application
	err := row.Scan(&i.ID, &i.Name, &i.Secret)
	return i, err
}

const insertApplication = `-- name: InsertApplication :one
INSERT INTO applications
    ( "name" ) VALUES ( $1 )
RETURNING "id"
`

func (q *Queries) InsertApplication(ctx context.Context, name string) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, insertApplication, name)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const listApplicaions = `-- name: ListApplicaions :many
SELECT id, name FROM applications
ORDER BY name ASC
`

type ListApplicaionsRow struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) ListApplicaions(ctx context.Context) ([]ListApplicaionsRow, error) {
	rows, err := q.db.Query(ctx, listApplicaions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListApplicaionsRow
	for rows.Next() {
		var i ListApplicaionsRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const renameApplication = `-- name: RenameApplication :exec
UPDATE applications
SET
    "name" = $2
WHERE
    id = $1
`

type RenameApplicationParams struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) RenameApplication(ctx context.Context, arg RenameApplicationParams) error {
	_, err := q.db.Exec(ctx, renameApplication, arg.ID, arg.Name)
	return err
}
