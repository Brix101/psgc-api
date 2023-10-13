# Variables
GO_BIN = ~/go/bin
APP_NAME = psgc
SRC_DIR = ./cmd/http
DATE := $(shell date +'%Y-%m-%d')
DATABASE = ../db/$(DATE)-data.db  # Define your database file here
GOOSE = goose  # Define the Goose binary (make sure it's in your PATH)
MIGRATIONS_DIR = migrations  # Define the directory where your migrations are located


.PHONY: dev-api
dev-api:
	$(GO_BIN)/air api

.PHONY: dev-gen
dev-gen:
	go run $(SRC_DIR) generate

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: build
build:
	go build -o $(APP_NAME) $(SRC_DIR)

.PHONY: docs
docs:
	$(GO_BIN)/swag fmt && $(GO_BIN)/swag init -d ./cmd/http,./internal/api,./internal/generator,./internal/domain && ./docs/fix.sh

.PHONY: migrate-up
migrate-up:
	cd migrations && $(GOOSE) sqlite3 $(DATABASE) up

.PHONY: migrate-down
migrate-down:
	cd migrations && $(GOOSE) sqlite3 $(DATABASE) down

.PHONY: migrate-reset
migrate-reset:
	cd migrations && $(GOOSE) sqlite3 $(DATABASE) reset

.PHONY: migrate-status
migrate-status:
	cd migrations && $(GOOSE) sqlite3 $(DATABASE) status
