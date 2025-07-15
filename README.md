# Service Catalog API

A simple, production-ready REST API to manage services and their versions. Built in Go using standard libraries and minimal dependencies to ensure clarity and explainability.

---

## Tech Stack

- **Language**: Go 1.21+
- **Router**: Gorilla Mux (simple, archived  used for clarity)
- **Database**: SQLite (via `sqlx`  pluggable interface)
- **Logging**: Zap (structured, production-ready)
- **Testing**: `testing` + `httptest` + `testify`

---

## How to Run
```bash
./scripts/make.sh             
```
## Or run manually

```bash
go run ./cmd/api --db-driver=sqlite3 --db-dsn=services.db --port=8080
```
The server will start on: http://localhost:8080

## API Authentication
All endpoints require a valid API key passed via the header:

```bash
X-API-Key: <your-key>
```
A default API key is auto-generated on first run and printed to the logs.

```bash
{"level":"info","ts":"2025-07-12T23:09:11.578+0530","caller":"api/main.go:31","msg":"Starting service catalog API"}
{"level":"info","ts":"2025-07-12T23:09:11.582+0530","caller":"api/main.go:122","msg":"Schema applied successfully"}
{"level":"info","ts":"2025-07-12T23:09:11.583+0530","caller":"api/main.go:78","msg":"Default API key generated: a6b7a5fb2106147bdde4193faf14a73b"}
{"level":"info","ts":"2025-07-12T23:09:11.583+0530","caller":"api/main.go:61","msg":"Listening on :8080"}
```

## Available Endpoints

| Method | Endpoint         | Description                |
| ------ | ---------------- | -------------------------- |
| GET    | `/services`      | List services (filterable) |
| POST   | `/services`      | Create a new service       |
| GET    | `/services/{id}` | Get service by ID          |
| PUT    | `/services/{id}` | Update a service           |
| DELETE | `/services/{id}` | Delete a service           |

Supports:

Filtering by name (?filter=foo)
Sorting (?sort=name)
Pagination (?page=1&limit=10)

## Project Structure
```bash
cmd/api/                  # Entry point (main.go)
internal/
  handler/                # HTTP handlers
  middleware/             # API key validation
  storage/                # Pluggable DB interface
  utils/                  # Helpers for JSON responses
  logger/                 # Zap logger setup
model/                    # Service & Version models
db/schema.sql             # SQLite schema
docs/service-catlog.yaml  # OpenAPI spec
scripts/                  # CLI and helper scripts
```

## Design Decisions
 No ORM: All DB access is raw SQL for clarity and control.

 Storage Interface: DB layer is pluggable (can support MySQL/Postgres).

 Zap Logging: Chosen for production-grade structured logs.

 API Key Auth: Simplest auth approach. Future-ready for JWT or caching.

 Makefile: Automates build/test/run. Included scripts/make.sh for ease.

## Extensibility Ideas

Add full Swagger UI via /docs

Implement version CRUD endpoints

Support multiple environments (via .env)

Add API key creation endpoint + role-based access

Write full integration test suite

Replace Gorilla Mux with chi or gin

##  Trade-offs
mux is archived but retained for readability.

SQLite used for simplicity — not ideal for production scale.

No Swagger-based code generation — routes are defined manually for clarity.

API key is stored plaintext in DB (not hashed for this demo).

## Quick API Test Scripy
To test the CRUD flow end-to-end, use the helper script:

Copy the api key from the logs 

Terminal 1:
```bash
./scripts/make.sh 
```
Terminal 2:
```bash
./scripts/test_crud.sh <apikey>
```

```bash
./scripts/test_crud.sh 96beb0bc5bd8766c7354e2da76bff479                                                                               ──(Sat,Jul12)─┘
Creating a new service...
{"code":200,"data":1,"error":"","success":true}
Created service with ID: 1

Fetching created service...
{
  "code": 200,
  "data": {
    "id": 1,
    "name": "Test Service",
    "description": "Test Description",
    "createdAt": "2025-07-12T17:50:04Z"
  },
  "error": "",
  "success": true
}

Updating service...
{
  "code": 200,
  "data": "service updated",
  "error": "",
  "success": true
}

Fetching updated service...
{
  "code": 200,
  "data": {
    "id": 1,
    "name": "Updated Service",
    "description": "Updated Description",
    "createdAt": "2025-07-12T17:50:04Z"
  },
  "error": "",
  "success": true
}

Deleting service...
{
  "code": 200,
  "data": "service deleted successfully",
  "error": "",
  "success": true
}

Verifying deletion...
{
  "code": 404,
  "data": null,
  "error": "Service not found",
  "success": false
}
```


## Note on Versions
The Versions field is part of the service model just to show that services can have versions.
For now, I’ve only handled CRUD for services since that’s what the assignment focuses on.
If needed, I’d add separate APIs to manage versions to keep things clean and modular.
