.PHONY: test lint fmt vet build examples clean help

# Default target
.DEFAULT_GOAL := help

# Variables
GO := go
GOFMT := gofmt
GOVET := go vet
GOTEST := go test
COVERAGE_DIR := coverage

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated in $(COVERAGE_DIR)/coverage.html"

test-pkg: ## Run tests for pkg packages only
	$(GOTEST) -v ./pkg/...

fmt: ## Format code with gofmt
	$(GOFMT) -l -w .

fmt-check: ## Check if code is formatted
	@test -z "$$($(GOFMT) -l . | tee /dev/stderr)" || (echo "Code is not formatted. Run 'make fmt' to fix." && exit 1)

vet: ## Run go vet
	$(GOVET) ./...

lint: fmt-check vet ## Run all linting checks

build: ## Build all examples
	@for dir in examples/*/*; do \
		if [ -d "$$dir" ]; then \
			echo "Building $$dir..."; \
			$(GO) build "./$$dir" || exit 1; \
		fi \
	done
	@echo "All examples built successfully"

examples: build ## Alias for build

clean: ## Clean build artifacts
	@rm -rf $(COVERAGE_DIR)
	@find . -name "*.test" -type f -delete
	@find examples -type f -name "main" -delete
	@echo "Cleaned build artifacts"

check: lint test ## Run all checks (lint + test)

ci: check ## Run CI checks (alias for check)

