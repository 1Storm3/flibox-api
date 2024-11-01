-- +goose Up
DROP TABLE user_films;
DROP TABLE users;
CREATE TABLE users
(
    id          UUID         DEFAULT uuid_generate_v4() PRIMARY KEY,
    nick_name   VARCHAR(255)                NOT NULL,
    name        VARCHAR(255)                NOT NULL,
    email       VARCHAR(255)                NOT NULL,
    password    VARCHAR(255)                NOT NULL,
    photo       VARCHAR(1000) NULL,
    role        VARCHAR(255) DEFAULT 'user' NOT NULL,
    is_verified BOOLEAN      DEFAULT false,
    created_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NULL
);
CREATE TABLE user_films
(
    user_id UUID NOT NULL,
    film_id INT  NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES films (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, film_id)
);
-- +goose Down
DROP TABLE users
DROP TABLE user_films
