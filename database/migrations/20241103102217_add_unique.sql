-- +goose Up
ALTER TABLE users
    ADD CONSTRAINT unique_email UNIQUE (email),
ADD CONSTRAINT unique_nick_name UNIQUE (nick_name);

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT IF EXISTS unique_email,
DROP CONSTRAINT IF EXISTS unique_nick_name;

