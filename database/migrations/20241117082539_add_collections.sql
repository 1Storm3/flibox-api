-- +goose Up
CREATE TABLE collections
(
    id          UUID      DEFAULT uuid_generate_v4() PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    cover_url   VARCHAR(1000),
    tags        TEXT[],
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    user_id     UUID REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE collection_films
(
    collection_id UUID NOT NULL,
    film_id       INT  NOT NULL,
    FOREIGN KEY (collection_id) REFERENCES collections (id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES films (id) ON DELETE CASCADE,
    PRIMARY KEY (collection_id, film_id)
);

-- +goose Down
DROP TABLE collection_films CASCADE;
DROP TABLE collections CASCADE;
