.PHONY: build test clean install run help

# Variables
BINARY_NAME=AeyeWire_mcp
BUILD_DIR=build
MAIN_FILE=src/AeyeWire_mcp.go

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the service
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run health check
health: build
	@./$(BUILD_DIR)/$(BINARY_NAME) health

# List supported languages
languages: build
	@./$(BUILD_DIR)/$(BINARY_NAME) languages

# Show version
version: build
	@./$(BUILD_DIR)/$(BINARY_NAME) version

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run || go vet ./...

# Help
help:
	@echo "AeyeWire MCP Service - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build      - Build the binary"
	@echo "  test       - Run tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install dependencies"
	@echo "  run        - Build and run the MCP service"
	@echo "  health     - Check service health"
	@echo "  languages  - List supported languages"
	@echo "  version    - Show version"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  help       - Show this help message"
