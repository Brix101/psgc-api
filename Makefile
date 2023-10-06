# Variables
GO_BIN = ~/go/bin
APP_NAME = psgc
SRC_DIR = ./cmd/http

.PHONY: dev-api
dev:
	$(GO_BIN)/air api

.PHONY: dev-gen
generate:
	go run $(SRC_DIR) generate

.PHONY: build
build:
	go build -o $(APP_NAME) $(SRC_DIR)

.PHONY: docs
docs:
	$(GO_BIN)/swag fmt && $(GO_BIN)/swag init -d ./cmd/http,./internal/api,./internal/generator,./internal/domain && ./docs/fix.sh
