// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: for_load.sql

package rdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const loadCreateUser = `-- name: LoadCreateUser :one
insert into users (id, office_id, name, password, role, employee_id)
values ($1, $2, $3, $4, $5, $6)
returning id, office_id, name, password, role, employee_id, created_at, updated_at
`

type LoadCreateUserParams struct {
	ID         int64       `json:"id"`
	OfficeID   int64       `json:"office_id"`
	Name       string      `json:"name"`
	Password   string      `json:"password"`
	Role       UserType    `json:"role"`
	EmployeeID pgtype.Int8 `json:"employee_id"`
}

func (q *Queries) LoadCreateUser(ctx context.Context, arg LoadCreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, loadCreateUser,
		arg.ID,
		arg.OfficeID,
		arg.Name,
		arg.Password,
		arg.Role,
		arg.EmployeeID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.OfficeID,
		&i.Name,
		&i.Password,
		&i.Role,
		&i.EmployeeID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
