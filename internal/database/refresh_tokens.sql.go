// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_tokens.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES (
    $1,
    $2,
    $3
)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRefreshTokenParams struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.Token, arg.UserID, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const queryRefreshToken = `-- name: QueryRefreshToken :one
SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token = $1
`

type QueryRefreshTokenRow struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

func (q *Queries) QueryRefreshToken(ctx context.Context, token string) (QueryRefreshTokenRow, error) {
	row := q.db.QueryRowContext(ctx, queryRefreshToken, token)
	var i QueryRefreshTokenRow
	err := row.Scan(&i.UserID, &i.ExpiresAt, &i.RevokedAt)
	return i, err
}
