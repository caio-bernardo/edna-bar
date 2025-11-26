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
FROM gcr.io/distroless/static:nonroot AS runner

WORKDIR /app
COPY --from=builder /build/main /app/main

# Service port
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/main"]
