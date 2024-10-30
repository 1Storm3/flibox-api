-- +goose Up
CREATE TABLE user_films
(
    user_id UUID NOT NULL,
    film_id INT          NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES films (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, film_id)
);
-- +goose Down
DROP TABLE user_films

