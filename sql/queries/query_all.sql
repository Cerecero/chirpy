-- name: QueryChirp :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: QueryChirpById :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at ASC;
