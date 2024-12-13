-- +goose Up
ALTER TABLE comments
    ALTER COLUMN content DROP NOT NULL;

-- +goose Down
ALTER TABLE comments
    ALTER COLUMN content SET NOT NULL;
