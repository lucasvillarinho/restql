# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format Go code using golangci-lint
	golangci-lint fmt

test: ## Run tests
	go test -v -race ./...

coverage: ## Generate coverage report
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

coverage-html: coverage ## Generate and open HTML coverage report
	go tool cover -html=coverage.out