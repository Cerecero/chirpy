-- name: QueryChirp :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: QueryChirpById :many
SELECT * FROM chirps WHERE id = $1 ORDER BY created_at ASC;
