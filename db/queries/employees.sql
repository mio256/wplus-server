-- name: CreateEmployee :one
insert into employees (name, workplace_id) values ($1, $2) returning *;

-- name: SoftDeleteEmployee :exec
update employees set deleted_at = now() where id = $1;
