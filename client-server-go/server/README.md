# TCP File Server in Go

A hand-rolled HTTP file server built on raw TCP sockets — no `net/http`. Written to understand how HTTP requests are parsed and handled at the network level, following [Beej's Guide to Network Programming](https://beej.us/guide/bgnet/).

## What it does

For each incoming TCP connection, the server:

1. **Parses the HTTP request** — reads the request line, headers, and body from the raw byte stream
2. **Resolves the file path** — strips path traversal sequences (`../`) so all paths are rooted at the server's configured directory
3. **Reads the file and determines its MIME type** — infers `Content-Type` from the file extension (`.html`, `.css`, `.js`, `.json`, `.png`, `.jpg`, plain text)
4. **Builds and sends the HTTP response** — writes the status line, headers, and file content back over the connection
5. **Closes the connection** — one request per connection (`Connection: close`)

Graceful shutdown is handled via OS interrupt signals — in-flight connections finish before the server exits.

## Project structure

```
cmd/server/
  main.go               Entry point; wires router, starts server, handles shutdown

internal/
  models/
    request.go          Request struct (method, path, HTTP version, headers, body)
    response.go         Response struct

  parser/
    parser.go           Reads from io.Reader; parses request line, headers, body

  errors/
    errors.go           HTTPError type with constructors (NewBadRequest, NewNotFound, etc.)

  response/
    response.go         ResponseWriter interface (Header, Write)
    conn_writer.go      net.Conn-backed implementation; serialises HTTP response to wire format

  handlers/
    handler.go          Handler interface (ServeHTTP)
    file_handler.go     Reads file from disk, sets Content-Type, writes response

  router/
    router.go           Maps (method, path) → Handler; 404 on no match

  server/
    server.go           TCP listener; accepts connections; context-aware graceful shutdown
```

## Running the server

```bash
# Build
go build -o fileserver ./cmd/server

# Run (serves files from ./cmd/server on port 28333)
./fileserver --port 28333 --dir ./cmd/server

# Request a file
curl http://localhost:28333/file1.txt
curl http://localhost:28333/file2.html

# Stop with Ctrl+C — in-flight connections drain before exit
```

## Running the tests

```bash
go test ./...
```

Tests cover:

| Package | What's tested |
|---|---|
| `parser` | Valid GET, POST with body, malformed request line, malformed headers, invalid `Content-Length` |
| `handlers` | File served with correct body and `Content-Type`, 404 for missing file, path traversal blocked |
| `router` | Route match, 404 for unregistered path, 404 on method mismatch |

No `net.Conn` or running server required for any test — the `ResponseWriter` interface lets handlers be tested with an in-memory mock.

## Adding a route

Register a handler in `main.go` before starting the server:

```go
r := router.NewRouter()
r.Handle("GET", "/hello.txt", &handlers.FileHandler{RootDir: "./static"})
```
