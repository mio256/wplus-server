-- name: LoadCreateUser :one
insert into users (id, office_id, name, password, role)
select $1, $2, $3, $4, $5 where not exists (
    select 1 from users where id = $1
)
returning *;