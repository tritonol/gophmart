-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "orders" (
    "id" BIGINT PRIMARY KEY,
    "user_id" INTEGER,
    "status" VARCHAR,
    "uploaded_at" TIMESTAMP without time zone default (now() at time zone 'utc'),

    CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP IF EXISTS orders;
-- +goose StatementEnd
