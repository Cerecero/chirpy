-- +goose Up
CREATE TABLE refresh_tokens(token TEXT PRIMARY KEY, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL REFERENCES USERS(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP
);

-- +goose Down 
DROP TABLE refresh_tokens;
