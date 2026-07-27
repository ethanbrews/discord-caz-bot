# Project configuration
BINARY_NAME=discord_caz_bot
MAIN_PACKAGE=.  # Change to ./cmd/discord_caz_bot if using a cmd directory
BUILD_DIR=bin

# Docker configuration
DOCKER_IMAGE_NAME=discord-caz-bot
DOCKER_TAG=latest

.PHONY: all build clean run test tidy docker-build docker-run help

## default: Builds the local binary
all: build

## build: Compiles the Go binary into bin/
build:
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

## run: Builds and executes the binary locally
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

## docker-build: Builds the Docker image
docker-build:
	@echo "Building Docker image: $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)..."
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

## docker-run: Runs the Docker container locally
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it --env-file .env $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

## test: Runs all tests with race detector
test:
	@echo "Running tests..."
	go test -v -race ./...

## tidy: Cleans up go.mod and go.sum dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

## clean: Removes built binaries and temporary files
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

## help: Shows available make commands
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'