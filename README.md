# Service Catalog API

A simple, production-ready REST API to manage services and their versions. Built in Go using standard libraries and minimal dependencies to ensure clarity and explainability.

---

## Tech Stack

* **Language**: Go 1.21+
* **Router**: Gorilla Mux (simple, archived - used for clarity)
* **Database**: SQLite (via `sqlx`, pluggable interface)
* **Logging**: Zap (structured, production-ready)
* **Testing**: `testing` + `httptest` + `testify`

---

## How to Run

```bash
./scripts/make.sh
```

### Or run manually

```bash
go run ./cmd/api --db-driver=sqlite3 --db-dsn=services.db --port=8080
```

The server will start on: [http://localhost:8080](http://localhost:8080)

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

| Method | Endpoint                  | Description                        |
| ------ | ------------------------- | ---------------------------------- |
| GET    | `/services`               | List services (filterable)         |
| POST   | `/services`               | Create a new service + version     |
| GET    | `/services/{id}`          | Get service by ID (with versions)  |
| PUT    | `/services/{id}`          | Update a service                   |
| DELETE | `/services/{id}`          | Delete a service                   |
| GET    | `/services/{id}/versions` | List versions for a service        |
| POST   | `/services/{id}/versions` | Create a new version for a service |
| GET    | `/versions/{id}`          | Get version by ID                  |
| DELETE | `/versions/{id}`          | Delete version by ID               |

Supports:

* Filtering by name/description (`?filter=dummy`)
* Sorting (`?sort=name` or `?sort=createdAt`)
* Pagination (`?page=1&limit=100`)

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

* **No ORM**: All DB access is raw SQL for clarity and control.
* **Storage Interface**: DB layer is pluggable (can support MySQL/Postgres).
* **Zap Logging**: Chosen for production-grade structured logs.
* **API Key Auth**: Simplest auth approach. Future-ready for JWT or caching.
* **Makefile**: Automates build/test/run. Included scripts/make.sh for ease.

## Extensibility Ideas

* Add full Swagger UI via /docs
* Implement full version CRUD (PUT coming soon)
* Support multiple environments (via .env)
* Add API key creation endpoint + role-based access
* Write a full integration test suite
* Replace Gorilla Mux with chi or gin

## Trade-offs

* mux is archived but used for readability.
* SQLite used for simplicity cannot be ideal for production scale.
* No Swagger-based code generation routes are defined manually for clarity.
* API key is stored plaintext in DB (not hashed for this test).

## Quick API Test Script

To test the CRUD flow end-to-end, use the helper script:

Copy the API key from the logs.

### Terminal 1:

```bash
./scripts/make.sh
```

### Terminal 2:

```bash
./scripts/test_crud.sh <apikey>
```
<details>
<summary>Click to expand test output</summary>

```bash
./scripts/test_crud.sh 8cf9a281dccea6d6c3129fe3d9a330db                              
==> [1] Creating 1000 services...
==> [2] Listing services without any query parameters (should return default 20)...
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "name": "Service 1",
      "description": "Description for service 1",
      "createdAt": "2025-07-15T18:06:50Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 2,
      "name": "Service 2",
      "description": "Description for service 2",
      "createdAt": "2025-07-15T18:06:50Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 3,
      "name": "Service 3",
      "description": "Description for service 3",
      "createdAt": "2025-07-15T18:06:50Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 4,
      "name": "Service 4",
      "description": "Description for service 4",
      "createdAt": "2025-07-15T18:06:50Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 5,
      "name": "Service 5",
      "description": "Description for service 5",
      "createdAt": "2025-07-15T18:06:50Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 6,
      "name": "Service 6",
      "description": "Description for service 6",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 7,
      "name": "Service 7",
      "description": "Description for service 7",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 8,
      "name": "Service 8",
      "description": "Description for service 8",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 9,
      "name": "Service 9",
      "description": "Description for service 9",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 10,
      "name": "Service 10",
      "description": "Description for service 10",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 11,
      "name": "Service 11",
      "description": "Description for service 11",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 12,
      "name": "Service 12",
      "description": "Description for service 12",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 13,
      "name": "Service 13",
      "description": "Description for service 13",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 14,
      "name": "Service 14",
      "description": "Description for service 14",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 15,
      "name": "Service 15",
      "description": "Description for service 15",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 16,
      "name": "Service 16",
      "description": "Description for service 16",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 17,
      "name": "Service 17",
      "description": "Description for service 17",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 18,
      "name": "Service 18",
      "description": "Description for service 18",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 19,
      "name": "Service 19",
      "description": "Description for service 19",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 20,
      "name": "Service 20",
      "description": "Description for service 20",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    }
  ],
  "error": "",
  "success": true
}
==> [3] Filtering services by name with sort=name...
{
  "code": 200,
  "data": [
    {
      "id": 2,
      "name": "Service 2",
      "description": "Description for service 2",
      "createdAt": "2025-07-15T18:06:50Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 20,
      "name": "Service 20",
      "description": "Description for service 20",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 200,
      "name": "Service 200",
      "description": "Description for service 200",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 201,
      "name": "Service 201",
      "description": "Description for service 201",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 202,
      "name": "Service 202",
      "description": "Description for service 202",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 203,
      "name": "Service 203",
      "description": "Description for service 203",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 204,
      "name": "Service 204",
      "description": "Description for service 204",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 205,
      "name": "Service 205",
      "description": "Description for service 205",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 206,
      "name": "Service 206",
      "description": "Description for service 206",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 207,
      "name": "Service 207",
      "description": "Description for service 207",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 208,
      "name": "Service 208",
      "description": "Description for service 208",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 209,
      "name": "Service 209",
      "description": "Description for service 209",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 21,
      "name": "Service 21",
      "description": "Description for service 21",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 210,
      "name": "Service 210",
      "description": "Description for service 210",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 211,
      "name": "Service 211",
      "description": "Description for service 211",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 212,
      "name": "Service 212",
      "description": "Description for service 212",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 213,
      "name": "Service 213",
      "description": "Description for service 213",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 214,
      "name": "Service 214",
      "description": "Description for service 214",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 215,
      "name": "Service 215",
      "description": "Description for service 215",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 216,
      "name": "Service 216",
      "description": "Description for service 216",
      "createdAt": "2025-07-15T18:06:52Z",
      "versions": [
        "v1.0.0"
      ]
    }
  ],
  "error": "",
  "success": true
}
==> [3] Filtering services by description with sort=createdAt...
{
  "code": 200,
  "data": [
    {
      "id": 3,
      "name": "Service 3",
      "description": "Description for service 3",
      "createdAt": "2025-07-15T18:06:50Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 30,
      "name": "Service 30",
      "description": "Description for service 30",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 31,
      "name": "Service 31",
      "description": "Description for service 31",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 32,
      "name": "Service 32",
      "description": "Description for service 32",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 33,
      "name": "Service 33",
      "description": "Description for service 33",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 34,
      "name": "Service 34",
      "description": "Description for service 34",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 35,
      "name": "Service 35",
      "description": "Description for service 35",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 36,
      "name": "Service 36",
      "description": "Description for service 36",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 37,
      "name": "Service 37",
      "description": "Description for service 37",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 38,
      "name": "Service 38",
      "description": "Description for service 38",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 39,
      "name": "Service 39",
      "description": "Description for service 39",
      "createdAt": "2025-07-15T18:06:51Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 300,
      "name": "Service 300",
      "description": "Description for service 300",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 301,
      "name": "Service 301",
      "description": "Description for service 301",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 302,
      "name": "Service 302",
      "description": "Description for service 302",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 303,
      "name": "Service 303",
      "description": "Description for service 303",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 304,
      "name": "Service 304",
      "description": "Description for service 304",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 305,
      "name": "Service 305",
      "description": "Description for service 305",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 306,
      "name": "Service 306",
      "description": "Description for service 306",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 307,
      "name": "Service 307",
      "description": "Description for service 307",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    },
    {
      "id": 308,
      "name": "Service 308",
      "description": "Description for service 308",
      "createdAt": "2025-07-15T18:06:53Z",
      "versions": [
        "v1.0.0"
      ]
    }
  ],
  "error": "",
  "success": true
}
==> [4] Get Service by ID (serviceId=10)...
{
  "code": 200,
  "data": {
    "id": 10,
    "name": "Service 10",
    "description": "Description for service 10",
    "createdAt": "2025-07-15T18:06:51Z",
    "versions": [
      "v1.0.0"
    ]
  },
  "error": "",
  "success": true
}
==> [5] Creating a patch version (v1.0.1) for serviceId=10...
==> New version created with ID: 1001
==> [6] Get Service by ID (serviceId=10) after patch version...
{
  "code": 200,
  "data": {
    "id": 10,
    "name": "Service 10",
    "description": "Description for service 10",
    "createdAt": "2025-07-15T18:06:51Z",
    "versions": [
      "v1.0.0",
      "v1.0.1"
    ]
  },
  "error": "",
  "success": true
}
==> [7] Deleting version ID 1001...
==> [8] Get Service by ID (serviceId=10) after deletion...
{
  "code": 200,
  "data": {
    "id": 10,
    "name": "Service 10",
    "description": "Description for service 10",
    "createdAt": "2025-07-15T18:06:51Z",
    "versions": [
      "v1.0.0"
    ]
  },
  "error": "",
  "success": true
} 
```
</details> 



## Note on Versions

The service version history is fully supported with dedicated endpoints. A service can have multiple versions. Version creation is triggered on service creation and can be managed independently via the `/versions` endpoints.

Each version contains:

* `version`: version tag (e.g., v1.0.1)
* `changelog`: optional notes
* `createdAt`: timestamp

This modular design keeps service metadata and version history decoupled and extensible.
