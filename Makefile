.PHONY: test lint fmt vet build examples clean help tools golangci-lint staticcheck fmt-check

# Default target
.DEFAULT_GOAL := help

# Variables
GO := go
GOFMT := gofmt
GOBIN_DIR := $(shell $(GO) env GOBIN)
GOPATH_DIR := $(shell $(GO) env GOPATH)
GO_BIN_DIRS := $(strip $(if $(GOBIN_DIR),$(GOBIN_DIR)) $(GOPATH_DIR)/bin)
define find_tool
$(firstword $(foreach dir,$(GO_BIN_DIRS),$(if $(wildcard $(dir)/$(1)),$(dir)/$(1),)) $(shell command -v $(1) 2>/dev/null))
endef
GOIMPORTS := $(call find_tool,goimports)
STATICCHECK := $(call find_tool,staticcheck)
GOLANGCI_LINT := $(call find_tool,golangci-lint)
GOVET := go vet
GOTEST := go test
COVERAGE_DIR := coverage
TOOLS_SCRIPT := scripts/install_tools.sh

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

tools: ## Install optional development tools
	@$(TOOLS_SCRIPT)

test: ## Run all tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated in $(COVERAGE_DIR)/coverage.html"

test-pkg: ## Run tests for pkg packages only
	$(GOTEST) -v ./pkg/...

fmt: ## Format code with gofmt or goimports if available
ifdef GOIMPORTS
	@echo "Running goimports formatting"
	$(GOIMPORTS) -w .
else
	@echo "goimports not found; falling back to gofmt"
	$(GOFMT) -l -w .
endif

fmt-check: ## Check if code is formatted
ifdef GOIMPORTS
	@test -z "$$($(GOIMPORTS) -l . | tee /dev/stderr)" || (echo "Code is not formatted. Run 'make fmt' to fix." && exit 1)
else
	@test -z "$$($(GOFMT) -l . | tee /dev/stderr)" || (echo "Code is not formatted. Run 'make fmt' to fix." && exit 1)
endif

vet: ## Run go vet
	$(GOVET) ./...

lint: fmt-check vet golangci-lint ## Run all linting checks

golangci-lint: ## Run golangci-lint
ifdef GOLANGCI_LINT
	$(GOLANGCI_LINT) run ./...
else
	@echo "golangci-lint not installed. Run 'make tools' or './scripts/install_tools.sh' first."
	@exit 1
endif

staticcheck: ## Run staticcheck directly
ifdef STATICCHECK
	$(STATICCHECK) ./...
else
	@echo "staticcheck not installed. Run 'make tools' or './scripts/install_tools.sh' first."
	@exit 1
endif

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
