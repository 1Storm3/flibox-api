-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS sequels
(
    sequel_id INT4 PRIMARY KEY,
    name_ru TEXT,
    name_original TEXT,
    poster_url TEXT,
    CONSTRAINT fk_film FOREIGN KEY (film_id) REFERENCES films (external_id)
    );

CREATE TABLE IF NOT EXISTS films_sequels
(
    film_id INT4,
    sequel_id INT4,
    CONSTRAINT fk_film FOREIGN KEY (film_id) REFERENCES films (external_id),
    CONSTRAINT fk_sequel FOREIGN KEY (sequel_id) REFERENCES sequels (sequel_id),
    PRIMARY KEY (film_id, sequel_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS films_sequels;
DROP TABLE IF EXISTS sequels;
-- +goose StatementEnd
