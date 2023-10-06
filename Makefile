# Variables
GO_BIN = ~/go/bin
APP_NAME = psgc
SRC_DIR = ./cmd/http

.PHONY: dev
dev:
	$(GO_BIN)/air

.PHONY: build
build:
	go build -o $(APP_NAME) $(SRC_DIR)

.PHONY: docs
docs:
	$(GO_BIN)/swag fmt && $(GO_BIN)/swag init -d ./cmd/http,./internal/api,./internal/generator && ./docs/fix.sh