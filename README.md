# Go URL Shortener

A simple, thread-safe URL shortener service written in Go.

## Features

- **Shorten URLs**: Converts long URLs into short keys.
- **Redirection**: Redirects users to the original URL.
- **Thread-Safe**: Uses `sync.RWMutex` to handle concurrent requests safely.
- **Persistence**: Saves data to a JSON file (`storage.json`) to survive server restarts.
- **Testing**: Includes unit tests for handlers.

## Usage

### 1. Run the server
```bash
go run .
The server will start on http://localhost:8080.

2. Save a URL
Send a POST request to /save with the URL in the body:

Bash

curl -X POST -d "[https://google.com](https://google.com)" http://localhost:8080/save
Response: http://localhost:8080/AbCdEfGh

3. Use the short link
Paste the received short URL into your browser, and you will be redirected to the original site.

Project Structure
main.go - Server entry point and setup.

handlers.go - HTTP handlers (saveURL, redirectURL).

structs.go - Data structures and file I/O logic.

utils.go - Key generator.

handlers_test.go - Unit tests.

Technologies
Go (Golang)

Standard library (net/http, encoding/json, sync)