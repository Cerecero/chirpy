-- name: QueryUser :one
SELECT id, created_at, updated_at, email FROM users WHERE email = $1;

-- name: QueryPassword :one
SELECT hashed_password FROM users;
