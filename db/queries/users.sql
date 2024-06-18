-- name: CreateUser :one
WITH MaxId AS (
    SELECT COALESCE(MAX(id), 0) AS max_id
    FROM users
    WHERE office_id = $1
),
NewId AS (
    SELECT max_id + 1 AS new_id
    FROM MaxId
)
INSERT INTO users (id, office_id, name, password, role, employee_id)
VALUES (
    (SELECT new_id FROM NewId),
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :one
select * from users where id = $1 and office_id = $2;