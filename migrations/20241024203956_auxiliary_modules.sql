-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tags"
(
    "id"         uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "name"       text,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_tags_name" ON "tags" ("name");

CREATE INDEX IF NOT EXISTS "idx_tags_deleted_at" ON "tags" ("deleted_at","deleted_at");

CREATE TABLE "collectors"
(
    "id"         uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "name"       text,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_collectors_name" ON "collectors" ("name");

CREATE INDEX IF NOT EXISTS "idx_collectors_deleted_at" ON "collectors" ("deleted_at","deleted_at");

CREATE TABLE "places"
(
    "id"         uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "name"       text,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_places_name" ON "places" ("name");

CREATE INDEX IF NOT EXISTS "idx_places_deleted_at" ON "places" ("deleted_at","deleted_at");

ALTER TABLE "expenses"
    ADD "name" text NOT NULL;

ALTER TABLE "expenses"
    ALTER COLUMN "value" SET NOT NULL;

ALTER TABLE "expenses"
    ALTER COLUMN "user_id" SET NOT NULL;

ALTER TABLE "expenses"
    ADD "collector_id" uuid NOT NULL;

ALTER TABLE "expenses"
    ADD "place_id" uuid NOT NULL;

ALTER TABLE "expenses"
    ADD CONSTRAINT "fk_expenses_collector" FOREIGN KEY ("collector_id") REFERENCES "collectors" ("id");

ALTER TABLE "expenses"
    ADD CONSTRAINT "fk_expenses_place" FOREIGN KEY ("place_id") REFERENCES "places" ("id");

CREATE TABLE "expense_tags"
(
    "model_id" uuid DEFAULT gen_random_uuid(),
    "tag_id"   uuid DEFAULT gen_random_uuid(),
    PRIMARY KEY ("model_id", "tag_id"),
    CONSTRAINT "fk_expense_tags_model" FOREIGN KEY ("model_id") REFERENCES "expenses" ("id"),
    CONSTRAINT "fk_expense_tags_tags" FOREIGN KEY ("tag_id") REFERENCES "tags" ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "tags";
DROP TABLE IF EXISTS "expense_tags";
DROP TABLE IF EXISTS "places";
DROP TABLE IF EXISTS "collectors";
-- +goose StatementEnd
