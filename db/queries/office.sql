-- name: AllOffices :many
select * from offices where deleted_at is null;

-- name: CreateOffice :one
insert into offices (name) values ($1) returning *;

-- name: SoftDeleteOffice :exec
update offices set deleted_at = now() where id = $1;
