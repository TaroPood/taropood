APP_NAME   ?= rule-engine
BIN_DIR    ?= .
GO         ?= go
DOCKER     ?= docker

BINARY      = $(BIN_DIR)/$(APP_NAME)
MIGRATE_BIN = $(BIN_DIR)/migrate
ATLAS       = atlas

DOCKER_COMPOSE = $(DOCKER) compose -f deployments/docker/docker-compose.yaml
DOCKER_FILE    = deployments/docker/Dockerfile.dev

GO_BUILD_LDFLAGS = -ldflags="-X main.version=$(VERSION)"

.PHONY: all build run test lint tidy gen fmt vet
.PHONY: migrate-diff migrate-apply migrate-down migrate-hash migrate-validate migrate-status
.PHONY: docker-up docker-down docker-build docker-logs dev clean help

all: tidy gen fmt vet build

# ---- Build ----

build:
	$(GO) build -o $(BINARY) ./cmd/$(APP_NAME)

build-migrate:
	$(GO) build -o $(MIGRATE_BIN) ./cmd/migrate

build-all: build build-migrate

# ---- Run ----

run:
	$(GO) run ./cmd/$(APP_NAME)/...

# ---- Development ----

dev:
	@command -v air >/dev/null 2>&1 || { \
		echo "Installing air..."; \
		$(GO) install github.com/air-verse/air@latest; \
	}
	air

# ---- Code Quality ----

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/"; \
		exit 1; \
	}
	golangci-lint run

tidy:
	$(GO) mod tidy

# ---- Code Generation ----

gen:
	$(GO) run ./tools/gentool/...

# ---- Database Migrations ----

migrate-diff:
	$(GO) run ./cmd/migrate/... diff NAME=$(filter-out $@,$(MAKECMDGOALS))

migrate-apply:
	$(GO) run ./cmd/migrate/... apply

migrate-down:
	$(GO) run ./cmd/migrate/... down $(N)

migrate-hash:
	$(GO) run ./cmd/migrate/... hash

migrate-validate:
	$(GO) run ./cmd/migrate/... validate

migrate-status:
	$(GO) run ./cmd/migrate/... status

# ---- Testing ----

test:
	$(GO) test ./... -count=1

test-race:
	$(GO) test ./... -race -count=1

test-cover:
	$(GO) test ./... -coverprofile=coverage.out -count=1
	$(GO) tool cover -html=coverage.out -o coverage.html

# ---- Docker ----

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-restart: docker-down docker-up

docker-build:
	$(DOCKER_COMPOSE) build

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-ps:
	$(DOCKER_COMPOSE) ps

# ---- Clean ----

clean:
	rm -f $(BINARY) $(MIGRATE_BIN)
	rm -f coverage.out coverage.html
	rm -rf tmp/

# ---- Help ----

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Build:"
	@echo "  build            Build rule-engine binary"
	@echo "  build-migrate    Build migrate CLI binary"
	@echo "  build-all        Build all binaries"
	@echo ""
	@echo "Run:"
	@echo "  run              Run application locally"
	@echo "  dev              Run with air hot-reload"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt              Format Go source files"
	@echo "  vet              Run go vet"
	@echo "  lint             Run golangci-lint"
	@echo "  tidy             Run go mod tidy"
	@echo ""
	@echo "Code Generation:"
	@echo "  gen              Regenerate gorm.io/gen query code"
	@echo ""
	@echo "Migrations (via atlas):"
	@echo "  migrate-diff NAME=xxx   Create a new migration"
	@echo "  migrate-apply           Apply pending migrations"
	@echo "  migrate-down N=1        Rollback N migrations"
	@echo "  migrate-hash            Hash migration directory"
	@echo "  migrate-validate        Validate migrations"
	@echo "  migrate-status          Show migration status"
	@echo ""
	@echo "Testing:"
	@echo "  test             Run all tests"
	@echo "  test-race        Run tests with race detector"
	@echo "  test-cover       Run tests with coverage"
	@echo ""
	@echo "Docker:"
	@echo "  docker-up        Start docker-compose services"
	@echo "  docker-down      Stop docker-compose services"
	@echo "  docker-restart   Restart docker-compose services"
	@echo "  docker-build     Build docker images"
	@echo "  docker-logs      Follow logs"
	@echo "  docker-ps        List containers"
	@echo ""
	@echo "Other:"
	@echo "  clean            Remove build artifacts"
	@echo "  all              tidy -> gen -> fmt -> vet -> build"
	@echo "  help             Show this message"

# Allow targets to accept extra args (e.g. NAME=xxx)
%:
	@:
