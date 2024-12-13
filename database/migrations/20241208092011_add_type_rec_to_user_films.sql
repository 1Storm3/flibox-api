-- +goose Up
ALTER TABLE user_films
    ALTER COLUMN type SET NOT NULL;
ALTER TABLE user_films DROP CONSTRAINT user_films_pkey;

ALTER TABLE user_films
    ADD PRIMARY KEY (user_id, film_id, type);

-- +goose Down
ALTER TABLE user_films DROP CONSTRAINT user_films_pkey;

ALTER TABLE user_films
    ADD PRIMARY KEY (user_id, film_id);