-- name: CreateWorkplace :one
insert into workplaces (name, office_id, work_type)
values ($1, $2, $3)
returning *;

-- name: SoftDeleteWorkplace :exec
update workplaces set deleted_at = now() where id = $1;
