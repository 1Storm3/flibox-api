-- +goose Up
Alter table users
    add column last_activity TIMESTAMP;

-- +goose Down
Alter table users
    drop column last_activity;
