-- +goose Up
ALTER TABLE films
    ADD COLUMN genres TEXT[];
-- +goose Down
ALTER TABLE films
DROP
COLUMN genres