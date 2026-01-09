.PHONY: help build run

BINARY_NAME=git-stack-watch
MAIN_PATH=.
BUILD_DIR=./bin

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  run      - Run the application"
	@echo "  vet      - Run go vet for static analysis"

build:
	@echo "Building..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

run:
	go run $(MAIN_PATH)

vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "go vet completed"
