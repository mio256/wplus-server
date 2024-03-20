-- name: LoadCreateUser :one
insert into users (id, office_id, name, password, role, employee_id)
values ($1, $2, $3, $4, $5, $6)
returning *;