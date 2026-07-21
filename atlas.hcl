locals {
  pg_user     = getenv("POSTGRES_USER")
  pg_password = getenv("POSTGRES_PASSWORD")
  pg_host     = getenv("POSTGRES_HOST")
  pg_port     = getenv("POSTGRES_PORT")
  pg_db       = getenv("POSTGRES_DB")
  pg_sslmode  = getenv("POSTGRES_SSLMODE")
}

data "external_schema" "gorm" {
  program = [
    "go", "tool", "atlas-provider-gorm",
    "load",
    "--path", "./internal/repository/postgres",
    "--dialect", "postgres",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  url = "postgres://${local.pg_user}:${local.pg_password}@${local.pg_host}:${local.pg_port}/${local.pg_db}?sslmode=${local.pg_sslmode}"
  dev = "docker://postgres:postgres@postgres/18.0/dev?search_path=public"

  migration {
    dir = "file://migrations"
  }

  diff {
    skip {
      drop_column = true
      drop_table  = true
    }
  }
}
