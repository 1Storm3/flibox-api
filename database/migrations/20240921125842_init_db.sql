-- +goose Up
CREATE TABLE films
(
    id               INT PRIMARY KEY,
    name_ru          VARCHAR(255),
    name_original    VARCHAR(255),
    year             INT,
    poster_url       VARCHAR(1000),
    description      VARCHAR(1000),
    logo_url         VARCHAR(1000),
    type             VARCHAR(255),
    rating_kinopoisk DECIMAL(3, 1),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE film_sequels
(
    film_id   INT NOT NULL,
    sequel_id INT NOT NULL,
    FOREIGN KEY (film_id) REFERENCES films (id) ON DELETE CASCADE,
    FOREIGN KEY (sequel_id) REFERENCES films (id) ON DELETE CASCADE,
    PRIMARY KEY (film_id, sequel_id)
);

CREATE TABLE users
(
    id         UUID      DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name VARCHAR(255) NULL,
    avatar     VARCHAR(1000) NULL,
    tg_id      VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NULL,
    user_token VARCHAR(255) NULL
);


-- +goose Down
DROP TABLE film_sequels;
DROP TABLE films;
DROP TABLE users;