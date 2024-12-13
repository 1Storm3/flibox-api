-- +goose Up
ALTER TABLE user_films
    ADD COLUMN type VARCHAR(255);
-- +goose Down
ALTER TABLE user_films
    DROP COLUMN type;
