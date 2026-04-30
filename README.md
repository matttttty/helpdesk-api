# Helpdesk API

A REST API for a helpdesk support system built with Go.
> ⚠️ This project is currently under active development. Some features may be incomplete or subject to change.
## Tech Stack

- **Go** - primary language
- **PostgreSQL** - database
- **Chi** - HTTP router
- **JWT** - authentication
- **Docker** - database setup

## Project Structure

```
helpdesk-api/
├── cmd/api/          # entry point
├── internal/
│   ├── handler/      # HTTP handlers
│   ├── service/      # business logic
│   ├── repository/   # database layer
│   └── model/        # data models
├── pkg/middleware/   # JWT auth middleware
├── migrations/       # SQL migrations
└── .env              # environment variables
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker

### Setup

1. Clone the repository:
```bash
git clone https://github.com/matttttty/helpdesk-api.git
cd helpdesk-api
```

2. Create `.env` file:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=helpdesk
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h
```

3. Start PostgreSQL:
```bash
docker run --name helpdesk-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=helpdesk \
  -p 5432:5432 \
  -d postgres:15
```

4. Run migrations:
```bash
docker exec -i helpdesk-postgres psql -U postgres -d helpdesk < migrations/001_init.sql
```

5. Start the server:
```bash
go run cmd/api/main.go
```

Server runs on `http://localhost:8080`

## API Endpoints

### Auth (public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Login and get JWT token |

### Tickets (requires JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/tickets` | Create a ticket |
| GET | `/tickets` | Get all tickets |
| GET | `/tickets/{id}` | Get ticket by ID |
| PUT | `/tickets/{id}` | Update ticket |
| DELETE | `/tickets/{id}` | Delete ticket |

## Usage Examples

### Register
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John", "email": "john@example.com", "password": "12345678"}'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "12345678"}'
```

### Create Ticket
```bash
curl -X POST http://localhost:8080/tickets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"title": "Issue with login", "description": "Cannot login to the system", "priority": "high"}'
```

## User Roles

- `client` - can create and view own tickets
- `agent` - can view and update all tickets
- `admin` - full access including delete