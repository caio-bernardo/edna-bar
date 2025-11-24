# Build stage
FROM golang:1.24.3-alpine AS builder

# Install build dependencies
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Create build directory
WORKDIR /build

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the binary with module-aware path and without problematic static linker flags.
# Use module-aware relative path to the package and keep CGO disabled for portability.
# Avoid -extldflags "-static" and -installsuffix cgo flags which can cause issues on Alpine.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o main ./cmd/api

# Production stage
FROM alpine:3.19 AS runner

# Install runtime dependencies and security updates
RUN apk update && apk add --no-cache \
    ca-certificates \
    curl \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create a non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/main /app/main

# Copy static assets if they exist
COPY --from=builder /build/assets /app/assets

# Change ownership of the app directory to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (configurable via PORT env var, defaults to 8080)
EXPOSE 8080

# Add metadata labels
LABEL maintainer="Edna Bar Team" \
      version="1.0" \
      description="Edna Bar Management System API" \
      org.opencontainers.image.title="edna-bar-api" \
      org.opencontainers.image.description="Go-based API for Edna Bar Management System" \
      org.opencontainers.image.version="1.0"

# Health check using the application's health endpoint
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:${PORT:-8080}/api/v1/health || exit 1

# Set environment variables for production
ENV TZ=UTC \
    PORT=8080

# Run the application
CMD ["/app/main"]
