# Go URL Shortener Service (v1.2)

A robust, thread-safe URL shortener service written in Go.
This project demonstrates **Clean Architecture**, **Concurrent Programming**, and **Graceful Shutdown** mechanisms without relying on external frameworks.

---

## 🚀 Features

* **REST API** – Simple endpoints to shorten URLs and redirect users
* **Layered Architecture** – Clean separation into `handlers`, `store`, `config`, `server`
* **Thread-Safety** – Concurrent access protected via `sync.RWMutex`
* **Data Persistence** – URL mappings are saved to disk in JSON format
* **Graceful Shutdown** – Handles `SIGINT` / `SIGTERM` to safely persist data
* **Configuration via Environment Variables**

---

## 🛠 Project Structure

```text
.
├── internal/
│   ├── config/       # Environment variable management
│   ├── handlers/     # HTTP request handlers & business logic
│   ├── server/       # HTTP server setup and lifecycle management
│   └── store/        # Data storage logic (thread-safe map + file I/O)
├── main.go           # Application entry point (Dependency Injection)
├── go.mod            # Go module definition
└── README.md         # Documentation
```

---

## ⚙️ Configuration

The application can be configured using environment variables.

| Variable      | Description                           | Default Value  |
| ------------- | ------------------------------------- | -------------- |
| `SERVER_PORT` | Port on which the HTTP server runs    | `:8080`        |
| `FILENAME`    | JSON file used for persistent storage | `storage.json` |

---

## 📦 Getting Started

### Prerequisites

* Go **1.20+** installed

### Installation

```bash
git clone https://github.com/your-username/your-repo-name.git
cd your-repo-name
```

---

## ▶️ Running the Application

### Default Run

```bash
go run .
```

### Custom Configuration (Linux / macOS)

```bash
SERVER_PORT=":9090" FILENAME="my_db.json" go run .
```

### Custom Configuration (Windows PowerShell)

```powershell
$env:SERVER_PORT=":9090"
$env:FILENAME="my_db.json"
go run .
```

---

## 🧪 Running Tests

```bash
go test ./... -v
```

---

## 🔌 API Usage

### 1. Save a URL

**Endpoint:**
`POST /save`

**Body:** raw string containing the URL

```bash
curl -X POST -d "https://www.google.com" http://localhost:8080/save
```

**Response:**

```
http://localhost:8080/AbCdEfGh
```

---

### 2. Redirect

**Endpoint:**
`GET /{short_key}`

Open the generated short URL in your browser or test via curl:

```bash
curl -v http://localhost:8080/AbCdEfGh
```

---

## 🛡 Graceful Shutdown

To test safe shutdown behavior:

1. Run the server
2. Create several short URLs
3. Press **Ctrl + C**

The application will:

* Intercept the shutdown signal
* Save all in-memory data to disk
* Exit cleanly without data loss

---

## 💻 Technologies Used

* **Language:** Go (Golang)
* **Standard Library Only**

  * `net/http` — HTTP server
  * `sync` — concurrency primitives
  * `encoding/json` — serialization
  * `os/signal`, `context` — graceful shutdown handling

---

✅ **Result:** A clean, idiomatic, production-style Go service demonstrating concurrency, architecture, and reliability.
