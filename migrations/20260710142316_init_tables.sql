-- Create "rule_models" table
CREATE TABLE "rule_models" (
  "id" text NOT NULL,
  "name" text NOT NULL,
  "condition_type" text NOT NULL,
  "condition_config" jsonb NULL,
  "priority" bigint NOT NULL DEFAULT 0,
  "enabled" boolean NOT NULL DEFAULT true,
  "tags" jsonb NULL,
  "metadata" jsonb NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_rules_enabled" to table: "rule_models"
CREATE INDEX "idx_rules_enabled" ON "rule_models" ("enabled");
-- Create index "idx_rules_name" to table: "rule_models"
CREATE INDEX "idx_rules_name" ON "rule_models" ("name");
-- Create "action_models" table
CREATE TABLE "action_models" (
  "id" text NOT NULL,
  "rule_id" text NOT NULL,
  "type" text NOT NULL,
  "config" jsonb NULL,
  "order" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_rule_models_actions" FOREIGN KEY ("rule_id") REFERENCES "rule_models" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_rule_actions_rule_id" to table: "action_models"
CREATE INDEX "idx_rule_actions_rule_id" ON "action_models" ("rule_id");
