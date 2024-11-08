-- +goose Up
CREATE TABLE comments
(
    id          UUID        DEFAULT uuid_generate_v4() PRIMARY KEY,
    film_id     INT         NOT NULL,
    user_id     UUID        NOT NULL,
    parent_id   UUID        NULL,
    content     TEXT        NOT NULL,
    created_at  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_film FOREIGN KEY (film_id) REFERENCES films (id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_parent_comment FOREIGN KEY (parent_id) REFERENCES comments (id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS comments;
