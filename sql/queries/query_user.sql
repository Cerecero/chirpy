-- name: QueryUser :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET email = $1, hashed_password = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, created_at, updated_at, email;

-- name: QueryAuthorUser :one
SELECT user_id FROM chirps WHERE id = $1;

-- name: QueryUpgradeUser :exec
UPDATE users SET is_chirpy_red = TRUE WHERE id = $1;
