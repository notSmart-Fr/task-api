# Task API

Task API is a Go REST API for managing tasks and authenticating users with Supabase. It uses Go's standard `net/http` server, SQLite for task persistence, structured JSON logging, and vertical feature modules for route registration and resource lifecycle management.

## Setup

### Prerequisites

- Go 1.26.6 or later
- A Supabase project with its project URL and publishable key

Create a `.env` file in the repository root. The server loads it automatically at startup:

```dotenv
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_PUBLISHABLE_KEY=your-publishable-key
DB_PATH=tasks.db
PORT=8000
```

`SUPABASE_URL` and `SUPABASE_PUBLISHABLE_KEY` are used by Supabase authentication and token validation. `DB_PATH` and `PORT` are optional; they default to `tasks.db` and `8000`.

### Run

From the repository root, run:

```bash
go run ./cmd/server
```

The API listens on `http://localhost:8000` unless `PORT` is changed.

## Why SQLite?

SQLite was chosen because the API needs durable task storage without requiring a separate database server. It keeps the project easy to run locally, stores the entire database in a portable `tasks.db` file, and still provides SQL queries, constraints, and reliable persistence for this small CRUD service.

## Stage 4: Verify SQLite Data

The Stage 4 database check used this query in SQLite DB Browser:

```sql
SELECT id, title, done
FROM tasks
ORDER BY id ASC;
```

### DB Browser Screenshot

![SQLite DB Browser showing the tasks table](docs/SqlLite.png)

## API Reference

These five endpoints cover the main health, authentication, profile, and task-list flows. Authenticated requests must send `Authorization: Bearer <supabase-access-token>`.

| Method | Endpoint | Description | Auth required |
| --- | --- | --- | --- |
| `GET` | `/health` | Check whether the API is running | No |
| `POST` | `/auth/signup` | Create a Supabase user account | No |
| `POST` | `/auth/login` | Sign in and receive access and refresh tokens | No |
| `GET` | `/protected/profile` | Get the current authenticated user's profile | Yes |
| `GET` | `/tasks` | List stored tasks with pagination | No |

The API also exposes task CRUD operations at `GET /tasks/{id}`, `POST /tasks`, `PUT /tasks/{id}`, and `DELETE /tasks/{id}`, plus `POST /auth/logout`, `GET /`, `GET /public/info`, and the interactive documentation at `/docs`.

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

With the server running, open [Scalar API documentation](http://localhost:8000/docs) for an interactive reference and request tester. The source OpenAPI 3.0 specification is available in [swagger.json](docs/swagger.json).

### Scalar UI Screenshot

<!-- Replace the path below with the captured screenshot stored in the repository. -->
![Scalar API documentation screenshot](docs/scalar-ui.png)

## Polite Scraper (Stage 0-6)

### 1. Target Classification & Rules

- **Target site:** [Books to Scrape](https://books.toscrape.com/)
- **Scope:** The first three catalogue pages, covering approximately 60 book detail pages.
- **Robots.txt result:** [`robots.txt`](https://books.toscrape.com/robots.txt) returned HTTP 404, meaning that no robots file was found. Permission is established by the explicit sandbox statement on the site homepage.
- **Ethics commitment:** I will not reuse this code on another site without checking its rules and terms first.

---

### 2. How to Run

- **Language and dependencies:** Go 1.26.6, with `golang.org/x/net` for HTML parsing.
- **Execution command:**

```powershell
go run ./cmd/scraper
```

### 3. Politeness & Safety Rules

- **User-Agent:** Honest identifying header: `FlyRankInternship-A9/1.0 (+https://github.com/yourusername/task-api)`.
- **Timeout:** Five seconds per request.
- **Rate-limit delay:** At least 500 milliseconds between live network requests.
- **Cache-first development:** Network responses are cached in `cache/`. Subsequent runs read from the local cache and log `CACHE HIT`, avoiding unnecessary requests to the sandbox server.

### 4. Schema Shape (`output/books.json`)

Every record in `books.json` follows this schema:

- `title` (string, required)
- `product_url` (canonical absolute HTTPS URL, required)
- `price_text` (raw price string, for example `"£51.77"`, required)
- `price_gbp` (numeric price, required)
- `availability_text` (string, required)
- `rating_text` (string, required)
- `description` (nullable string; `null` if missing)
- `source_page` (provenance URL, required)
- `fetched_at` (ISO-8601 UTC timestamp, required)

Validation failures are written to `output/errors.json`.

### 5. Sample Run Report (`output/run-report.json`)

```json
{
  "start_time": "2026-09-12T17:15:00Z",
  "end_time": "2026-09-12T17:15:02Z",
  "duration_ms": 2150,
  "catalogue_pages": 3,
  "total_discovered": 61,
  "pages_fetched": 0,
  "cache_hits": 0,
  "valid_records": 60,
  "invalid_records": 0,
  "failed_pages": 1
}
```

The sample includes the deliberately broken detail URL used to verify failure handling, which accounts for the one failed page.

### 6. Architectural Notes & Ethics

- **Why no headless browser (Playwright/Puppeteer)?** The target site serves static, server-rendered HTML. All required data is present in the HTTP response, so client-side JavaScript execution is unnecessary. Avoiding a headless browser also reduces CPU and memory overhead.
- **Honest limitations:** The parser relies on specific HTML structures such as `.product_main` and `.price_color`. Major DOM changes on the target site would require the parsing selectors to be updated.
- **Ethics note:** Use an official REST or GraphQL API when one is available. Never bypass login barriers, paywalls, or rate blocks, and collect only data relevant to the pipeline.

### Final Checkpoint & Stage 6 Commit

Verify that the Git history contains at least seven meaningful commits:

```powershell
git log --oneline
```

The history should include commits for Stage 0, Stage 1, Stage 2, Stage 3, Stage 4, Stage 5, and Stage 6.

## License

MIT