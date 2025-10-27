# Edna Bar - Book Printing Management API

A comprehensive RESTful API for managing book printing operations, built with Go following clean architecture principles.

## 🏗️ Architecture Overview

This application implements clean architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   HTTP      │  │ Middleware  │  │   Handlers  │        │
│  │   Server    │  │   Stack     │  │             │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────────┐
│                   Application Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Use Cases  │  │    DTOs     │  │  Services   │        │
│  │             │  │             │  │             │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────────┐
│                     Domain Layer                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Entities   │  │ Value       │  │ Repository  │        │
│  │             │  │ Objects     │  │ Interfaces  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                               │
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ PostgreSQL  │  │ Repository  │  │  Database   │        │
│  │ Database    │  │ Implements  │  │  Service    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Features

### Business Domain
- **Book Management**: Complete CRUD operations for books with ISBN validation
- **Author Management**: Author profiles and book associations
- **Publisher Management**: Publishing house information and catalogs
- **Printing Companies**: Support for both private and contracted printing services
- **Contract Management**: Printing contracts with financial tracking
- **Printing Jobs**: Job scheduling, copy tracking, and delivery management

### Technical Features
- **Clean Architecture**: Domain-driven design with dependency inversion
- **RESTful API**: Well-structured HTTP endpoints with proper status codes
- **Database Integration**: PostgreSQL with migrations and connection pooling
- **Middleware Stack**: CORS, logging, recovery, timeouts, and content validation
- **Health Monitoring**: Comprehensive health checks for all system components
- **Graceful Shutdown**: Proper resource cleanup and connection handling

## 📋 Requirements

- **Go**: 1.21 or higher
- **PostgreSQL**: 13 or higher
- **Migration Tool**: golang-migrate (optional)

## 🛠️ Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd edna-bar
```

### 2. Environment Configuration
Create a `.env` file in the root directory:

```env
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

### 3. Database Setup
```bash
# Create database
createdb edna_bar

# Run migrations (if using golang-migrate)
migrate -path ./migrations -database "postgres://username:password@localhost/edna_bar?sslmode=disable" up

# Or run the SQL files manually
psql -d edna_bar -f migrations/000001_add_editora.up.sql
psql -d edna_bar -f migrations/000002_add_livro.up.sql
# ... continue with all migration files
```

### 4. Install Dependencies
```bash
go mod download
```

### 5. Run the Application
```bash
go run ./cmd/api/main.go
```

The server will start on `http://localhost:8080`

### 6. Explore the API Documentation
Visit the interactive Swagger UI at `http://localhost:8080/docs/` to:
- Browse all available endpoints
- Test API calls directly from the browser
- View request/response schemas and examples
- Download the OpenAPI specification

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/
```

### 📚 Interactive API Documentation

The API includes comprehensive Swagger UI documentation for easy exploration and testing:

**Swagger UI**: [http://localhost:8080/docs/](http://localhost:8080/docs/)

Features:
- **Interactive Testing**: Try API calls directly from the browser
- **Complete Schema Documentation**: View all request/response models
- **Example Payloads**: See sample data for all endpoints
- **Authentication Testing**: Test secured endpoints (when implemented)
- **Export Options**: Download OpenAPI specification

**Quick Start with Swagger UI**:
1. Start the server: `go run ./cmd/api/main.go`
2. Open [http://localhost:8080/docs/](http://localhost:8080/docs/)
3. Select any endpoint to view documentation
4. Click "Try it out" to test API calls
5. Fill in parameters and click "Execute"

### Core Resources

#### Books (`/api/livros`)
- `GET /api/livros` - List all books
- `POST /api/livros` - Create a new book
- `GET /api/livros/{isbn}` - Get book by ISBN
- `PUT /api/livros/{isbn}` - Update book
- `DELETE /api/livros/{isbn}` - Delete book
- `GET /api/livros/{isbn}/authors` - Get book authors
- `POST /api/livros/{isbn}/authors` - Add author to book

#### Authors (`/api/autores`)
- `GET /api/autores` - List all authors
- `POST /api/autores` - Create a new author
- `GET /api/autores/{rg}` - Get author by RG
- `PUT /api/autores/{rg}` - Update author
- `DELETE /api/autores/{rg}` - Delete author
- `GET /api/autores/{rg}/books` - Get author's books

#### Publishers (`/api/editoras`)
- `GET /api/editoras` - List all publishers
- `POST /api/editoras` - Create a new publisher
- `GET /api/editoras/{id}` - Get publisher by ID
- `PUT /api/editoras/{id}` - Update publisher
- `DELETE /api/editoras/{id}` - Delete publisher

#### Printing Companies (`/api/graficas`)
- `GET /api/graficas` - List all printing companies
- `POST /api/graficas` - Create a new printing company
- `GET /api/graficas/{id}` - Get printing company by ID
- `PUT /api/graficas/{id}` - Update printing company
- `DELETE /api/graficas/{id}` - Delete printing company

#### Contracts (`/api/contratos`)
- `GET /api/contratos` - List all contracts
- `POST /api/contratos` - Create a new contract
- `GET /api/contratos/{id}` - Get contract by ID
- `PUT /api/contratos/{id}` - Update contract
- `DELETE /api/contratos/{id}` - Delete contract

#### Printing Jobs (`/api/printing-jobs`)
- `GET /api/printing-jobs` - List all printing jobs
- `POST /api/printing-jobs` - Create a new printing job
- `GET /api/printing-jobs/{id}` - Get printing job details
- `PUT /api/printing-jobs/{id}` - Update printing job
- `DELETE /api/printing-jobs/{id}` - Delete printing job

### System Endpoints

#### Health Checks
- `GET /health` - Comprehensive system health
- `GET /health/db` - Database connection status
- `GET /health/app` - Application service status
- `GET /status` - Server runtime information

#### API Information & Documentation
- `GET /api/` - API documentation and available endpoints
- `GET /docs/` - Interactive Swagger UI
- `GET /docs/swagger.yaml` - OpenAPI specification

#### Interactive Documentation
- `GET /docs/` - Swagger UI for interactive API testing
- `GET /swagger/` - Alternative Swagger UI interface
- `GET /docs/swagger.yaml` - OpenAPI 3.0 specification

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests (requires database)
go test ./internal/infra/server/...

# Skip database-dependent tests
TEST_DB_SKIP=true go test ./...

# Generate and view API documentation
make swagger
```

### Example Test Commands
```bash
# Test specific package
go test ./internal/domain/...

# Run benchmarks
go test -bench=. ./internal/infra/server/...

# Verbose test output
go test -v ./internal/applications/...

# Generate Swagger documentation
make swagger

# View Swagger UI (with server running)
make swagger-ui
```

## 🔧 Development

### Project Structure
```
edna-bar/
├── cmd/api/                 # Application entry point
├── internal/
│   ├── domain/              # Business entities, value objects, interfaces
│   ├── applications/        # Use cases and application services
│   ├── adapters/
│   │   ├── handlers/        # HTTP handlers and middleware
│   │   └── repository/      # Database repository implementations
│   └── infra/
│       ├── database/        # Database connection and service
│       └── server/          # HTTP server and routing
├── migrations/              # Database migration files
├── examples/                # Usage examples and demos
└── README.md
```

### Adding New Features

1. **Domain Layer**: Define entities and business rules
2. **Repository Layer**: Create database operations
3. **Application Layer**: Implement use cases
4. **Handler Layer**: Create HTTP endpoints
5. **Integration**: Wire components in server setup

### Code Style
- Follow Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Write comprehensive tests for business logic
- Document public APIs and complex functions

## 🐳 Docker Deployment

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./main"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_DATABASE=edna_bar
      - DB_USERNAME=postgres
      - DB_PASSWORD=postgres
      - PORT=8080
    depends_on:
      - postgres

  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=edna_bar
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## 📊 Database Schema

### Core Tables
- **Editora**: Publishers (id, nome, endereco)
- **Autor**: Authors (RG, nome, endereco)
- **Livro**: Books (ISBN, titulo, data_de_publicacao, editora_id)
- **Grafica**: Printing companies (id, nome)

### Specialized Tables
- **Particular**: Private printing companies (grafica_id)
- **Contratada**: Contracted printing companies (grafica_id, endereco)
- **Contrato**: Printing contracts (id, valor, nome_responsavel, grafica_cont_id)

### Relationship Tables
- **Escreve**: Author-Book relationships (ISBN, RG)
- **Imprime**: Printing jobs (lisbn, grafica_id, nto_copias, data_entrega)

## 🔍 Monitoring & Observability

### Health Check Endpoints
Monitor system health with these endpoints:

```bash
# Overall system health
curl http://localhost:8080/health

# Database connectivity
curl http://localhost:8080/health/db

# Application services
curl http://localhost:8080/health/app
```

### API Documentation & Testing
Explore and test the API using the interactive documentation:

```bash
# Open Swagger UI in browser
open http://localhost:8080/docs/

# Get OpenAPI specification
curl http://localhost:8080/docs/swagger.yaml

# Alternative Swagger UI
open http://localhost:8080/swagger/
```

### Logging
The application provides structured logging for:
- HTTP requests and responses
- Database operations
- Error conditions
- System startup and shutdown

## 🚀 Production Deployment

### Environment Setup
- Configure production database credentials
- Set up SSL/TLS certificates
- Configure load balancer
- Set up monitoring and alerting

### Performance Considerations
- Database connection pooling
- Request timeouts and rate limiting
- Resource monitoring
- Backup procedures

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

### Development Guidelines
- Follow clean architecture principles
- Write comprehensive tests
- Document API changes
- Use meaningful commit messages

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- **Interactive API Docs**: [http://localhost:8080/docs/](http://localhost:8080/docs/)
- **API Information**: [http://localhost:8080/api/](http://localhost:8080/api/)
- **Usage Examples**: Review examples in the `/examples` directory  
- **Health Monitoring**: Check health endpoints for system status
- **Troubleshooting**: Review logs for detailed error information
- **OpenAPI Spec**: Download from [http://localhost:8080/docs/swagger.yaml](http://localhost:8080/docs/swagger.yaml)

## 🎯 Roadmap

### Completed ✅
- [x] **Complete REST API** with full CRUD operations
- [x] **Clean Architecture** implementation
- [x] **Interactive API Documentation** with Swagger UI
- [x] **Health Monitoring** endpoints
- [x] **Database Integration** with PostgreSQL
- [x] **Comprehensive Testing** suite

### Planned Features 🚀
- [ ] **Authentication and Authorization** (JWT-based)
- [ ] **API Rate Limiting** and throttling
- [ ] **Metrics Collection** (Prometheus integration)
- [ ] **Distributed Tracing** (Jaeger support)  
- [ ] **WebSocket Support** for real-time updates
- [ ] **Caching Layer** (Redis integration)
- [ ] **Message Queue Integration** for async processing
- [ ] **API Versioning** support
- [ ] **Advanced Search** and filtering capabilities
- [ ] **Bulk Operations** for efficient data processing

---

**Edna Bar Book Printing API** - A complete solution for managing book printing operations with clean architecture and modern Go practices.