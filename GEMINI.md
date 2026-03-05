# UK Housing Market & Rental Intelligence API

## Project Overview

This project provides a comprehensive API for analyzing and accessing UK housing market and rental data, leveraging `dataset.csv` as the primary data source. The system is designed for high performance, scalability, and maintainability using Go.

## Architecture & Design

The project follows a strictly layered architecture to ensure clean separation of concerns and high testability.

### Layered Architecture

**Router → Handler → Service → Repository → Database**

- **Router:** Defines API endpoints and connects them to handlers.
- **Handler:** Manages HTTP/GraphQL request parsing, validation, and response formatting.
- **Service:** Contains the core business logic. It is the only layer that should coordinate between different repositories or external services.
- **Repository:** Responsible for all data access logic. It interacts directly with the database via GORM.
- **Database:** PostgreSQL for persistent storage.

### Design Principles

- **Separation of Concerns:** Each layer has a single, well-defined responsibility.
- **Dependency Injection:** Dependencies are passed into constructors as interfaces, facilitating loose coupling and easier testing.
- **Interface-based Repository Abstraction:** All repositories must implement an interface defined in the service layer (or a shared repository package).
- **Testability-first Design:** Logic is decoupled from infrastructure, allowing for comprehensive unit and integration testing.
- **Minimal Coupling:** Layers only communicate with the layer immediately below them, and only through well-defined interfaces.

## Technology Stack
- **Backend:** Go (Golang)
- **Database:** PostgreSQL with GORM ORM
- **API Formats:** REST (standard endpoints) and GraphQL (flexible data querying)
- **API Documentation:** Swagger (OpenAPI 2.0)
- **Logging:** Structured logging with `log/slog`

## Documentation Standards
- **Swagger:** All REST API endpoints must be documented using Swagger annotations.
- **Generation:** Documentation is generated using `swag init -g cmd/main.go` from the `backend/` directory.
- **Access:** Interactive documentation is available at `/swagger/index.html` when the server is running.


## Coding Style Guidelines

Following the `example/stock-manage` pattern:

- **Project Structure:**
  - `cmd/main.go`: Application entry point.
  - `internal/application/`: Application-level orchestration and lifecycle management.
  - `internal/config/`: Configuration management.
  - `internal/db/`: Database connection and GORM initialization.
  - `internal/models/`: GORM domain models reflecting the database schema.
  - `internal/repository/`: Data access implementations.
  - `internal/services/`: Core business logic implementations.
  - `internal/routes/`: API routing and handler implementations.
- **Idiomatic Go:** Follow standard Go naming conventions and error handling patterns.
- **Constructor Injection:** Use `New...` functions for initializing structs with their dependencies.

## Logging Standards

Following the `example/GEMS` pattern:

- **Library:** Use the standard library `log/slog` for structured logging.
- **Outputs:** Logs should be written to both `os.Stdout` (for container logs) and `backend.log`.
- **Formatting:** Use JSON formatting for logs to ensure they are easily parsable by log management systems (like Loki).
- **Middleware:** Implement request-level logging middleware to capture HTTP method, path, status code, and latency.
- **Contextual Logging:** Always include relevant metadata in logs (e.g., `user_id`, `request_id`, `error` details).

## Data Management

- **Dataset:** `dataset.csv` serves as the initial data source.
- **Ingestion:** Implementation should include a mechanism to seed/sync the PostgreSQL database from the CSV file.
- **Migrations:** Use GORM's `AutoMigrate` or a dedicated migration tool for schema management.
