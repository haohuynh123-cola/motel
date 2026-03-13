# GEMINI.md - Tro-Go API Context

This document provides essential context and instructions for AI agents working on the **Tro-Go API** project.

## 🚀 Project Overview

**Tro-Go API** is a high-performance Backend API built with **Go (Golang)** for managing motels and rooms. It follows **Clean Architecture (Hexagonal Architecture)** to ensure maintainability, scalability, and testability.

### Core Tech Stack
- **Language:** Go 1.24+
- **Web Framework:** [Echo v4](https://echo.labstack.com/)
- **Database:** PostgreSQL 15 (using `jackc/pgx/v5`)
- **Migrations:** `golang-migrate`
- **Authentication:** JWT (JSON Web Tokens)
- **Authorization:** RBAC (Role-Based Access Control)
- **DevOps:** Docker & Docker Compose with `Air` for hot-reloading.

## 🏗 Architecture & Directory Structure

The project is organized into layers to separate concerns:

- `cmd/api/`: Application entry point (`main.go`). Handles dependency injection and server startup.
- `db/migrations/`: SQL migration files for database schema versioning.
- `internal/`:
    - `domain/`: Business entities (structs) and core logic.
    - `port/`: Interface definitions for UseCases and Repositories (the "contracts").
    - `usecase/`: Business logic implementation. Orchestrates data flow between ports.
    - `adapter/`:
        - `handler/`: HTTP transport layer (Echo handlers).
        - `repository/`: Data access layer (PostgreSQL implementations).
        - `db/`: Database connection pool management.
- `pkg/`: Utility packages (e.g., `config` for environment variables).

## 🛠 Building and Running

### Development Environment (Docker)
The preferred way to run the project locally is via Docker Compose, which sets up the database, runs migrations, and starts the API with hot-reloading.

```bash
# Start the entire stack (Postgres + API with Air)
docker compose up

# Stop the stack
docker compose down
```

### Manual Commands (Local Go)
If running outside of Docker (requires local Postgres):
```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Run the API (requires .env configuration)
go run cmd/api/main.go
```

## 📝 Development Conventions

- **Clean Architecture:** Always define interfaces in `internal/port/` before implementing them in `usecase/` or `adapter/`.
- **Dependency Injection:** Dependencies are injected in `cmd/api/main.go`. Avoid global state or `init()` functions for database connections or configurations.
- **Error Handling:** Use Go's idiomatic error handling. Prefer wrapping errors with context if necessary, but keep domain errors (like `ErrNotFound`) in `internal/port/`.
- **Database Access:** Use raw SQL with `pgx` in the repository layer. Avoid ORMs to maintain control over performance and complexity.
- **Hot-Reloading:** The development environment uses `Air`. Configuration is in `.air.toml`.
- **Environment Variables:** All configurations should be handled via `pkg/config` and read from `.env` or system environment variables.

## 🧪 Testing Strategy

- **Unit Tests:** Located alongside the implementation files (e.g., `internal/usecase/user_usecase_test.go`).
- **Mocking:** Use the interfaces in `internal/port/` to create mocks for testing UseCases in isolation.
- **Database Tests:** Use a test database or transactional tests for repository layer validation.

---
*This file is intended for AI context. For human-readable documentation, see `README.md`.*
