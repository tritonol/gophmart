-- +goose Up
-- +goose StatementBegin
ALTER TABLE balance DROP COLUMN type;
ALTER TABLE balance DROP COLUMN "from";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
