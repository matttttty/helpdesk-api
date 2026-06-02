# syntax=docker/dockerfile:1

# ============================================================================
# Stage 1: Builder — compile the Go binary
# ============================================================================
FROM golang:1.25-alpine AS builder

WORKDIR /src

# Copy dependency files first — Docker caches this layer if go.mod/go.sum don't change
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a fully static binary (no glibc dependency)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /out/helpdesk-api \
    ./cmd/api

# ============================================================================
# Stage 2: Runtime — minimal image with just the binary
# ============================================================================
FROM alpine:3.20

# Add ca-certificates for HTTPS (will be needed when we add Web3 features)
RUN apk add --no-cache ca-certificates

# Create non-root user to run the app
RUN adduser -D -u 1000 helpdesk

WORKDIR /app

# Copy only the compiled binary from builder
COPY --from=builder /out/helpdesk-api .

# Switch to non-root user
USER helpdesk

EXPOSE 8080

ENTRYPOINT ["/app/helpdesk-api"]