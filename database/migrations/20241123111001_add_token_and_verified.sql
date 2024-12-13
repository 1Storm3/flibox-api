-- +goose Up
ALTER TABLE users
    ADD COLUMN verified_token VARCHAR(255),
    ADD COLUMN is_blocked BOOLEAN DEFAULT false;
-- +goose Down
ALTER TABLE users
    DROP COLUMN verified_token,
    DROP COLUMN is_blocked;

