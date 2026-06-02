# Helpdesk API

A REST API for a helpdesk support system, built with Go. Tickets are created by clients and handled by support staff, with role-based access control, JWT authentication, and a layered architecture covered by unit and integration tests.

## Tech Stack

- **Go 1.25** — primary language
- **PostgreSQL 15** — database (pure SQL, no ORM)
- **chi v5** — HTTP router
- **JWT** (`golang-jwt/jwt v5`) — stateless authentication
- **bcrypt** — password hashing
- **Docker + Compose** — one-command local environment
- **Swagger** (`swaggo`) — interactive API docs
- **testify + testcontainers-go** — unit and integration tests

## Architecture

Layered design with a clear dependency direction:

```
handler  →  service  →  repository  →  PostgreSQL
(HTTP)      (business     (data
            logic, RBAC)   access)
```

The repository layer sits behind an interface defined in the service package, so the business logic is unit-tested against mocks while the repository itself is verified with real-database integration tests.

## Project Structure

```
helpdesk-api/
├── cmd/api/              # entry point, graceful shutdown
├── internal/
│   ├── config/           # env loading (port, DB, JWT, timeouts)
│   ├── handler/          # HTTP handlers, router, JSON helpers
│   ├── service/          # business logic, RBAC, sentinel errors
│   ├── repository/       # database layer (+ integration tests)
│   └── model/            # data models
├── pkg/middleware/       # JWT generation/validation, auth middleware
├── migrations/           # SQL migrations
├── docs/                 # generated Swagger spec
├── Dockerfile            # multi-stage, non-root, distroless-style alpine
├── docker-compose.yml    # api + postgres + healthcheck + auto-migrations
└── Makefile              # up / down / restart / logs / test / lint
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.25+ (only if running outside Docker)

### Run with Docker (recommended)

```bash
cp .env.example .env     # adjust values if needed
make up
```

This builds the API image, starts PostgreSQL with a healthcheck, applies migrations automatically, and launches the server. Everything comes up in a few seconds.

- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/index.html`

Other commands:

```bash
make logs       # tail container logs
make down       # stop everything
make restart    # restart the stack
make test       # run the full test suite
```

## API Endpoints

### Auth (public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Log in and receive a JWT |

### Tickets (JWT required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/tickets` | Create a ticket |
| GET | `/tickets` | List tickets (scoped by role) |
| GET | `/tickets/{id}` | Get a ticket by ID |
| PUT | `/tickets/{id}` | Update a ticket |
| DELETE | `/tickets/{id}` | Delete a ticket (admin only) |

Full request/response schemas are available in Swagger UI.

## Authentication & Authorization

Authentication is stateless via JWT. After login, the token is sent as `Authorization: Bearer <token>`; middleware validates it and injects the user ID and role into the request context.

Three roles, with access enforced in the service layer:

| Role | Tickets they can see / edit | Delete |
|--------|------------------------------|--------|
| `client` | only their own | no |
| `agent` | all | no |
| `admin` | all | yes |

A few deliberate security decisions:

- **No IDOR.** A client requesting a ticket that belongs to someone else gets the same `404 Not Found` as a ticket that does not exist — access checks happen in the service, not just in the route.
- **Anti-enumeration.** Login returns an identical error for a non-existent email and a wrong password, so the API never reveals which emails are registered.
- **Hard delete is admin-only** by design, pending a future soft-delete feature.

## Testing

```bash
make test        
```

- **Unit tests** (service layer) run against a mocked repository and cover the RBAC matrix, default-value logic, and the access-control edge cases (e.g. a stranger editing a foreign ticket returns `ErrTicketNotFound`).
- **Integration tests** (repository layer) spin up a real PostgreSQL via `testcontainers-go`, so SQL and error mapping are verified against an actual database. They are idempotent and pass under `-shuffle=on`.

## Usage Examples

Register:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John", "email": "john@example.com", "password": "12345678"}'
```

Log in (returns a token):

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "12345678"}'
```

Create a ticket:

```bash
curl -X POST http://localhost:8080/tickets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title": "Cannot log in", "description": "Login page returns 500", "priority": "high"}'
```

## Data Model

```
users     (id, name, email[unique], password, role, created_at)
tickets   (id, title, description, status, priority,
           author_id → users, assignee_id → users, created_at, updated_at)
comments  (id, ticket_id → tickets [cascade], user_id → users, text, created_at)
```

## Roadmap

- Soft delete + owner-initiated delete
- Ticket assignment to agents (`assignee_id` is already in the schema)
- Comments endpoints
- Deployment to a cloud host