// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: work_entries.sql

package rdb

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createWorkEntry = `-- name: CreateWorkEntry :one
insert into work_entries (employee_id, workplace_id, date, hours, start_time, end_time, attendance, comment)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning id, employee_id, workplace_id, date, hours, start_time, end_time, attendance, comment, deleted_at, created_at, updated_at
`

type CreateWorkEntryParams struct {
	EmployeeID  int64       `json:"employee_id"`
	WorkplaceID int64       `json:"workplace_id"`
	Date        pgtype.Date `json:"date"`
	Hours       pgtype.Int2 `json:"hours"`
	StartTime   pgtype.Time `json:"start_time"`
	EndTime     pgtype.Time `json:"end_time"`
	Attendance  pgtype.Bool `json:"attendance"`
	Comment     pgtype.Text `json:"comment"`
}

func (q *Queries) CreateWorkEntry(ctx context.Context, arg CreateWorkEntryParams) (WorkEntry, error) {
	row := q.db.QueryRow(ctx, createWorkEntry,
		arg.EmployeeID,
		arg.WorkplaceID,
		arg.Date,
		arg.Hours,
		arg.StartTime,
		arg.EndTime,
		arg.Attendance,
		arg.Comment,
	)
	var i WorkEntry
	err := row.Scan(
		&i.ID,
		&i.EmployeeID,
		&i.WorkplaceID,
		&i.Date,
		&i.Hours,
		&i.StartTime,
		&i.EndTime,
		&i.Attendance,
		&i.Comment,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getWorkEntriesByEmployee = `-- name: GetWorkEntriesByEmployee :many
select workplaces.name as workplace_name, employees.name as employee_name, work_entries.id, work_entries.employee_id, work_entries.workplace_id, work_entries.date, work_entries.hours, work_entries.start_time, work_entries.end_time, work_entries.attendance, work_entries.comment, work_entries.deleted_at, work_entries.created_at, work_entries.updated_at
from work_entries
join employees on work_entries.employee_id = employees.id
join workplaces on work_entries.workplace_id = workplaces.id
where employees.id = $1 and work_entries.deleted_at is null
`

type GetWorkEntriesByEmployeeRow struct {
	WorkplaceName string           `json:"workplace_name"`
	EmployeeName  string           `json:"employee_name"`
	ID            int64            `json:"id"`
	EmployeeID    int64            `json:"employee_id"`
	WorkplaceID   int64            `json:"workplace_id"`
	Date          pgtype.Date      `json:"date"`
	Hours         pgtype.Int2      `json:"hours"`
	StartTime     pgtype.Time      `json:"start_time"`
	EndTime       pgtype.Time      `json:"end_time"`
	Attendance    pgtype.Bool      `json:"attendance"`
	Comment       pgtype.Text      `json:"comment"`
	DeletedAt     pgtype.Timestamp `json:"deleted_at"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) GetWorkEntriesByEmployee(ctx context.Context, id int64) ([]GetWorkEntriesByEmployeeRow, error) {
	rows, err := q.db.Query(ctx, getWorkEntriesByEmployee, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkEntriesByEmployeeRow
	for rows.Next() {
		var i GetWorkEntriesByEmployeeRow
		if err := rows.Scan(
			&i.WorkplaceName,
			&i.EmployeeName,
			&i.ID,
			&i.EmployeeID,
			&i.WorkplaceID,
			&i.Date,
			&i.Hours,
			&i.StartTime,
			&i.EndTime,
			&i.Attendance,
			&i.Comment,
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

const getWorkEntriesByOffice = `-- name: GetWorkEntriesByOffice :many
select workplaces.name as workplace_name, employees.name as employee_name, work_entries.id, work_entries.employee_id, work_entries.workplace_id, work_entries.date, work_entries.hours, work_entries.start_time, work_entries.end_time, work_entries.attendance, work_entries.comment, work_entries.deleted_at, work_entries.created_at, work_entries.updated_at
from work_entries
join employees on work_entries.employee_id = employees.id
join workplaces on work_entries.workplace_id = workplaces.id
join offices on workplaces.office_id = offices.id
where offices.id = $1 and work_entries.deleted_at is null
`

type GetWorkEntriesByOfficeRow struct {
	WorkplaceName string           `json:"workplace_name"`
	EmployeeName  string           `json:"employee_name"`
	ID            int64            `json:"id"`
	EmployeeID    int64            `json:"employee_id"`
	WorkplaceID   int64            `json:"workplace_id"`
	Date          pgtype.Date      `json:"date"`
	Hours         pgtype.Int2      `json:"hours"`
	StartTime     pgtype.Time      `json:"start_time"`
	EndTime       pgtype.Time      `json:"end_time"`
	Attendance    pgtype.Bool      `json:"attendance"`
	Comment       pgtype.Text      `json:"comment"`
	DeletedAt     pgtype.Timestamp `json:"deleted_at"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) GetWorkEntriesByOffice(ctx context.Context, id int64) ([]GetWorkEntriesByOfficeRow, error) {
	rows, err := q.db.Query(ctx, getWorkEntriesByOffice, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkEntriesByOfficeRow
	for rows.Next() {
		var i GetWorkEntriesByOfficeRow
		if err := rows.Scan(
			&i.WorkplaceName,
			&i.EmployeeName,
			&i.ID,
			&i.EmployeeID,
			&i.WorkplaceID,
			&i.Date,
			&i.Hours,
			&i.StartTime,
			&i.EndTime,
			&i.Attendance,
			&i.Comment,
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

const getWorkEntriesByWorkplace = `-- name: GetWorkEntriesByWorkplace :many
select workplaces.name as workplace_name, employees.name as employee_name, work_entries.id, work_entries.employee_id, work_entries.workplace_id, work_entries.date, work_entries.hours, work_entries.start_time, work_entries.end_time, work_entries.attendance, work_entries.comment, work_entries.deleted_at, work_entries.created_at, work_entries.updated_at
from work_entries
join employees on work_entries.employee_id = employees.id
join workplaces on work_entries.workplace_id = workplaces.id
where workplaces.id = $1 and work_entries.deleted_at is null
`

type GetWorkEntriesByWorkplaceRow struct {
	WorkplaceName string           `json:"workplace_name"`
	EmployeeName  string           `json:"employee_name"`
	ID            int64            `json:"id"`
	EmployeeID    int64            `json:"employee_id"`
	WorkplaceID   int64            `json:"workplace_id"`
	Date          pgtype.Date      `json:"date"`
	Hours         pgtype.Int2      `json:"hours"`
	StartTime     pgtype.Time      `json:"start_time"`
	EndTime       pgtype.Time      `json:"end_time"`
	Attendance    pgtype.Bool      `json:"attendance"`
	Comment       pgtype.Text      `json:"comment"`
	DeletedAt     pgtype.Timestamp `json:"deleted_at"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) GetWorkEntriesByWorkplace(ctx context.Context, id int64) ([]GetWorkEntriesByWorkplaceRow, error) {
	rows, err := q.db.Query(ctx, getWorkEntriesByWorkplace, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWorkEntriesByWorkplaceRow
	for rows.Next() {
		var i GetWorkEntriesByWorkplaceRow
		if err := rows.Scan(
			&i.WorkplaceName,
			&i.EmployeeName,
			&i.ID,
			&i.EmployeeID,
			&i.WorkplaceID,
			&i.Date,
			&i.Hours,
			&i.StartTime,
			&i.EndTime,
			&i.Attendance,
			&i.Comment,
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

const getWorkEntry = `-- name: GetWorkEntry :one
select id, employee_id, workplace_id, date, hours, start_time, end_time, attendance, comment, deleted_at, created_at, updated_at from work_entries where id = $1 and deleted_at is null
`

func (q *Queries) GetWorkEntry(ctx context.Context, id int64) (WorkEntry, error) {
	row := q.db.QueryRow(ctx, getWorkEntry, id)
	var i WorkEntry
	err := row.Scan(
		&i.ID,
		&i.EmployeeID,
		&i.WorkplaceID,
		&i.Date,
		&i.Hours,
		&i.StartTime,
		&i.EndTime,
		&i.Attendance,
		&i.Comment,
		&i.DeletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const softDeleteWorkEntriesByEmployee = `-- name: SoftDeleteWorkEntriesByEmployee :exec
update work_entries set deleted_at = now() where employee_id = $1 and deleted_at is null
`

func (q *Queries) SoftDeleteWorkEntriesByEmployee(ctx context.Context, employeeID int64) error {
	_, err := q.db.Exec(ctx, softDeleteWorkEntriesByEmployee, employeeID)
	return err
}

const softDeleteWorkEntry = `-- name: SoftDeleteWorkEntry :exec
update work_entries set deleted_at = now() where id = $1
`

func (q *Queries) SoftDeleteWorkEntry(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, softDeleteWorkEntry, id)
	return err
}
