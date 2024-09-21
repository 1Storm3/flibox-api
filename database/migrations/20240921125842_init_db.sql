-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS films
(
    external_id INT4  PRIMARY KEY,
    name_ru TEXT,
    name_original TEXT,
    year INT,
    poster_url TEXT,
    rating_kinopoisk FLOAT,
    description TEXT,
    logo_url TEXT,
    type TEXT
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS films;
-- +goose StatementEnd
