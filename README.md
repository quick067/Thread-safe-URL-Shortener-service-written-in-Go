# Thread-Safe URL Shortener Service

A production-grade URL shortener written in **Go (Golang)**.
This project focuses on **clean architecture**, **concurrency safety**, **security**, and **observability**, demonstrating how to build a maintainable and scalable backend service using mostly the Go standard library and a small set of well-chosen dependencies.

---

## Overview

The service provides a REST API for shortening URLs and redirecting users.
It is designed to be safe under high concurrent load, easy to test, and ready for real-world usage with authentication, monitoring, and persistent storage.

Key design goals:

* Clear separation of concerns
* Predictable behavior under concurrency
* Strong emphasis on correctness and testability
* Production-oriented features (metrics, rate limiting, graceful failure)

---

## Key Features

### Security and Authentication

* **JWT Authentication**
  Secure access to protected endpoints using JSON Web Tokens (access tokens).

* **Password Hashing**
  User passwords are hashed with `bcrypt` before being stored in the database.

* **Rate Limiting**
  IP-based rate limiting using the Token Bucket algorithm to mitigate abuse and denial-of-service attacks.

* **Input Validation**
  Strict validation of URLs and custom aliases to prevent invalid or malicious input.

---

### Observability and Monitoring

* **Prometheus Metrics**
  Built-in Prometheus integration via the `/metrics` endpoint, exposing:

  * Total number of HTTP requests
  * Response status code counters (2xx, 4xx, 5xx)
  * Request duration / latency

* **Structured Logging**
  HTTP middleware logs request method, path, status, duration, and remote IP.

---

### Storage and Data Management

* **PostgreSQL Backend**
  Persistent storage for users and URL mappings.

* **Concurrency Safety**
  All operations are safe under concurrent access, ensuring data consistency.

* **Conflict Handling**
  Graceful handling of duplicate aliases with proper HTTP 409 (Conflict) responses.

---

### Configuration and Architecture

* **Flexible Configuration**
  Configuration can be provided via:

  * Environment variables
  * Command-line flags
    Priority order: Environment Variables → Flags → Defaults.

* **Clean Architecture**
  Clear separation between layers:

  * `handlers` — HTTP layer
  * `service` — business logic
  * `store` — data access layer

* **Dependency Injection**
  Services and storage implementations are injected, enabling easy testing and modularity.

---

### Testing

* **High Test Coverage**
  Unit tests cover handlers, services, and middleware.

* **Mock-Based Testing**
  Database interactions are mocked to isolate business logic during tests.

* **Edge Case Coverage**
  Tests include validation failures, conflicts, and internal error scenarios.

---

## Technology Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **HTTP:** `net/http`
* **Authentication:** `golang-jwt/jwt/v5`
* **Metrics:** `prometheus/client_golang`
* **Rate Limiting:** `golang.org/x/time/rate`
* **Configuration:** `flag`, `os`

---

## Getting Started

### Prerequisites

* Go 1.22 or newer
* PostgreSQL instance

---

### Installation

Clone the repository:

```bash
git clone https://github.com/quick067/Thread-safe-URL-Shortener-service-written-in-Go.git
cd Thread-safe-URL-Shortener-service-written-in-Go
```

---

### Configuration

Set environment variables (or use a `.env` file):

```bash
export SERVER_ADDRESS=":8080"
export BASE_URL="http://localhost:8080"
export DATABASE_DSN="postgres://user:password@localhost:5432/shortener?sslmode=disable"
export JWT_SECRET="your_super_secret_key"
```

---

### Running the Service

Using environment variables:

```bash
go run cmd/shortener/main.go
```

Using command-line flags:

```bash
go run cmd/shortener/main.go -a :9090 -b http://my-shortener.com
```

---

## API Endpoints

| Method | Endpoint    | Auth Required | Description                          |
| -----: | ----------- | ------------- | ------------------------------------ |
|   POST | `/register` | No            | Create a new user account            |
|   POST | `/login`    | No            | Authenticate and receive a JWT token |
|   POST | `/save`     | Yes           | Shorten a URL (JSON body)            |
|    GET | `/{alias}`  | No            | Redirect to the original URL         |
|    GET | `/metrics`  | No            | Prometheus metrics endpoint          |

---

## Testing

Run all tests:

```bash
go test ./...
```

Check test coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Project Structure

```text
.
├── cmd
│   └── shortener
│       └── main.go        # Application entry point
├── internal
│   ├── config             # Configuration (env + flags)
│   ├── handlers           # HTTP handlers
│   ├── middleware         # Auth, logging, rate limiting, metrics
│   ├── service            # Business logic
│   ├── store              # PostgreSQL data access layer
│   └── mocks              # Mocks for unit testing
├── go.mod
└── README.md
```

---

## Author

Developed by **quick067** as part of an advanced Go backend engineering learning path.
