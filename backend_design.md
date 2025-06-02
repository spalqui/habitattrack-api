# Backend Code Structure and Design

This document outlines the proposed backend code structure and design for the HabitatTrack Financial Management API. It emphasizes SOLID principles, modularity, maintainability, and testability.

## 1. Layered Architecture

We will adopt a layered architecture to separate concerns and promote modularity. The primary layers will be:

*   **Presentation/API Layer (Handlers/Controllers):**
    *   **Responsibilities:** Receives HTTP requests, validates them (superficially, e.g., path/query parameters, request body format), deserializes request bodies into Data Transfer Objects (DTOs) or request models, calls appropriate methods in the Application/Service Layer, serializes responses (or error responses) from the Service Layer into HTTP responses.
    *   **Focus:** HTTP communication, request/response lifecycle, authentication/authorization hooks.
    *   **Example:** `TransactionHandler` would parse an incoming JSON for creating a transaction, call `TransactionService.CreateTransaction()`, and then format the returned transaction (or error) as a JSON HTTP response.

*   **Application/Service Layer:**
    *   **Responsibilities:** Contains the core business logic and use cases of the application. Orchestrates calls to domain objects and repositories. Handles transaction management (database transactions), validation of business rules, and coordination between different domain entities if necessary. Maps DTOs from the API layer to domain models and vice-versa.
    *   **Focus:** Business workflows, use case implementation, data transformation.
    *   **Example:** `TransactionService` would contain methods like `CreateTransaction()`, `GetTransactionByID()`, `ListTransactions()`. `CreateTransaction()` would validate the input against business rules (e.g., ensuring `categoryId` matches the transaction `type`), interact with `TransactionRepository` and potentially `CategoryRepository` or `PropertyRepository`.

*   **Domain Layer:**
    *   **Responsibilities:** Represents the core business concepts, entities, value objects, and domain-specific logic. Entities encapsulate data and behavior related to that data. This layer should be independent of infrastructure concerns.
    *   **Focus:** Business entities, rules directly associated with those entities, and their states.
    *   **Example:** `Transaction` entity (with fields like `ID`, `Amount`, `TransactionDate`, `Type`, `CategoryID`, `PropertyID`, and methods if any specific to a transaction's lifecycle), `TransactionCategory` entity, `Property` entity.

*   **Infrastructure Layer:**
    *   **Responsibilities:** Deals with external concerns such as database interactions, file system access, network communication (e.g., calling external APIs), message queues, caching, and logging implementations. Implements interfaces defined in the Application or Domain layers (e.g., Repository interfaces).
    *   **Focus:** Data persistence, external service integration, low-level technical details.
    *   **Example:** `PostgresTransactionRepository` (implementing a `TransactionRepository` interface defined in the Application/Service layer) would handle SQL queries to a PostgreSQL database for transaction data. `S3FileStorage` for storing receipts, `StdOutLogger` for logging.

## 2. SOLID Principles Application

The design will adhere to the SOLID principles to ensure a robust, maintainable, and flexible codebase.

*   **Single Responsibility Principle (SRP):**
    *   **Application:**
        *   **Handlers/Controllers:** Solely responsible for handling HTTP request/response concerns (parsing, serialization, routing to services). They will not contain business logic.
        *   **Services:** Each service (e.g., `TransactionService`, `CategoryService`, `PropertyService`) will be responsible for a specific domain's business logic and use cases. For instance, `TransactionService` handles all operations related to transactions but doesn't know about user authentication details (handled by middleware or a dedicated auth service).
        *   **Repositories:** Each repository (e.g., `TransactionRepository`) will be responsible only for data access and persistence logic for its specific entity.
        *   **Domain Models:** Entities like `Transaction` will encapsulate data and behavior specific to that entity, not how it's stored or presented.

*   **Open/Closed Principle (OCP):**
    *   **Application:** The system will be designed to be open for extension but closed for modification.
        *   **New Validation Rules:** Transaction validation logic within `TransactionService` could use a strategy pattern. New validation rules (e.g., for specific transaction types or amounts) can be added as new strategy implementations without modifying the core service code.
        *   **Reporting Capabilities:** If new reporting formats or data sources are needed, new `ReportGenerator` implementations could be added, adhering to a common `IReportGenerator` interface, without altering existing reporting services.
        *   **Payment Gateways/Notification Services:** Abstract external service integrations (like payment processing or notifications) behind interfaces. New providers can be added by implementing these interfaces.

*   **Liskov Substitution Principle (LSP):**
    *   **Application:** While direct inheritance hierarchies for core domain entities might be minimal in this Go project (favoring composition), if interfaces are used to define common behaviors (e.g., an `Auditable` interface with `GetCreatedAt()` and `GetUpdatedAt()` methods implemented by `Transaction`, `Category`, etc.), any implementation of such an interface must be substitutable where the interface is expected.
    *   If, for example, we had different types of `NotificationService` (e.g., `EmailNotificationService`, `SMSNotificationService`) implementing an `INotificationService` interface, any instance should be usable wherever an `INotificationService` is required, without unexpected behavior.

*   **Interface Segregation Principle (ISP):**
    *   **Application:** Interfaces will be lean and focused. Clients (e.g., a service using a repository) should not be forced to depend on methods they do not use.
        *   **Repository Interfaces:** Instead of one large `IRepository` interface, we'll have specific interfaces like `ITransactionRepository` (with methods like `Create`, `GetByID`, `List`, `Update`, `Delete` for transactions) and `ICategoryRepository` for categories. A service that only reads transactions would only need to depend on a part of the interface (e.g., `ITransactionReader`) if we further segregate read/write operations.
        *   **Service Interfaces:** Services themselves will expose focused interfaces to the Presentation layer. For example, `ITransactionService` will only contain methods relevant to transaction management.

*   **Dependency Inversion Principle (DIP):**
    *   **Application:** High-level modules (e.g., Services) will depend on abstractions (interfaces), not on low-level concrete implementations (e.g., specific database repositories). Low-level modules will also depend on abstractions.
        *   **Dependency Injection (DI):** Concrete implementations (like `PostgresTransactionRepository`) will be "injected" into services (like `TransactionService`) that depend on the `ITransactionRepository` interface. This will be managed using a DI container or manual DI at the application's composition root (e.g., in `main.go`).
        *   `TransactionService` (high-level) -> `ITransactionRepository` (abstraction) <- `PostgresTransactionRepository` (low-level detail).

## 3. Modular Directory Structure (Organized by Feature/Domain)

A modular directory structure, primarily organized by feature/domain, will be used to enhance clarity and maintainability.

```
/habitattrack-api
|-- /cmd
|   |-- /api                 # Main application entry point
|   |   |-- main.go
|-- /internal                # Private application and library code
|   |-- /config              # Configuration loading and management
|   |   |-- config.go
|   |-- /core                # Core domain entities, value objects - potentially shared across features
|   |   |-- /entity
|   |   |   |-- transaction.go
|   |   |   |-- category.go
|   |   |   |-- property.go
|   |   |-- /valueobject     # e.g., Money, DateRange
|   |-- /platform            # Infrastructure implementations (database, logging, etc.)
|   |   |-- /database        # Database connection, migrations
|   |   |   |-- postgres.go
|   |   |-- /logger          # Logging setup
|   |   |   |-- zap_logger.go
|   |   |-- /auth            # Auth client, token validation logic
|   |-- /features            # Feature-specific modules
|   |   |-- /transactions
|   |   |   |-- handler.go       # HTTP handlers (Presentation Layer)
|   |   |   |-- service.go       # Business logic (Application Layer)
|   |   |   |-- repository.go    # Repository interface (Application Layer)
|   |   |   |-- postgres_repo.go # PostgreSQL impl of repository (Infrastructure Layer, specific to this feature)
|   |   |   |-- dto.go           # Data Transfer Objects for API requests/responses
|   |   |   |-- validation.go    # Feature-specific validation rules
|   |   |-- /categories
|   |   |   |-- handler.go
|   |   |   |-- service.go
|   |   |   |-- repository.go
|   |   |   |-- postgres_repo.go
|   |   |   |-- dto.go
|   |   |-- /properties
|   |   |   |-- handler.go
|   |   |   |-- service.go
|   |   |   |-- repository.go    # (Assuming existing property structure might be refactored or integrated)
|   |   |   |-- postgres_repo.go
|   |   |   |-- dto.go
|   |-- /shared              # Shared utilities, DTOs, interfaces if truly cross-cutting
|   |   |-- /middleware      # HTTP middleware (auth, logging, CORS)
|   |   |-- /utils           # Common utility functions
|   |   |-- /apierrors       # Standardized API error types/responses
|-- /api                     # OpenAPI specification files
|   |-- openapi.yaml
|-- /pkg                     # Public library code (if any, to be shared with other projects)
|-- /scripts                 # Build, deploy, utility scripts
|-- /test                    # E2E tests, shared test utilities
|   |-- /e2e
|   |-- /mocks
|-- go.mod
|-- go.sum
|-- README.md
|-- backend_design.md      # This file
```

**Explanation of Key Directories:**

*   **`cmd/api/main.go`:** The entry point of the application. Responsible for initializing dependencies (config, logger, database connections, repositories, services, handlers), setting up HTTP routes, and starting the server. This is the "composition root".
*   **`internal/config`:** Handles loading application configuration from environment variables, config files, etc.
*   **`internal/core/entity`:** Contains the primary domain model definitions (structs for `Transaction`, `TransactionCategory`, `Property`). These are plain Go structs, potentially with methods that enforce business rules.
*   **`internal/core/valueobject`:** Contains value objects like `Money` (to handle currency and amount precisely) or `DateRange`.
*   **`internal/platform/*`:** Contains concrete implementations for infrastructure concerns. For example, `database` might contain GORM setup or raw `sql.DB` helpers, and `logger` might configure a logging library like Zap or Logrus.
*   **`internal/features/*`:** Each feature (transactions, categories, properties) gets its own module.
    *   `handler.go`: Implements the Presentation/API Layer for that feature.
    *   `service.go`: Implements the Application/Service Layer.
    *   `repository.go`: Defines the repository interface(s) for the feature (dependency of the service).
    *   `postgres_repo.go` (or similar): Implements the repository interface using a specific database (Infrastructure Layer, but co-located for feature cohesion or placed in `platform/persistence/transaction_repo.go` etc.).
    *   `dto.go`: Data Transfer Objects specific to this feature's API.
*   **`internal/shared/*`:** For code that is genuinely shared across multiple features and doesn't belong to a specific one, like HTTP middleware, common utility functions, or standardized API error structures.
*   **`api/`:** Contains the OpenAPI specification.

This structure promotes modularity by feature, making it easier to navigate, develop, and test individual parts of the application.

## 4. Maintainable & Testable Design

*   **Error Handling:**
    *   **Standardized Error Responses:** The API will return standardized JSON error responses as defined in `openapi.yaml` (`ErrorResponse` schema with `code`, `message`, `details`).
    *   **Custom Error Types:** Define custom error types in Go (e.g., `ErrNotFound`, `ErrValidation`, `ErrUnauthorized`, `ErrDuplicateEntry`) that services can return.
    *   **Error Wrapping:** Use Go 1.13+ error wrapping (`fmt.Errorf("service: %w", err)`) to provide context without losing the original error type.
    *   **Centralized Error Handling Middleware:** An HTTP middleware can catch errors returned by services, log them appropriately, and convert them into the standard API error JSON response with correct HTTP status codes.
        *   `ErrNotFound` -> HTTP 404
        *   `ErrValidation` -> HTTP 400 or 422
        *   `ErrUnauthorized` -> HTTP 401
        *   `ErrDuplicateEntry` -> HTTP 409
        *   Unhandled/unexpected errors -> HTTP 500 (with a generic message to the client and detailed log server-side).

*   **Logging Strategy:**
    *   **Structured Logging:** Use a structured logging library (e.g., Zap, Logrus) to produce logs in JSON or another machine-parsable format.
    *   **Log Levels:** Implement standard log levels (DEBUG, INFO, WARN, ERROR, FATAL).
    *   **Contextual Logging:** Include contextual information in logs, such as request ID, user ID (if applicable), operation name, etc. This can be achieved by injecting a logger instance with context into services or using context-aware loggers.
    *   **What to Log:**
        *   **Requests:** Key details of incoming requests (method, path, but not sensitive data in body unless specifically for debugging in dev).
        *   **Errors:** All errors, especially at the point they are handled or cross layer boundaries. Include stack traces for unexpected errors.
        *   **Key Business Events:** Important state changes or business operations (e.g., "Transaction created: ID=xyz").
        *   **External Service Calls:** Requests and responses (or errors) when interacting with external services.
    *   **Audit Logging:** For sensitive operations, dedicated audit logs might be necessary, capturing who did what and when.

*   **Separation of Concerns:**
    *   Strict adherence to the layered architecture and SRP ensures components have distinct responsibilities, reducing coupling and making the system easier to understand and modify.
    *   DI further promotes this by decoupling components from their concrete dependencies.

*   **Testing Strategy:**
    *   **Unit Tests:**
        *   **Coverage:** Business logic in services, complex logic in domain models, utility functions, validation rules.
        *   **Technique:** Test individual functions/methods in isolation. Use mocks/stubs for dependencies (e.g., mock repository interfaces when testing services). Go's built-in `testing` package is standard. Table-driven tests are encouraged for covering multiple scenarios.
        *   **Location:** `_test.go` files alongside the code being tested (e.g., `service_test.go` in the `transactions` feature directory).
        *   **Tools:** Go's `testing` package, `testify/mock` for mocking, `testify/assert` or `testify/require` for assertions.
    *   **Integration Tests:**
        *   **Coverage:** Interactions between layers, typically service-repository interactions, and API endpoint behavior against a real (test) database.
        *   **Technique:** Test a component with its actual dependencies (e.g., test a service method that calls a real repository interacting with a test database). For API integration tests, make HTTP calls to the running application (or a test server instance) and verify responses, status codes, and database state changes.
        *   **Location:** Can be in `_test.go` files (e.g., `handler_integration_test.go`) or a separate `test/integration` directory.
        *   **Tools:** Go's `testing` package, `net/http/httptest` for API tests, Docker for spinning up test databases (e.g., using `testcontainers-go`).
    *   **End-to-End (E2E) Tests:**
        *   **Coverage:** Test complete user flows through the deployed application, simulating real user scenarios.
        *   **Technique:** Interact with the API as a client would, verifying the entire system behavior. These are typically slower and run less frequently (e.g., in a CI/CD pipeline against a staging environment).
        *   **Location:** `/test/e2e` directory.
        *   **Tools:** Go's `testing` package with HTTP client, or dedicated E2E testing frameworks if UI is involved (not applicable here).

*   **Configuration Management:**
    *   **Environment-Specific Configurations:** The application must support different configurations for various environments (development, staging, production).
    *   **Source:** Configuration can be loaded from environment variables (primary, for 12-factor app compliance), and optionally from configuration files (e.g., `.env`, `config.yaml`) for local development convenience.
    *   **Structure:** A `Config` struct (e.g., in `internal/config/config.go`) will hold all configuration parameters (database DSN, server port, JWT secrets, external service URLs, log level).
    *   **Loading:** Use a library like `spf13/viper` or a custom solution to load and unmarshal configuration at startup.
    *   **Secrets Management:** Sensitive data (API keys, database passwords) should be managed securely, ideally through environment variables injected by the deployment platform or a secrets management tool (e.g., HashiCorp Vault, AWS Secrets Manager), not hardcoded or committed to version control.

This comprehensive approach to backend design aims to create a system that is scalable, maintainable, and easy to test, aligning with modern software engineering best practices.