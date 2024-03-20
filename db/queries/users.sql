-- name: CreateUser :one
insert into users (id, office_id, name, password, role, employee_id)
VALUES ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetUser :one
select * from users where id = $1 and office_id = $2;