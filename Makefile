# The Go module lives in cmd/.
MODULE  := github.com/EduardKav/bento-layout/cmd
# Tags for a module in cmd/ are prefixed (cmd/v1.2.3); strip it. Empty output (no git) falls back to dev.
VERSION ?= $(or $(shell git describe --tags --always --dirty 2>/dev/null | sed 's|^cmd/||'),dev)
LDFLAGS := -X $(MODULE)/internal/version.Version=$(VERSION)
BIN     := bento-layout

.PHONY: build test lint fmt check clean help

help: ## List available targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-6s %s\n", $$1, $$2}'

build: ## Build the binary
	cd cmd && go build -ldflags "$(LDFLAGS)" -o ../$(BIN) ./bento-layout

test: ## Run Go tests (with the race detector)
	cd cmd && go test -race ./...

lint: ## Run golangci-lint (also checks formatting)
	cd cmd && golangci-lint run ./...

fmt: ## Format Go code
	cd cmd && golangci-lint fmt ./...

check: lint test ## Run lint and tests (what CI runs)

clean: ## Remove build output
	rm -rf dist $(BIN)
