MODULE  := github.com/EduardKav/bento-layout
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X $(MODULE)/internal/version.Version=$(VERSION)
BIN     := bento-layout

.PHONY: build test lint fmt check clean help

help: ## List available targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-6s %s\n", $$1, $$2}'

build: ## Build the binary
	go build -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/bento-layout

test: ## Run Go tests (with the race detector)
	go test -race ./...

lint: ## Run golangci-lint (also checks formatting)
	golangci-lint run ./...

fmt: ## Format Go code
	golangci-lint fmt ./...

check: lint test ## Run lint and tests (what CI runs)

clean: ## Remove build output
	rm -rf dist $(BIN)
