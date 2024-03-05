-- name: GetEmployees :many
select * from employees where workplace_id = $1 and deleted_at is null;

-- name: CreateEmployee :one
insert into employees (name, workplace_id) values ($1, $2) returning *;

-- name: SoftDeleteEmployee :exec
update employees set deleted_at = now() where id = $1;
