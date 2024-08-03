-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE todos
RENAME COLUMN description to task;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE todos
RENAME COLUMN task to description;