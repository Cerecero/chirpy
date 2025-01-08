-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password TEXT DEFAULT 'unset';
-- +goose Down 
DROP TABLE users;
