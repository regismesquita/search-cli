.PHONY: build clean test test-coverage test-verbose install lint run setup help local-install

# Build variables
BINARY_NAME=search
BUILD_DIR=bin
CMD_DIR=cmd/search
COVERAGE_DIR=coverage
GO_BIN=$(HOME)/go/bin

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Build the project
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Clean build files
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) $(COVERAGE_DIR)
	$(GOCLEAN)

# Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated in $(COVERAGE_DIR)/coverage.html"

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	$(GOTEST) -v ./...

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run

# Install the binary
install:
	@echo "Installing to $(GO_BIN)..."
	@mkdir -p $(GO_BIN)
	$(GOBUILD) -o $(GO_BIN)/$(BINARY_NAME) ./$(CMD_DIR)
	@if ! echo "$$PATH" | grep -q "$(GO_BIN)"; then \
		echo "Warning: $(GO_BIN) is not in your PATH"; \
		echo "Add this line to your ~/.bashrc or ~/.zshrc:"; \
		echo "  export PATH=\$$PATH:$(GO_BIN)"; \
	fi

# Install to /usr/local/bin (requires sudo)
local-install: build
	@echo "Installing to /usr/local/bin..."
	@sudo install -m 755 $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Run the application (requires query argument)
run:
	@if [ -z "$(query)" ]; then \
		echo "Usage: make run query='your search query'"; \
		exit 1; \
	fi
	@$(BUILD_DIR)/$(BINARY_NAME) $(query)

# Setup project structure
setup:
	@echo "Setting up project structure..."
	@mkdir -p cmd/search internal/adapters internal/formatter
	@if [ -f main.go ]; then \
		echo "Moving main.go to cmd/search/..."; \
		rm cmd/search/main.go || true; \
		mv main.go cmd/search/; \
	fi
	@if [ -d adapters ]; then \
		echo "Moving adapter files to internal/adapters/..."; \
		mv adapters/* internal/adapters/ 2>/dev/null || true; \
		rmdir adapters; \
	fi
	@if [ -f main_test.go ]; then \
		echo "Moving main_test.go to internal/formatter/formatter_test.go..."; \
		mv main_test.go internal/formatter/formatter_test.go; \
	fi
	@if [ -f search.py ]; then \
		echo "Removing Python version..."; \
		rm search.py; \
	fi
	@echo "Updating go.mod..."
	@$(GOMOD) init github.com/regismesquita/search-cli || true
	@$(GOMOD) tidy
	@echo "Project structure setup complete"

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build the project"
	@echo "  make clean          - Clean build files"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make deps           - Download dependencies"
	@echo "  make lint           - Run linter"
	@echo "  make install        - Install binary to ~/go/bin"
	@echo "  make local-install  - Install binary to /usr/local/bin"
	@echo "  make run           - Run the application (requires query='your search query')"
	@echo "  make setup         - Setup project structure"

# Default target
.DEFAULT_GOAL := help