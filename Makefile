# GoBase Makefile

.PHONY: help build test clean version install lint release

# Variables
BINARY_NAME=gobase
CMD_DIR=./cmd/gobase
BUILD_DIR=./build
VERSION_FILE=cmd/gobase/version.go

# Get version from version.go
VERSION := $(shell grep -o 'Version.*=.*"v[^"]*"' $(VERSION_FILE) | sed 's/.*"v\([^"]*\)".*/\1/')
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-X main.BuildDate=$(BUILD_DATE) -X main.GitCommit=$(GIT_COMMIT)

help: ## Show this help message
	@echo "GoBase v$(VERSION) - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the CLI binary
	@echo "Building GoBase CLI v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ Built: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build for all platforms
	@echo "Building GoBase CLI v$(VERSION) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	
	@echo "✅ Built all platform binaries in $(BUILD_DIR)/"

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests completed"

test-coverage: test ## Run tests and show coverage
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "✅ Cleaned"

install: build ## Install the CLI to GOPATH/bin
	@echo "Installing GoBase CLI..."
	go install -ldflags "$(LDFLAGS)" $(CMD_DIR)
	@echo "✅ Installed to GOPATH/bin"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run
	@echo "✅ Linting completed"

# Version management commands (similar to poetry)
version: ## Show current version
	@echo "Current version: v$(VERSION)"
	@echo "Build date: $(BUILD_DATE)"
	@echo "Git commit: $(GIT_COMMIT)"

version-patch: ## Bump patch version (1.0.0 -> 1.0.1)
	@./scripts/version.sh patch

version-minor: ## Bump minor version (1.0.0 -> 1.1.0)
	@./scripts/version.sh minor

version-major: ## Bump major version (1.0.0 -> 2.0.0)
	@./scripts/version.sh major

version-alpha: ## Set prerelease to alpha (1.0.0 -> 1.0.0-alpha)
	@./scripts/version.sh prerelease alpha

version-beta: ## Set prerelease to beta (1.0.0 -> 1.0.0-beta)
	@./scripts/version.sh prerelease beta

version-rc: ## Set prerelease to rc (1.0.0 -> 1.0.0-rc)
	@./scripts/version.sh prerelease rc

# Release commands
tag: ## Create git tag for current version
	@echo "Creating tag v$(VERSION)..."
	git tag v$(VERSION)
	@echo "✅ Created tag v$(VERSION)"
	@echo "Push with: git push origin main --tags"

release: clean test lint build-all tag ## Full release process
	@echo "🚀 Release v$(VERSION) prepared!"
	@echo ""
	@echo "Files created:"
	@ls -la $(BUILD_DIR)/
	@echo ""
	@echo "Next steps:"
	@echo "1. git push origin main --tags"
	@echo "2. Create GitHub release with binaries from $(BUILD_DIR)/"

# Development commands
dev-build: ## Quick development build
	go build -o $(BINARY_NAME) $(CMD_DIR)

dev-run: dev-build ## Build and run with version flag
	./$(BINARY_NAME) -version

# Docker commands (if needed later)
docker-build: ## Build Docker image
	docker build -t gobase:v$(VERSION) .

.DEFAULT_GOAL := help
