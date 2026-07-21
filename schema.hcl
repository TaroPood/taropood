schema "public" {
  comment = "standard public schema for taropood"
}

table "rule_models" {
  schema = schema.public
  column "id" {
    null = false
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "condition_type" {
    null = false
    type = text
  }
  column "condition_config" {
    null = true
    type = jsonb
  }
  column "priority" {
    null = false
    type = integer
    default = 0
  }
  column "enabled" {
    null = false
    type = boolean
    default = true
  }
  column "tags" {
    null = true
    type = jsonb
  }
  column "metadata" {
    null = true
    type = jsonb
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  column "updated_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_rules_name" {
    columns = [column.name]
  }
  index "idx_rules_enabled" {
    columns = [column.enabled]
  }
}

table "action_models" {
  schema = schema.public
  column "id" {
    null = false
    type = text
  }
  column "rule_id" {
    null = false
    type = text
  }
  column "type" {
    null = false
    type = text
  }
  column "config" {
    null = true
    type = jsonb
  }
  column "order" {
    null = false
    type = integer
    default = 0
  }
  column "created_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_action_models_rule" {
    columns     = [column.rule_id]
    ref_columns = [table.rule_models.column.id]
    on_delete   = CASCADE
  }
  index "idx_rule_actions_rule_id" {
    columns = [column.rule_id]
  }
}
