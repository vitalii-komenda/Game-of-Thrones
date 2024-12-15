-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters DROP COLUMN updated_at;
-- +goose StatementEnd
