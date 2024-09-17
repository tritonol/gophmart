-- +goose Up
-- +goose StatementBegin
ALTER TABLE balance 
ADD COLUMN processed_at TIMESTAMP without time zone default (now() at time zone 'utc')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
