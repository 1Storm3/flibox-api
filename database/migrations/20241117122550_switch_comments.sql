-- +goose Up

ALTER TABLE comments
DROP
CONSTRAINT IF EXISTS fk_parent_comment;

ALTER TABLE comments
    ADD CONSTRAINT fk_parent_comment FOREIGN KEY (parent_id)
        REFERENCES comments (id) ON DELETE CASCADE;

-- +goose Down

ALTER TABLE comments
DROP
CONSTRAINT IF EXISTS fk_parent_comment;

ALTER TABLE comments
    ADD CONSTRAINT fk_parent_comment FOREIGN KEY (parent_id)
        REFERENCES comments (id) ON DELETE SET NULL;