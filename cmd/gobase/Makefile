# GoBase CLI Makefile

.PHONY: build install clean test help

# Build the CLI binary
build:
	go build -o gobase-cli .

# Install the CLI globally
install:
	go install .

# Clean build artifacts
clean:
	rm -f gobase-cli

# Run tests
test:
	cd ../../ && go test -v

# Show help
help:
	@echo "Available commands:"
	@echo "  build    - Build the CLI binary"
	@echo "  install  - Install the CLI globally"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  help     - Show this help message"

# Default target
all: build
