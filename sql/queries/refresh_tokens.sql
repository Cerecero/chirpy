-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: QueryRefreshToken :one
SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token = $1;
