-- +goose Up
CREATE TABLE history_films
(
    id         UUID      DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    film_id    INT NOT NULL REFERENCES films (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL
);

-- +goose Down
DROP TABLE history_films CASCADE;
