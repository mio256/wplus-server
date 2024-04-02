-- name: GetWorkEntriesByEmployee :many
select workplaces.name as workplace_name, employees.name as employee_name, work_entries.*
from work_entries
join employees on work_entries.employee_id = employees.id
join workplaces on work_entries.workplace_id = workplaces.id
where employees.id = $1 and work_entries.deleted_at is null;

-- name: GetWorkEntriesByOffice :many
select workplaces.name as workplace_name, employees.name as employee_name, work_entries.*
from work_entries
join employees on work_entries.employee_id = employees.id
join workplaces on work_entries.workplace_id = workplaces.id
join offices on workplaces.office_id = offices.id
where offices.id = $1 and work_entries.deleted_at is null;

-- name: GetWorkEntriesByWorkplace :many
select workplaces.name as workplace_name, employees.name as employee_name, work_entries.*
from work_entries
join employees on work_entries.employee_id = employees.id
join workplaces on work_entries.workplace_id = workplaces.id
where workplaces.id = $1 and work_entries.deleted_at is null;

-- name: CreateWorkEntry :one
insert into work_entries (employee_id, workplace_id, date, hours, start_time, end_time, attendance, comment)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: SoftDeleteWorkEntry :exec
update work_entries set deleted_at = now() where id = $1;

-- name: SoftDeleteWorkEntriesByEmployee :exec
update work_entries set deleted_at = now() where employee_id = $1 and deleted_at is null;