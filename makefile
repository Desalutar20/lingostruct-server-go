MAIN_GO_DIR=./cmd/api
BINARY_NAME=lingostruct-server-go
BUILD_DIR=./bin

GO_VERSION := $(shell go version | cut -d ' ' -f 3)

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building project..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_GO_DIR)

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./tests/e2e/... -count=1

.PHONY: clean
clean:
	@echo "Cleaning the build..."
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: run
run: build
	@echo "Running the project..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	@goose up

fmt:
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make build        - Build the project"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean the build"
	@echo "  make run          - Run the project"
	@echo "  make migrate      - Run the database migrations"
	@echo "  make fmt          - Format the code"
	@echo "  make lint         - Run the linter"
	@echo "  make deps         - Install dependencies"
