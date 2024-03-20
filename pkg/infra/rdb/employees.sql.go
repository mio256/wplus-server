// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: employees.sql

package rdb

import (
	"context"
)

const createEmployee = `-- name: CreateEmployee :one
insert into employees (name, workplace_id) values ($1, $2) returning id, name, workplace_id, deleted_at, created_at, updated_at
`

type CreateEmployeeParams struct {
	Name        string `json:"name"`
	WorkplaceID int64  `json:"workplace_id"`
}

func (q *Queries) CreateEmployee(ctx context.Context, arg CreateEmployeeParams) (Employee, error) {
	row := q.db.QueryRow(ctx, createEmployee, arg.Name, arg.WorkplaceID)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.WorkplaceID,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getEmployee = `-- name: GetEmployee :one
select id, name, workplace_id, deleted_at, created_at, updated_at from employees where workplace_id = $1 and id = $2 and deleted_at is null
`

type GetEmployeeParams struct {
	WorkplaceID int64 `json:"workplace_id"`
	ID          int64 `json:"id"`
}

func (q *Queries) GetEmployee(ctx context.Context, arg GetEmployeeParams) (Employee, error) {
	row := q.db.QueryRow(ctx, getEmployee, arg.WorkplaceID, arg.ID)
	var i Employee
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.WorkplaceID,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getEmployeeOffice = `-- name: GetEmployeeOffice :one
select workplaces.office_id
from employees
    join workplaces on employees.workplace_id = workplaces.id
where employees.id = $1 and employees.deleted_at is null
`

func (q *Queries) GetEmployeeOffice(ctx context.Context, id int64) (int64, error) {
	row := q.db.QueryRow(ctx, getEmployeeOffice, id)
	var office_id int64
	err := row.Scan(&office_id)
	return office_id, err
}

const getEmployees = `-- name: GetEmployees :many
select id, name, workplace_id, deleted_at, created_at, updated_at from employees where workplace_id = $1 and deleted_at is null
`

func (q *Queries) GetEmployees(ctx context.Context, workplaceID int64) ([]Employee, error) {
	rows, err := q.db.Query(ctx, getEmployees, workplaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Employee
	for rows.Next() {
		var i Employee
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.WorkplaceID,
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

const softDeleteEmployee = `-- name: SoftDeleteEmployee :exec
update employees set deleted_at = now() where id = $1
`

func (q *Queries) SoftDeleteEmployee(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, softDeleteEmployee, id)
	return err
}

const updateEmployeeWorkplace = `-- name: UpdateEmployeeWorkplace :exec
update employees set workplace_id = $2 where id = $1 and deleted_at is null
`

type UpdateEmployeeWorkplaceParams struct {
	ID          int64 `json:"id"`
	WorkplaceID int64 `json:"workplace_id"`
}

func (q *Queries) UpdateEmployeeWorkplace(ctx context.Context, arg UpdateEmployeeWorkplaceParams) error {
	_, err := q.db.Exec(ctx, updateEmployeeWorkplace, arg.ID, arg.WorkplaceID)
	return err
}
