# Makefile

# Variables
DOCKER_IMAGE_NAME = go-boilerplate
DOCKERFILE_PATH = Dockerfile

# Default target
.PHONY: all
all: go-mod-tidy go-format go-test

# Run the application
.PHONY: go-run
go-run:
	@echo "Running the application..."
	go run cmd/main.go

# Tidy Go modules
.PHONY: go-mod-tidy
go-mod-tidy:
	@echo "Tidying Go modules..."
	go mod tidy

# Run Go tests
.PHONY: go-test
go-test:
	@echo "Running tests..."
	go test ./...

# Format Go code
.PHONY: go-format
go-format:
	@echo "Formatting Go code..."
	go fmt ./...

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) -f $(DOCKERFILE_PATH) .

# Clean up Docker images
.PHONY: docker-clean
docker-clean:
	@echo "Cleaning up Docker images..."
	docker rmi $(DOCKER_IMAGE_NAME)
