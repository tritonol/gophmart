-- +goose Up
-- +goose StatementBegin
ALTER TABLE balance 
ALTER COLUMN "value" TYPE float;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
