// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: applications.sql

package postgresql

import (
	"context"

	"github.com/google/uuid"
)

const getApplication = `-- name: GetApplication :one
SELECT id, name, secret, keys FROM applications
WHERE
    id = $1
`

func (q *Queries) GetApplication(ctx context.Context, id uuid.UUID) (Application, error) {
	row := q.db.QueryRow(ctx, getApplication, id)
	var i Application
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Secret,
		&i.Keys,
	)
	return i, err
}

const getApplicationByName = `-- name: GetApplicationByName :one
SELECT id, name, secret, keys FROM applications
WHERE
    name = $1
`

func (q *Queries) GetApplicationByName(ctx context.Context, name string) (Application, error) {
	row := q.db.QueryRow(ctx, getApplicationByName, name)
	var i Application
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Secret,
		&i.Keys,
	)
	return i, err
}

const insertApplication = `-- name: InsertApplication :one
INSERT INTO applications
    ( "name", "keys" ) VALUES
    ( $1, $2 )
RETURNING "id"
`

type InsertApplicationParams struct {
	Name string
	Keys []string
}

func (q *Queries) InsertApplication(ctx context.Context, arg InsertApplicationParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, insertApplication, arg.Name, arg.Keys)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
