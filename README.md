# bfm-example

Go module **`bfm-example`** (see [`go.mod`](go.mod)).

**bfmcp** means **Backend for Model Context Protocol**. This repo is a small **Backend-for-Frontend (BFF)** proof of concept in Go: the BFF exposes a single workflow—**process order**—as a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) tool over HTTP. It validates input, checks that the user exists via the backend, forwards the order to the backend, then enriches the result (user details and BFF metadata).

The codebase follows a layered layout: domain (entities, repositories, use cases), infrastructure (HTTP clients with header forwarding), and adapters (MCP handlers, middleware, transport). Wiring uses [Uber Fx](https://uber-go.github.io/fx/).

## Requirements

- [Go](https://go.dev/dl/) **1.25.1** or compatible (see `go.mod`).

## Configuration

The server reads environment variables (via [Viper](https://github.com/spf13/viper)):

| Variable | Purpose | Default |
|----------|---------|---------|
| `PORT` | Address the BFF HTTP server listens on | `:8081` |
| `BACKEND_BASE_URL` | Base URL of the REST backend (no trailing slash required) | `http://localhost:8082` |

The BFF serves MCP over **streamable HTTP** at **`/mcp`** on the configured port.

## Run the BFF server

From the repository root:

```bash
go run ./cmd/server
```

With custom port and backend:

```bash
PORT=:9090 BACKEND_BASE_URL=http://localhost:3000 go run ./cmd/server
```

Keep this process running while you use a client or the demo below.

## Run the end-to-end demo

The repository root `main.go` is a **demo client**: it starts a **mock REST API** on port **8082** (`POST /api/process_order`, `GET /api/users/{userID}`) and connects to the BFF as an MCP client, calling the `process_order` tool with sample arguments and tenant/correlation headers.

1. **Terminal 1** — start the BFF (expects the backend at `http://localhost:8082` by default):

   ```bash
   go run ./cmd/server
   ```

2. **Terminal 2** — start the mock backend and run the demo client against the BFF at `http://localhost:8081/mcp`:

   ```bash
   go run .
   ```

You should see JSON printed for the structured tool result.

## Tests

```bash
go test ./...
```

## Project layout (high level)

- `cmd/server` — Fx-wired MCP HTTP server (`main.go`).
- `internal/config` — Loads `PORT` and `BACKEND_BASE_URL` (Viper).
- `internal/domain` — Entities, repository ports, and use cases (e.g. `ProcessOrder`).
- `internal/infrastructure/httpclient` — HTTP implementations of repositories.
- `internal/adapter` — MCP tool handlers, logging middleware, and transport registration.
- `pkg/headerforward` — Allowlisted headers from context into outbound HTTP and MCP client calls.
- `pkg/httpjson` — Bounded JSON decoding from HTTP responses.
- `pkg/validate` — Shared [go-playground/validator](https://github.com/go-playground/validator) setup (e.g. custom tags).
- `main.go` (root) — Mock backend on **:8082** + MCP demo client targeting the BFF (not the production server entrypoint).
