## Getting Started

### Prerequisites

- Go 1.18 or later
- Docker (optional, for database setup)
- [Echo Framework](https://echo.labstack.com/)

### Installation

1. Clone this repository:

    ```bash
    git clone https://github.com/bossncn/restaurant-reservation-service.git
    cd restaurant-reservation-service
    ```

2. Install Go dependencies:

    ```bash
    go mod tidy
    ```

3. Run the application:

    ```bash
    go run main.go
    ```
   
The application should now be running at `http://localhost:8080`.

Optional

Run with Docker

1. Build

   ```bash
   make docker-build
   ```
   
2. Run docker

   ```bash
   make docker-run
   ```

### Make Commands

Here are the available `make` commands you can use to manage the project:

- **`make go-run`**: Runs the application.
- **`make go-mod-tidy`**: Tidies up Go modules.
- **`make go-test`**: Runs all tests in the project.
- **`make go-format`**: Formats all Go code.
- **`make docker-build`**: Builds the Docker image for the application.
- **`make docker-run`**: Run Docker
- **`make docker-clean`**: Cleans up the Docker image.
- **`make generate-docs`**: Generate Swagger Documents
- **`make generate-mock`**: Generate Mock


### Project Structure

```bash
/go-hexagonal-boilerplate
├── cmd
│   ├── app                 # Initializes HTTP server and acts as the entry point for the application
│   │   └── app             # Core application components for HTTP setup and dependency injection
│   └── main                # Runs the application (starts the HTTP server initialized in 'app')
├── config                  # Configuration files for the app (e.g., environment variables, settings)
├── internal
│   ├── adapter             # Implementations of external systems (e.g., HTTP, DB adapters)
│   │   ├── http            # HTTP handler implementations using Echo framework
│   │   └── memory          # In-memory storage implementation of repository
│   │   └── event           # Event App Process Command request
│   ├── core                # Core business logic
│   │   ├── model           # Core models (e.g., User, Product, etc.)
│   │   ├── repository      # Interfaces for repositories (e.g., data storage logic)
│   │   └── service         # Core service implementations, business logic
│   └── middleware          # Application middleware (e.g., logging, authentication)
├── test                    # Test files
│   └── integration         # Integration tests to test app components together
```

### Running Tests
To run tests:

   ```bash
   make go-test
   ```

## Features

- **Go Boilerplate**: A basic structure to get started with a clean, organized Go project.
- **Hexagonal Architecture**: Clear separation between core business logic and external systems (HTTP, databases).
- **Echo Framework**: Uses the Echo framework for HTTP server and routing.
- **Swagger Documentation**: Auto-generated OpenAPI documentation with [swaggo/swag](https://github.com/swaggo/swag).
- **Zap Logger**: Efficient and structured logging with the [Zap](https://github.com/uber-go/zap) logging library.
- **Docker Support**: Easy Docker integration for running the application in containers.
- **Unit & Integration Tests**: Tests for ensuring functionality and integration across the system.

## Hexagonal Architecture Overview

Hexagonal Architecture is an architectural pattern where the core business logic (inside the "hexagon") is isolated from external systems. This separation helps to achieve:

- **Flexibility**: Easy to swap out external systems (e.g., database, HTTP client) without impacting core logic.
- **Testability**: The core business logic is decoupled from external dependencies, making it easier to unit test.
- **Maintainability**: With clear boundaries between the core and external layers, it's easier to understand and manage code.

### The basic components in this architecture are:

### 1. **Core (Domain)**
The core business logic and domain models are located in the `internal/core` directory. It contains:
- **Models**: Domain models that represent core entities (e.g., `User`, `Product`).
- **Repositories (Ports)**: Interfaces for data access, located in the `internal/core/repository` directory. These repository interfaces (e.g., `UserRepository.go`) define how the core communicates with external data stores, without binding it to a specific implementation.
- **Services**: Business logic services that orchestrate domain operations.

### 2. **Ports**
Ports define the communication channels between the core and external systems (adapters). They are **interfaces** that are located within the core business logic, usually in the `internal/core/repository` or `internal/core/service` files. In our case, these ports represent the abstract definitions for interacting with external systems.

For example:
- The **`UserRepository` interface** defines how the core communicates with data storage (e.g., in-memory or a database). It doesn't specify how the data is fetched, only the methods available for use. This allows the core logic to be agnostic to the underlying data source.
- The **`UserService` interface** could be defined as a service that coordinates with different repositories, which is another example of a port where the core business logic communicates via an interface to external services.

### 3. **Adapters**
Adapters implement the **ports** and act as bridges between the core business logic and external systems. They are located in the `internal/adapter` directory and include:
- **HTTP Adapters**: Implementations for handling HTTP requests (e.g., in the `internal/adapter/http` directory using the Echo framework).
- **Memory Adapters**: In-memory implementations of repositories or services (e.g., `internal/adapter/memory`).
- **External Services**: Adapters for connecting to external services like databases, APIs, etc.

Adapters implement the interfaces (ports) defined by the core, translating data between external systems and the business logic. For example:
- The **`UserRepository` adapter** in `internal/adapter/http` or `internal/adapter/memory` provides concrete implementations for storing and retrieving `User` data, while adhering to the `UserRepository` interface defined in the core (`internal/core/repository/UserRepository.go`).

### 4. **Testability and Separation**
- **Unit Tests**: Unit tests for the core components (services, repositories, models) should be placed in the `internal/test` directory.
- **Integration Tests**: Integration tests for the application as a whole (e.g., HTTP requests and interactions between components) should be placed in the `test/integration` directory.

By separating these concerns into distinct layers, it becomes much easier to write unit tests for the core logic without worrying about external dependencies, and it simplifies integration testing of the entire system.

---

In this setup, the core logic is always isolated from the external systems, making it easier to maintain and test. The adapters allow external systems to interact with the core logic without tightly coupling the application to any particular implementation.


## Contributing
Feel free to open issues or create pull requests to improve the boilerplate. Contributions are always welcome!

## License

This project is licensed under the [MIT License](LICENSE). See the [LICENSE](LICENSE) file for more details.
