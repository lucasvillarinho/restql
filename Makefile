# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: install-tools
install-tools: ## Install development tools (gotestfmt, etc.)
	@echo "Installing gotestfmt..."
	@go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	@echo "✓ Tools installed successfully"

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format Go code using golangci-lint
	golangci-lint fmt

test: ## Run tests (shows only failures and summary)
	@command -v gotestfmt >/dev/null 2>&1 || { echo "gotestfmt not found. Run 'make install-tools' first."; exit 1; }
	go test -json -race -cover ./... | gotestfmt -hide successful-tests

test-verbose: ## Run tests with verbose output (no formatting)
	go test -v -race -cover ./...

coverage: ## Generate coverage report
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

coverage-html: coverage ## Generate and open HTML coverage report
	go tool cover -html=coverage.out