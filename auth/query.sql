-- name: GetUser :one
SELECT * FROM auth.users
WHERE email = $1 LIMIT 1;