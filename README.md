# Backend Architecture (Go)

* **[Project description](https://github.com/alexander-kovalev2110/full-stack-web-proj_REST-API-BANK/blob/master/PHP-test.pdf)**
* **[Swagger (OpenAPI)](https://alexander-kovalev2110.github.io/full-stack-web-proj_REST-API-BANK/Swagger-OpenAPI/dist/index.html)**

This project is built with **Go + Chi Router + MySQL (`database/sql`) + JWT** and follows the principles of **Clean Architecture (Ports and Adapters)** and **Domain-Driven Design (DDD)**. It features a strict separation between HTTP presentation, application business logic, domain models, and infrastructure data persistence.

The architecture ensures predictable data flow:

**Command (Write) Flow:**
```
FE → Router (Chi) → Handler → UseCase → Repository Port → MySQL Adapter → JSON Response → FE
```

**Query (Read) Flow:**
```
FE → Router (Chi) → JWTMiddleware → Handler → UseCase → Repository Port → MySQL Adapter → TransactionListView DTO → JSON Response → FE
```

---

## Architectural Principles

* **Clean Architecture (Ports & Adapters)** — The core domain and application layers are isolated from external technical dependencies. They communicate with the outside world using abstract interfaces ("Ports"). Frameworks, databases, security, and external HTTP routers are treated as external "Adapters" located in the Infrastructure and Delivery layers.
* **Domain-Driven Design (DDD)** — The codebase is structured around business contexts (`Customer` and `Transaction`). Domain models protect business invariants, while repository interfaces abstract data persistence.
* **Separation of Concerns** — Each layer has a single responsibility.
* **Thin Handlers / Rich UseCases** — HTTP handlers only handle request decoding, response encoding, and basic HTTP validation; business rules and orchestration live in UseCase layers.
* **DTO-driven boundaries** — Request and response models (`TransactionView`, `TransactionListView`) are explicit and decoupled from raw database models.
* **Domain-first approach** — Application usecases operate on domain entities and interfaces, abstracting database infrastructure.

---

## Request Lifecycle Overview

### 1. Server Entry Point (`cmd/server/main.go`)
* **Purpose:** Bootstraps the application environment, database connection pool, repositories, use cases, HTTP handlers, Chi router middlewares, and starts the HTTP server.
* **Flow:** `Incoming HTTP Request → Chi Router`

### 2. Security & Middleware Layer
* **Purpose:** Handles CORS, request logging, panic recovery, and JWT authentication.
* **Responsibilities:**
  * [`JWTMiddleware`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/middleware/jwt_middleware.go) parses and validates JWT tokens from the `Authorization: Bearer <token>` header.
  * Injects the authenticated `customer_id` into Go's `context.Context` for downstream handlers.

### 3. DTO Layer (Request Models)
* **Purpose:** Strongly typed input definitions for request payloads.
* **Examples:**
  * `RegisterRequest` and `LoginRequest` in [`auth_handler.go`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/handler/auth_handler.go)
  * `AmountTransactionRequest` in [`transaction_handler.go`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/handler/transaction_handler.go)

### 4. Handler Layer (HTTP Orchestration)
* **Purpose:** Coordinates request handling without containing core business logic.
* **Responsibilities:**
  * Decodes JSON request payloads or query params (`page`, `limit`).
  * Extracts user context from `r.Context()`.
  * Delegates execution to UseCase services.
* **Structure:**
  * [`AuthHandler`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/handler/auth_handler.go) (`/customers/register`, `/customers/login`)
  * [`TransactionHandler`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/handler/transaction_handler.go) (`/transactions`, `/transactions/{transactionId}`)

### 5. UseCase / Application Layer
* **Purpose:** Implements application-level business logic and orchestration.
* **Responsibilities:**
  * Retrieve domain entities from repository interfaces.
  * Hash passwords via `bcrypt` and issue signed JWT tokens.
  * Enforce customer ownership boundaries for transaction access.
  * Construct immutable response DTOs (`TransactionListView`).
* **Structure:**
  * [`CustomerUseCase`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/usecase/customer_usecase.go)
  * [`TransactionUseCase`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/usecase/transaction_usecase.go)

### 6. Repository Layer (Data Access Ports & Adapters)
* **Purpose:** Encapsulate database retrieval and persistence.
* **Structure:**
  * **Domain Ports (Interfaces):** `CustomerRepository` in [`customer.go`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/domain/customer.go) and `TransactionRepository` in [`transaction.go`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/domain/transaction.go).
  * **Infrastructure Adapters:** MySQL implementations using `database/sql` in [`customer_repository.go`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/repository/mysql/customer_repository.go) and [`transaction_repository.go`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/repository/mysql/transaction_repository.go).

### 7. Domain Layer (Entities & Errors)
* **Purpose:** Core business models and domain error definitions.
* **Structure:**
  * [`internal/domain/customer.go`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/domain/customer.go)
  * [`internal/domain/transaction.go`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/domain/transaction.go)
* **Characteristics:**
  * Framework-agnostic.
  * Defines domain errors (`ErrCustomerAlreadyExists`, `ErrInvalidCredentials`, `ErrTransactionNotFound`).

### 8. Response Layer (Output View DTOs)
* **Purpose:** Defines structured JSON payloads returned by the API.
* **Examples:**
  * `TransactionView` (single transaction payload)
  * `TransactionListView` (paginated transaction list payload containing `transactions` array and `total` count)

---

## Exception Handling & Validation

1. **Request Payload Validation:**
   Input fields are checked for blank or non-positive values in handlers. Any violations trigger `middleware.WriteValidationError()` returning a structured `400 Bad Request` JSON response.
2. **Domain Error Mapping:**
   Domain errors returned by UseCases are intercepted by [`middleware.WriteError()`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/internal/middleware/error_handler.go#L26-L41) and mapped to standard HTTP statuses (`401 Unauthorized`, `404 Not Found`, `409 Conflict`, `500 Internal Server Error`).

---

## Folder Structure Overview

```
cmd/
  server/                   # Application entry point (main.go)
internal/
  domain/                   # Core domain models, DTOs, errors & repository interfaces (Ports)
  usecase/                  # Application business logic & use cases
  handler/                  # Primary (Driving) HTTP handlers / controllers
  repository/
    mysql/                  # Secondary (Driven) MySQL database repository adapters
  middleware/               # HTTP Middleware (JWT authentication, CORS, error handling)
```

---

## Running the Project & Tests

### Running the Server

Make sure MySQL is running and configured in [`.env`](file:///d:/KA/tests/github/REST-API-BANK_backend_go/.env).

```bash
go run ./cmd/server
```

### Running Tests

The test suite covers customer authentication and transaction logic using in-memory mock repositories.

#### On Linux / macOS / Git Bash:
* Run all tests:
  ```bash
  go test ./...
  ```
* Run with verbose output:
  ```bash
  go test -v ./...
  ```

#### On Windows (Command Prompt - CMD):
* Run all tests:
  ```cmd
  go test ./...
  ```

#### On Windows (PowerShell):
* Run all tests:
  ```powershell
  go test ./...
  ```

---

## Core Philosophy

Backend is responsible for:
* Enforcing business rules and protecting domain integrity.
* Validating all input parameter structures.
* Providing clean port-adapter boundaries for ease of testing.
* Returning structured contracts (`TransactionListView`) instead of raw database entities.
