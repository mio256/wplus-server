// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: office.sql

package rdb

import (
	"context"
)

const allOffices = `-- name: AllOffices :many
select id, name, deleted_at, created_at, updated_at from offices where deleted_at is null
`

func (q *Queries) AllOffices(ctx context.Context) ([]Office, error) {
	rows, err := q.db.Query(ctx, allOffices)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Office
	for rows.Next() {
		var i Office
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.DeletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const createOffice = `-- name: CreateOffice :one
insert into offices (name) values ($1) returning id, name, deleted_at, created_at, updated_at
`

func (q *Queries) CreateOffice(ctx context.Context, name string) (Office, error) {
	row := q.db.QueryRow(ctx, createOffice, name)
	var i Office
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const softDeleteOffice = `-- name: SoftDeleteOffice :exec
update offices set deleted_at = now() where id = $1
`

func (q *Queries) SoftDeleteOffice(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, softDeleteOffice, id)
	return err
}
