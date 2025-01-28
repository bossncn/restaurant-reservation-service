# Makefile

# Variables
DOCKER_IMAGE_NAME = restaurant-reservation-service
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
	go test ./... -v -cover

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

# Run Docker image
.PHONY: docker-run
docker-run:
	@echo "Running Docker image..."
	docker run -d -p 8080:8080 --name $(DOCKER_IMAGE_NAME) $(DOCKER_IMAGE_NAME)

# Clean up Docker images
.PHONY: docker-clean
docker-clean:
	@echo "Cleaning up Docker images..."
	docker rmi $(DOCKER_IMAGE_NAME)

# Generate Swagger Documents
.PHONY: generate-docs
generate-docs:
	@echo "Generating Swagger"
	swag init -d cmd,internal/adapters/http --parseDependency --parseInternal

.PHONY: install-tools
install-tools:
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: generate-mock
generate-mock:
	mockgen -source=internal/core/repository/tables.go -destination=internal/core/repository/mock/mock_table_repository.go
	mockgen -source=internal/core/repository/reservations.go -destination=internal/core/repository/mock/mock_reservation_repository.go
	mockgen -source=internal/core/service/tables.go -destination=internal/core/service/mock/mock_table_service.go
	mockgen -source=internal/core/service/reservations.go -destination=internal/core/service/mock/mock_reservation_service.go
