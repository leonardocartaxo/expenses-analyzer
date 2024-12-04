-- +goose Up
-- +goose StatementBegin
ALTER TABLE "users" ADD "email" text;
ALTER TABLE "users" ADD "password" text;
CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_email" ON "users" ("email");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
