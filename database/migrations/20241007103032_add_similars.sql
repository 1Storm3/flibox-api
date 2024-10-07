-- +goose Up
CREATE TABLE film_similars
(
    film_id    INT NOT NULL,
    similar_id INT NOT NULL,
    FOREIGN KEY (film_id) REFERENCES films (id) ON DELETE CASCADE,
    FOREIGN KEY (similar_id) REFERENCES films (id) ON DELETE CASCADE,
    PRIMARY KEY (film_id, similar_id)
);
-- +goose Down
DROP TABLE film_similars
