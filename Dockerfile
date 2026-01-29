# Build stage - using Go 1.24 for compatibility with go.mod (pgx requires updated toolchain)
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Install build dependencies for PostgreSQL
RUN apt-get update && apt-get install -y gcc libc6-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application (CGO disabled for PostgreSQL-only)
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags='-s -w' -o main ./cmd/api

# Runtime stage - minimal image
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies (PostgreSQL client libs)
RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  tzdata \
  && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Set environment variables
ENV APP_ENV=production
ENV APP_PORT=8080

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
