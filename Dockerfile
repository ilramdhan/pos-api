# Build stage - using Debian for better CGO compatibility
FROM golang:1.21-bookworm AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y gcc libc6-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main ./cmd/api

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  tzdata \
  sqlite3 \
  && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Create data directory
RUN mkdir -p /app/data

# Set environment variables
ENV APP_ENV=production
ENV APP_PORT=8080
ENV DB_PATH=/app/data/pos.db

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
