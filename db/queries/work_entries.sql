-- name: GetWorkEntriesByEmployee :many
select * from work_entries where employee_id = $1 and deleted_at is null;

-- name: CreateWorkEntry :one
insert into work_entries (employee_id, workplace_id, date, hours, start_time, end_time, attendance, comment)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: SoftDeleteWorkEntry :exec
update work_entries set deleted_at = now() where id = $1;
