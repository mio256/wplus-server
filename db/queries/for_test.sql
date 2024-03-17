-- name: TestCreateOffice :one
insert into offices (name) values ($1) returning *;

-- name: TestDeleteOffice :exec
delete from offices where id = $1;

-- name: TestCheckDeletedOffice :one
select deleted_at from offices where id = $1;

-- name: TestCreateWorkplace :one
insert into workplaces (name, office_id, work_type)
values ($1, $2, $3)
returning *;

-- name: TestDeleteWorkplace :exec
delete from workplaces where id = $1;

-- name: TestCheckDeletedWorkplace :one
select deleted_at from workplaces where id = $1;

-- name: TestCreateEmployee :one
insert into employees (name, workplace_id) values ($1, $2) returning *;

-- name: TestDeleteEmployee :exec
delete from employees where id = $1;

-- name: TestCheckDeletedEmployee :one
select deleted_at from employees where id = $1;

-- name: TestCreateWorkEntry :one
insert into work_entries (employee_id, workplace_id, date, hours, start_time, end_time, attendance, comment)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: TestDeleteWorkEntry :exec
delete from work_entries where id = $1;

-- name: TestCheckDeletedWorkEntry :one
select deleted_at from work_entries where id = $1;

-- name: TestCreateUser :one
insert into users (id, office_id, name, password, role) values ($1, $2, $3, $4, $5) returning *;

-- name: TestDeleteUser :exec
delete from users where id = $1;