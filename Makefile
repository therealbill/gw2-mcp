# GW2 MCP Server Makefile

.DELETE_ON_ERROR:
.DEFAULT_GOAL := help

# ── Tools ────────────────────────────────────────────────────────────────────
GO          := go
GOFUMPT     := gofumpt
GOLANGCI    := golangci-lint
GOVULNCHECK := govulncheck

# ── Platform ─────────────────────────────────────────────────────────────────
ifeq ($(OS),Windows_NT)
  EXT := .exe
else
  EXT :=
endif

BINARY := bin/gw2-mcp$(EXT)

# ── Version info (lazy evaluation — only resolved when used) ─────────────────
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  = $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# ── Targets ──────────────────────────────────────────────────────────────────

.PHONY: all
all: format vet lint test build ## Run full pipeline (format, vet, lint, test, build)

.PHONY: build
build: ## Build the server binary
	@echo "Building GW2 MCP Server..."
	$(GO) build -o $(BINARY) -ldflags "$(LDFLAGS)" ./

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	$(GO) clean

.PHONY: test
test: ## Run tests with race detection
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out ./...

.PHONY: coverage
coverage: test ## Generate HTML coverage report
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html

.PHONY: bench
bench: ## Run benchmarks
	$(GO) test -bench=. -benchmem -run=^$$ ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

.PHONY: lint
lint: ## Run linter (golangci-lint)
	@echo "Running linter..."
	$(GOLANGCI) run ./...

.PHONY: format
format: ## Format code and tidy modules
	@echo "Formatting code..."
	$(GOFUMPT) -w .
	$(GO) mod tidy

.PHONY: deps
deps: ## Install/update dependencies
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod verify

.PHONY: run
run: build ## Build and run the server
	@echo "Starting GW2 MCP Server..."
	./$(BINARY)

.PHONY: dev
dev: ## Run in development mode with race detection
	@echo "Starting GW2 MCP Server in development mode..."
	$(GO) run -race ./

.PHONY: tools
tools: ## Install development tools
	@echo "Installing development tools..."
	$(GO) install mvdan.cc/gofumpt@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: security
security: ## Check for security vulnerabilities (govulncheck)
	@echo "Checking for security vulnerabilities..."
	$(GOVULNCHECK) ./...

.PHONY: docs
docs: ## List package documentation
	@echo "Generating documentation..."
	@$(GO) list ./... | while read pkg; do echo "=== $$pkg ==="; $(GO) doc -all "$$pkg"; done

.PHONY: generate
generate: ## Run go generate
	$(GO) generate ./...

.PHONY: check
check: vet lint test security ## Run all quality checks (vet, lint, test, security)

.PHONY: release
release: ## Build optimized release binary
	@echo "Building release version..."
	CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) ./

.PHONY: docker
docker: ## Build Docker image
	@echo "Building Docker image..."
	docker build --build-arg VERSION=$(VERSION) -t gw2-mcp:latest .

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# ── Debugging ────────────────────────────────────────────────────────────────
print-%: ## Print any variable (e.g., make print-VERSION)
	@echo '$*=$($*)'
