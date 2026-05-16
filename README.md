# Go-Play 🚀

![Build Status](https://github.com/chetanuchiha16/go-play/actions/workflows/ci.yml/badge.svg)

A professional, production-ready Go template focusing on **Pragmatic Clean Architecture**, type-safe database interactions, and robust development workflows.

## 🏗 Architecture

The project follows a **Feature-Based (Domain-Driven)** organization combined with **Pragmatic Clean Architecture** principles.

### The "Aura" Hierarchy (Interface vs Struct)
To achieve decoupling and high testability, the project uses a deliberate pattern for layers:

| Layer | Type | Responsibility | Clean Architecture Role |
| :--- | :--- | :--- | :--- |
| **Handler** | `Struct` | Manages HTTP transport, OpenAPI mapping, and validation. | Interface Adapter |
| **Service** | `Interface` | Encapsulates business logic (e.g., hashing, orchestration). | Use Cases |
| **Store** | `Interface` | Abstract data persistence (Repository pattern). | Interface Adapter |
| **Models** | `Generated` | Type-safe entities generated via SQLC. | Entities (Shared) |

> [!TIP]
> This "Interface vs Struct" pattern allows the Handler to be unit-tested using a Mock Service, and the Service to be tested using a Mock Store, without needing a real database or network.

## 🛠 Tech Stack

*   **Language**: Go 1.25.5
*   **Database**: PostgreSQL with `pgx/v5`
*   **API Framework**: Standard `net/http` with OpenAPI generation (`oapi-codegen`)
*   **Code Generation**: [SQLC](https://sqlc.dev/) (Type-safe SQL), [Mockery](https://github.com/vektra/mockery) (Testing)
*   **Documentation**: Swagger UI via OpenAPI 3.0
*   **Configuration**: [Viper](https://github.com/spf13/viper) (Env/Config management)
*   **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
*   **Logging**: [Zerolog](https://github.com/rs/zerolog) (Structured logging)

## 📁 Project Structure

*   **`cmd/server/`**: Application entry point and router initialization.
*   **`api/`**: OpenAPI specifications (`openapi.yaml`) and generator configs.
*   **`db/`**: Migrations, SQL queries, and SQLC-generated code.
*   **`internal/domain/`**: Feature-based core logic (User, Auth, etc.).
*   **`internal/api/`**: Generated OpenAPI server stubs and types.
*   **`internal/middleware/`**: Custom HTTP middlewares (Logging, Auth).
*   **`tests/`**: Integration tests and test suites.

## ⚙️ Getting Started

### Prerequisites
*   Go (v1.25.5+)
*   Docker & Docker Compose
*   [sqlc](https://sqlc.dev/) & [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) (for generation)

### Setup
1.  **Clone the repository**.
2.  **Environment**: Copy `.env` and configure your `DATABASE_URL`.
3.  **Spin up Infrastructure**:
    ```bash
    make up
    ```

## 🛠 Development Commands

| Command | Description |
| :--- | :--- |
| `make dev` | Start local development with hot-reloading (**Air**) |
| `make generate` | Generate OpenAPI server stubs and types |
| `make sqlc` | Generate type-safe Go code from SQL queries |
| `make up` | Start DB and App using Docker Compose |
| `make test` | Run all unit tests |
| `make test-integration` | Run integration tests (requires DB) |
| `make mock` | Generate mocks using Mockery |
| `make test-cover` | Generate and view HTML coverage report |

## 🚀 API Documentation
The API is documented using OpenAPI. You can view the Swagger UI by navigating to `/docs` (when the server is running).

To update the API:
1.  Modify `api/openapi.yaml`.
2.  Run `make generate` to update the server stubs.

## 🧪 Testing
The project distinguishes between unit and integration tests:
*   **Unit Tests**: Located alongside the code (e.g., `service_test.go`). Run with `make test`.
*   **Integration Tests**: Located in `tests/integration/`. They verify the full HTTP -> Database flow. Run with `make test-integration`.

---