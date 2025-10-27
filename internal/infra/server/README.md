# Server Integration Layer

This directory contains the integrated server implementation that connects all components of the Edna Bar Book Printing API following clean architecture principles.

## Overview

The server layer serves as the main entry point that wires together:
- **Database Service** - PostgreSQL connection management
- **Repository Layer** - Data persistence implementations
- **Application Layer** - Business logic use cases
- **Presentation Layer** - HTTP handlers and middleware

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Request  │ -> │   Middleware    │ -> │    Handlers     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        v
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Application    │ <- │   Use Cases     │ <- │   Domain        │
│   Service       │    │                 │    │   Entities      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │
        v
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Repository     │ -> │   Database      │ -> │   PostgreSQL    │
│   Registry      │    │   Service       │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Files

### `server.go`
Main server implementation containing:
- Server struct with all component dependencies
- Component initialization and wiring
- Health check endpoints
- Graceful shutdown handling

### `routes.go`
Route registration and HTTP handler setup:
- API route registration via handler registry
- Health check endpoints
- Legacy endpoint compatibility
- CORS and middleware integration

## Key Features

### 🔧 Complete Integration
- All repositories, use cases, and handlers properly wired
- Dependency injection throughout the stack
- Clean separation of concerns

### 🏥 Health Monitoring
- `/health` - Comprehensive system health
- `/health/db` - Database connection status
- `/health/app` - Application service status
- `/status` - Server runtime information

### 🛡️ Production Ready
- Graceful shutdown handling
- Request timeouts and limits
- CORS support
- Recovery middleware
- Request logging

### 📊 Monitoring & Observability
- Database connection metrics
- Request/response logging
- Error tracking
- Health status reporting

## Configuration

### Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=edna_bar
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_SCHEMA=public

# Server Configuration
PORT=8080
```

### Default Values
- **Port**: 8080 (if PORT not set)
- **Request Timeout**: 30 seconds
- **Read Timeout**: 10 seconds
- **Write Timeout**: 30 seconds
- **Idle Timeout**: 1 minute

## API Endpoints

### Core API
All business logic endpoints are mounted under `/api/`:

- **Books**: `/api/livros`
- **Authors**: `/api/autores`
- **Publishers**: `/api/editoras`
- **Printing Companies**: `/api/graficas`
- **Contracts**: `/api/contratos`
- **Printing Jobs**: `/api/printing-jobs`

### System Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Basic server information |
| `/health` | GET | Comprehensive health check |
| `/health/db` | GET | Database health status |
| `/health/app` | GET | Application service health |
| `/status` | GET | Server runtime information |
| `/api/` | GET | API documentation and info |

### Legacy Endpoints
- `/legacy/hello` - Simple hello world response
- `/legacy/health` - Basic health check

## Usage

### Starting the Server

```go
package main

import (
    "log"
    "edna/internal/infra/server"
)

func main() {
    // Create configured server
    httpServer := server.NewServer()
    
    // Start server
    log.Printf("Server starting on %s", httpServer.Addr)
    if err := httpServer.ListenAndServe(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### Running with Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Health Check Examples

```bash
# Check overall system health
curl http://localhost:8080/health

# Check database connectivity
curl http://localhost:8080/health/db

# Check application services
curl http://localhost:8080/health/app

# Get server status
curl http://localhost:8080/status
```

## Component Initialization Order

1. **Database Service** - Establishes PostgreSQL connection
2. **Repository Registry** - Creates all repository implementations
3. **Application Service** - Initializes use cases with repositories
4. **Handler Registry** - Creates HTTP handlers with use cases
5. **HTTP Server** - Configures routes and middleware

## Error Handling

### Startup Errors
- Database connection failures result in immediate shutdown
- Invalid configuration causes startup failure
- Missing dependencies are detected during initialization

### Runtime Errors
- Database disconnections are handled gracefully
- Request timeouts prevent hanging connections
- Panic recovery prevents server crashes
- Health checks detect component failures

## Middleware Stack

The following middleware is applied to all API routes:

1. **Recovery** - Panic recovery and logging
2. **Logger** - Request/response logging
3. **CORS** - Cross-origin resource sharing
4. **ContentType** - Content-type validation
5. **Timeout** - Request timeout handling

## Testing

### Integration Testing

```go
func TestServerIntegration(t *testing.T) {
    // Setup test database
    testDB := setupTestDatabase(t)
    defer testDB.Close()
    
    // Create test server
    server := createTestServer(testDB)
    
    // Test endpoints
    resp, err := http.Get(server.URL + "/health")
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

### Load Testing

```bash
# Test server performance
ab -n 1000 -c 10 http://localhost:8080/api/livros

# Test health endpoint
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/health
```

## Monitoring

### Metrics Available
- Database connection pool stats
- Request latency and throughput
- Error rates by endpoint
- Memory and CPU usage

### Health Check Responses

#### Healthy System
```json
{
  "service": "Edna Bar Book Printing API",
  "overall_status": "healthy",
  "timestamp": "2024-01-15T10:00:00Z",
  "components": {
    "database": {
      "status": "up",
      "message": "It's healthy"
    },
    "application": {
      "status": "healthy",
      "message": "All systems operational"
    },
    "server": {
      "status": "running",
      "port": 8080
    }
  }
}
```

#### Unhealthy System
```json
{
  "service": "Edna Bar Book Printing API",
  "overall_status": "unhealthy",
  "timestamp": "2024-01-15T10:00:00Z",
  "components": {
    "database": {
      "status": "down",
      "message": "Connection failed"
    }
  }
}
```

## Development

### Local Development Setup

1. **Prerequisites**
   ```bash
   # Install Go 1.21+
   # Install PostgreSQL 13+
   # Set up environment variables
   ```

2. **Database Setup**
   ```bash
   # Run migrations
   migrate -path ./migrations -database "postgres://user:pass@localhost/edna_bar?sslmode=disable" up
   ```

3. **Start Server**
   ```bash
   go run ./cmd/api/main.go
   ```

### Adding New Endpoints

1. Create domain entity (if needed)
2. Implement repository interface
3. Create use case
4. Add handler
5. Register routes in handler registry
6. Update server integration

## Deployment

### Production Checklist
- [ ] Environment variables configured
- [ ] Database migrations applied
- [ ] SSL/TLS certificates installed
- [ ] Load balancer configured
- [ ] Health checks configured
- [ ] Monitoring alerts set up
- [ ] Backup procedures in place

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: edna-bar-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: edna-bar-api
  template:
    metadata:
      labels:
        app: edna-bar-api
    spec:
      containers:
      - name: api
        image: edna-bar-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
        - name: PORT
          value: "8080"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/app
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Security Considerations

- Input validation via middleware
- SQL injection prevention via prepared statements
- CORS policy configuration
- Request rate limiting (recommended addition)
- Authentication/authorization (future enhancement)

## Performance Optimization

- Connection pooling for database
- Request timeout handling
- Graceful shutdown for zero-downtime deployments
- Health check caching (if needed)
- Static asset serving (if applicable)

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check environment variables
   - Verify database is running
   - Check network connectivity

2. **Port Already in Use**
   - Change PORT environment variable
   - Check for running instances

3. **High Memory Usage**
   - Check database connection pool size
   - Monitor for connection leaks
   - Review query performance

### Debug Mode

```bash
# Enable detailed logging
export LOG_LEVEL=debug
go run ./cmd/api/main.go
```

## Future Enhancements

- [ ] Metrics collection (Prometheus)
- [ ] Distributed tracing (Jaeger)
- [ ] Rate limiting middleware
- [ ] Authentication/authorization
- [ ] API versioning support
- [ ] WebSocket support for real-time updates
- [ ] Caching layer (Redis)
- [ ] Message queue integration