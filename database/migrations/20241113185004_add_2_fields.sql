-- +goose Up
ALTER TABLE films
    ADD COLUMN cover_url VARCHAR(1000), ADD COLUMN trailer_url VARCHAR(1000);

-- +goose Down
Alter table films
    drop column cover_url, drop column trailer_url
