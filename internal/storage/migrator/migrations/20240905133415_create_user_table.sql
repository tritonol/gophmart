-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "login" VARCHAR,
    "password" VARCHAR
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_login ON users(login);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP IF EXISTS user;

-- +goose StatementEnd
