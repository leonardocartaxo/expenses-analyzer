-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users"
(
    "id"         uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "name"       text,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_users_deleted_at" ON "users" ("deleted_at","deleted_at");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_name" ON "users" ("name");

CREATE TABLE "expenses"
(
    "id"         uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "value"      decimal,
    "user_id"    uuid,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_expenses_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_expenses_deleted_at" ON "expenses" ("deleted_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "expenses";
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
