# HTTP Handlers

This package contains the HTTP handlers (presentation layer) for the Edna Bar book printing application. The handlers implement a RESTful API following clean architecture principles and provide endpoints for all business operations.

## Overview

The handlers layer is responsible for:
- HTTP request/response handling
- Request validation and parameter parsing
- Response formatting and error handling
- Route registration and middleware application
- API documentation and health checks

## Architecture

```
handlers/
├── README.md                 # This documentation
├── handler_registry.go       # Main handler coordinator
├── middleware.go            # HTTP middleware functions
├── utils.go                 # Common HTTP utilities
├── livro_handler.go         # Book operations
├── autor_handler.go         # Author operations
├── editora_handler.go       # Publisher operations
├── grafica_handler.go       # Printing company operations
├── contrato_handler.go      # Contract operations
└── imprime_handler.go       # Printing job operations
```

## API Endpoints

All API endpoints are prefixed with `/api/` and support JSON request/response format.

### Books (`/api/livros`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/livros` | List all books |
| POST | `/livros` | Create a new book |
| GET | `/livros/{isbn}` | Get book by ISBN |
| GET | `/livros/{isbn}/authors` | Get book with authors |
| PUT | `/livros/{isbn}` | Update book |
| DELETE | `/livros/{isbn}` | Delete book |
| GET | `/livros/editora/{editoraId}` | Get books by publisher |
| GET | `/livros/search/date-range?start=YYYY-MM-DD&end=YYYY-MM-DD` | Get books by publication date range |
| POST | `/livros/{isbn}/authors/{authorRG}` | Add author to book |
| DELETE | `/livros/{isbn}/authors/{authorRG}` | Remove author from book |

### Authors (`/api/autores`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/autores` | List all authors |
| POST | `/autores` | Create a new author |
| GET | `/autores/{rg}` | Get author by RG |
| GET | `/autores/{rg}/books` | Get author with books |
| PUT | `/autores/{rg}` | Update author |
| DELETE | `/autores/{rg}` | Delete author |
| GET | `/autores/search/name?name=NAME` | Search authors by name |
| POST | `/autores/{rg}/books/{isbn}` | Add author to book |
| DELETE | `/autores/{rg}/books/{isbn}` | Remove author from book |

### Publishers (`/api/editoras`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/editoras` | List all publishers |
| POST | `/editoras` | Create a new publisher |
| GET | `/editoras/{id}` | Get publisher by ID |
| GET | `/editoras/{id}/books` | Get publisher with books |
| PUT | `/editoras/{id}` | Update publisher |
| DELETE | `/editoras/{id}` | Delete publisher |
| GET | `/editoras/search/name?name=NAME` | Search publishers by name |

### Printing Companies (`/api/graficas`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/graficas` | List all printing companies |
| POST | `/graficas` | Create a new printing company |
| GET | `/graficas/{id}` | Get printing company by ID |
| GET | `/graficas/{id}/contracts` | Get company with contracts |
| GET | `/graficas/{id}/jobs` | Get company with printing jobs |
| PUT | `/graficas/{id}` | Update printing company |
| DELETE | `/graficas/{id}` | Delete printing company |
| GET | `/graficas/search/name?name=NAME` | Search companies by name |
| GET | `/graficas/search/type?type=particular\|contratada` | Get companies by type |

### Contracts (`/api/contratos`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/contratos` | List all contracts |
| POST | `/contratos` | Create a new contract |
| GET | `/contratos/{id}` | Get contract by ID |
| PUT | `/contratos/{id}` | Update contract |
| DELETE | `/contratos/{id}` | Delete contract |
| GET | `/contratos/grafica/{graficaId}` | Get contracts by printing company |
| GET | `/contratos/search/responsavel?name=NAME` | Search by responsible person |
| GET | `/contratos/search/value-range?min=VALUE&max=VALUE` | Get contracts by value range |
| GET | `/contratos/analysis` | Get contract value analysis |

### Printing Jobs (`/api/printing-jobs`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/printing-jobs` | List all printing jobs |
| POST | `/printing-jobs` | Create a new printing job |
| GET | `/printing-jobs/{isbn}/{graficaId}` | Get specific printing job |
| PUT | `/printing-jobs/{isbn}/{graficaId}` | Update printing job |
| DELETE | `/printing-jobs/{isbn}/{graficaId}` | Delete printing job |
| GET | `/printing-jobs/book/{isbn}` | Get jobs for a book |
| GET | `/printing-jobs/grafica/{graficaId}` | Get jobs for a printing company |
| GET | `/printing-jobs/search/date-range?start=YYYY-MM-DD&end=YYYY-MM-DD` | Get jobs by delivery date |
| GET | `/printing-jobs/overdue` | Get overdue jobs |
| GET | `/printing-jobs/pending` | Get pending jobs |
| GET | `/printing-jobs/statistics?start=YYYY-MM-DD&end=YYYY-MM-DD` | Get printing statistics |
| POST | `/printing-jobs/{isbn}/{graficaId}/complete` | Mark job as completed |

### System Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/` | API information and documentation |
| GET | `/api/health` | Health check endpoint |

## Request/Response Format

### Request Format

All POST and PUT requests should include a `Content-Type: application/json` header and a JSON body with the required fields.

Example book creation request:
```json
POST /api/livros
Content-Type: application/json

{
  "isbn": "978-3-16-148410-0",
  "titulo": "Clean Architecture",
  "data_de_publicacao": "2017-09-20T00:00:00Z",
  "editora_id": 1
}
```

### Response Format

All responses are in JSON format with consistent structure:

**Success Response:**
```json
{
  "isbn": "978-3-16-148410-0",
  "titulo": "Clean Architecture",
  "data_de_publicacao": "2017-09-20T00:00:00Z",
  "editora_id": 1
}
```

**Error Response:**
```json
{
  "error": "Book not found",
  "code": 404,
  "timestamp": "2023-12-07T10:30:00Z",
  "details": "Book with ISBN 978-3-16-148410-0 not found"
}
```

## HTTP Status Codes

| Code | Description | Usage |
|------|-------------|-------|
| 200 | OK | Successful GET, PUT operations |
| 201 | Created | Successful POST operations |
| 400 | Bad Request | Invalid request data |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server errors |

## Middleware

The API includes several middleware layers:

### Recovery Middleware
- Catches panics and returns a proper error response
- Logs stack traces for debugging

### Logger Middleware
- Logs all HTTP requests with method, path, status code, and duration
- Helps with monitoring and debugging

### CORS Middleware
- Handles Cross-Origin Resource Sharing
- Allows requests from any origin (configurable)

### Content-Type Middleware
- Sets appropriate content type headers
- Ensures JSON responses for API endpoints

### Timeout Middleware
- Adds request timeout to prevent hanging requests
- Default timeout: 30 seconds

## Error Handling

The handlers provide comprehensive error handling:

### Validation Errors
```json
{
  "error": "Invalid request body",
  "code": 400,
  "timestamp": "2023-12-07T10:30:00Z",
  "details": "Field 'titulo' is required"
}
```

### Business Rule Errors
```json
{
  "error": "Failed to delete author",
  "code": 400,
  "timestamp": "2023-12-07T10:30:00Z",
  "details": "cannot delete author with associated books"
}
```

### Not Found Errors
```json
{
  "error": "Book not found",
  "code": 404,
  "timestamp": "2023-12-07T10:30:00Z",
  "details": "Book with ISBN 978-3-16-148410-0 not found"
}
```

## Usage Example

### Initialize Handler Registry

```go
// Create application service
appService := applications.NewApplicationService(config)

// Create handler registry
handlerConfig := handlers.HandlerConfig{
    ApplicationService: appService,
    RequestTimeout:     30 * time.Second,
}
handlerRegistry := handlers.NewHandlerRegistry(handlerConfig)

// Create HTTP server
mux := http.NewServeMux()
handlerRegistry.RegisterRoutes(mux)

server := &http.Server{
    Addr:    ":8080",
    Handler: mux,
}

log.Fatal(server.ListenAndServe())
```

### Making API Requests

```bash
# Create a new book
curl -X POST http://localhost:8080/api/livros \
  -H "Content-Type: application/json" \
  -d '{
    "isbn": "978-3-16-148410-0",
    "titulo": "Clean Architecture",
    "data_de_publicacao": "2017-09-20T00:00:00Z",
    "editora_id": 1
  }'

# Get book with authors
curl http://localhost:8080/api/livros/978-3-16-148410-0/authors

# Search books by date range
curl "http://localhost:8080/api/livros/search/date-range?start=2020-01-01&end=2023-12-31"

# Get printing job statistics
curl "http://localhost:8080/api/printing-jobs/statistics?start=2023-01-01&end=2023-12-31"

# Check API health
curl http://localhost:8080/api/health
```

## Testing

### Unit Testing

Each handler should be tested with:
- Valid request scenarios
- Invalid request validation
- Error handling scenarios
- Middleware functionality

```go
func TestLivroHandler_Create(t *testing.T) {
    // Setup
    mockUsecase := &MockLivroUsecase{}
    handler := NewLivroHandler(mockUsecase)
    
    // Test cases
    t.Run("successful creation", func(t *testing.T) {
        // Test implementation
    })
    
    t.Run("validation error", func(t *testing.T) {
        // Test implementation
    })
}
```

### Integration Testing

Test the complete HTTP flow:
- Request parsing
- Use case execution
- Response formatting
- Error handling

## Security Considerations

### Input Validation
- All inputs are validated at the handler level
- JSON decoder disallows unknown fields
- Parameter validation with appropriate error messages

### CORS Configuration
- Currently allows all origins (development setting)
- Should be restricted in production

### Request Timeouts
- All requests have timeouts to prevent resource exhaustion
- Configurable timeout values

### Error Information
- Error messages don't expose sensitive information
- Stack traces only in development mode

## Performance Considerations

### Request Timeouts
- 30-second default timeout for all requests
- Prevents hanging requests and resource leaks

### JSON Encoding
- Streaming JSON encoder for better memory usage
- Efficient request/response processing

### Middleware Chain
- Lightweight middleware for minimal overhead
- Proper error recovery without affecting other requests

## Monitoring and Observability

### Logging
- All requests are logged with timing information
- Error conditions are logged with details
- Structured logging format for easy parsing

### Health Checks
- `/api/health` endpoint for service monitoring
- Returns service status and version information

### Metrics
- Request duration logging
- Error rate tracking through logs
- Status code distribution

## Extension Points

### Adding New Endpoints
1. Create handler method following naming convention
2. Add route registration in `RegisterRoutes`
3. Include proper validation and error handling
4. Update documentation

### Custom Middleware
1. Implement `MiddlewareFunc` interface
2. Add to middleware chain in `handler_registry.go`
3. Ensure proper error handling and recovery

### Response Formats
1. Use common utilities in `utils.go`
2. Maintain consistent error response format
3. Include appropriate HTTP status codes