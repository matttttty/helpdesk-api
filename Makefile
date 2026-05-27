.PHONY: help up down restart logs build run test lint fmt clean migrate-up migrate-down

# Default target
help:
	@echo "Available commands:"
	@echo "  make up         - Start all services (postgres + api)"
	@echo "  make down       - Stop all services"
	@echo "  make restart    - Restart all services"
	@echo "  make logs       - Tail api logs"
	@echo "  make build      - Build the Go binary locally"
	@echo "  make run        - Run the app locally (postgres must be up)"
	@echo "  make test       - Run all tests"
	@echo "  make lint       - Run golangci-lint"
	@echo "  make fmt        - Format code with gofmt"
	@echo "  make clean      - Remove containers, volumes, build artifacts"

# Docker Compose
up:
	docker compose up -d --build
	@echo ""
	@echo "API:      http://localhost:8080"
	@echo "Postgres: localhost:5432"
	@echo ""
	@echo "Tail logs: make logs"

down:
	docker compose down

restart: down up

logs:
	docker compose logs -f api

# Local development (no Docker for API, only Postgres in container)
build:
	go build -o bin/helpdesk-api ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v -race ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

clean:
	docker compose down -v
	rm -rf bin/