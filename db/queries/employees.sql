-- name: GetEmployee :one
select * from employees where id = $1 and deleted_at is null;

-- name: GetEmployees :many
select * from employees where workplace_id = $1 and deleted_at is null;

-- name: GetEmployeesByOffice :many
select employees.*
from employees
join workplaces on employees.workplace_id = workplaces.id
where workplaces.office_id = $1 and employees.deleted_at is null;

-- name: GetEmployeeOffice :one
select workplaces.office_id
from employees join workplaces on employees.workplace_id = workplaces.id
where employees.id = $1 and employees.deleted_at is null;

-- name: CreateEmployee :one
insert into employees (name, workplace_id) values ($1, $2) returning *;

-- name: SoftDeleteEmployee :exec
update employees set deleted_at = now() where id = $1;

-- name: UpdateEmployeeWorkplace :exec
update employees set workplace_id = $2 where id = $1 and deleted_at is null;
