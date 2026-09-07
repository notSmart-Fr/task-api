# Task API

A high-performance, in-memory REST API for managing a to-do list. The service is built with Go 1.22+ and the standard library's `net/http` package, and is organized using Vertical Slice Architecture.

## Features and Architecture

- **Zero framework overhead:** Uses Go's standard library and enhanced Go 1.22 path routing.
- **Vertical Slice Architecture:** The `tasks`, `health`, and `docs` features are isolated into self-contained modules.
- **Thread-safe in-memory store:** Uses `sync.RWMutex` for safe concurrent access.
- **Structured logging:** Request metadata is emitted as JSON using `log/slog`.
- **Interactive API reference:** Scalar serves the OpenAPI documentation at `/docs`.

## Getting Started

### Prerequisites

- Go 1.22 or later

### Run the server

From the repository root:

```bash
go run ./cmd/server
```

The API listens on `http://localhost:8000`.

## API Endpoints

| Method | Path | Description | Responses |
| --- | --- | --- | --- |
| `GET` | `/` | API metadata and endpoint index | `200 OK` |
| `GET` | `/health` | Server health check | `200 OK` |
| `GET` | `/docs` | Scalar interactive API documentation | `200 OK` |
| `GET` | `/tasks` | List all tasks | `200 OK` |
| `GET` | `/tasks/{id}` | Get a task by ID | `200 OK`, `400 Bad Request`, `404 Not Found` |
| `POST` | `/tasks` | Create a task | `201 Created`, `400 Bad Request` |
| `PUT` | `/tasks/{id}` | Update a task title and/or completion state | `200 OK`, `400 Bad Request`, `404 Not Found` |
| `DELETE` | `/tasks/{id}` | Delete a task | `204 No Content`, `400 Bad Request`, `404 Not Found` |

## Sample Request

Create a task with PowerShell:

```powershell
$body = @{ title = "Finish backend assignment" } | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8000/tasks" -Method Post -Body $body -ContentType "application/json"
```

Example response:

```json
{
  "id": 4,
  "title": "Finish backend assignment",
  "done": false
}
```

Update a task:

```powershell
$updateBody = @{ title = "Finish assignment week 2"; done = $true } | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8000/tasks/4" -Method Put -Body $updateBody -ContentType "application/json"
```

## API Documentation

With the server running, open [Scalar API documentation](http://localhost:8000/docs) for an interactive reference and request tester. The source OpenAPI 3.0 specification is available in [openapi.json](openapi.json).

### Scalar UI Screenshot

<!-- Replace the path below with the captured screenshot stored in the repository. -->
![Scalar API documentation screenshot](docs/scalar-ui.png)

## License

MIT

## Stage 6: Publish and Docs

To commit and push the documentation:

```powershell
git add README.md
git commit -m "Stage 6: publish and docs"
git push origin main
```