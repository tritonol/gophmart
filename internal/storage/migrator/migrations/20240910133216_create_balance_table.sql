-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "balance" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER,
    "from" VARCHAR,
    "from_id" BIGINT,
    "value" BIGINT,
    "type" INTEGER,

    CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP IF EXISTS balance;
-- +goose StatementEnd
