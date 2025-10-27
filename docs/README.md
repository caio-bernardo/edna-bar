# API Documentation

This directory contains the OpenAPI 3.0 specification and related documentation for the Edna Bar Book Printing API.

## 📚 Swagger UI

The API documentation is available through Swagger UI at multiple endpoints:

- **Primary Swagger UI**: [http://localhost:8080/docs/](http://localhost:8080/docs/)
- **Alternative UI**: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- **OpenAPI Specification**: [http://localhost:8080/docs/swagger.yaml](http://localhost:8080/docs/swagger.yaml)

## 📁 Files

- **`swagger.yaml`** - Complete OpenAPI 3.0 specification
- **`README.md`** - This documentation file

## 🚀 Quick Start

1. **Start the server**:
   ```bash
   go run ./cmd/api/main.go
   ```

2. **Open Swagger UI**:
   Navigate to [http://localhost:8080/docs/](http://localhost:8080/docs/)

3. **Explore the API**:
   - Browse available endpoints by tags
   - Try out API calls directly from the UI
   - View request/response schemas
   - See example payloads

## 📖 API Overview

The Edna Bar Book Printing API provides comprehensive functionality for managing book printing operations:

### 📚 Core Resources

- **Books (`/api/livros`)**
  - Create, read, update, delete books
  - Manage book-author relationships
  - Search by title, ISBN, publication date

- **Authors (`/api/autores`)**
  - Manage author profiles
  - View author's published books
  - Search by name

- **Publishers (`/api/editoras`)**
  - Publishing house management
  - Publisher catalog operations

- **Printing Companies (`/api/graficas`)**
  - Support for private and contracted companies
  - Company type classification
  - Address management for contracted companies

- **Contracts (`/api/contratos`)**
  - Printing contract management
  - Financial tracking
  - Responsible person assignment

- **Printing Jobs (`/api/printing-jobs`)**
  - Job scheduling and tracking
  - Copy count management
  - Delivery date monitoring
  - Overdue job detection

### 🔧 System Endpoints

- **Health Checks**: `/health`, `/health/db`, `/health/app`
- **Server Status**: `/status`
- **API Information**: `/api/`

## 🛠️ Using the API

### Authentication

Currently, the API does not require authentication. Future versions may include JWT-based authentication.

### Content Type

All requests and responses use `application/json` content type.

### Error Handling

The API returns structured error responses:

```json
{
  "error": "Resource not found",
  "code": "NOT_FOUND",
  "details": "Book with ISBN 978-0123456789 was not found",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Common HTTP Status Codes

- **200 OK** - Successful GET request
- **201 Created** - Successful POST request
- **204 No Content** - Successful DELETE request
- **400 Bad Request** - Invalid input data
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error

## 📋 Example Requests

### Create a Book

```bash
curl -X POST http://localhost:8080/api/livros \
  -H "Content-Type: application/json" \
  -d '{
    "isbn": "978-0307474728",
    "titulo": "Cem Anos de Solidão",
    "data_de_publicacao": "1967-06-05",
    "editora_id": 1
  }'
```

### Get All Authors

```bash
curl http://localhost:8080/api/autores
```

### Create a Printing Job

```bash
curl -X POST http://localhost:8080/api/printing-jobs \
  -H "Content-Type: application/json" \
  -d '{
    "isbn": "978-0307474728",
    "grafica_id": 1,
    "nto_copias": 5000,
    "data_entrega": "2024-02-15"
  }'
```

## 🔄 Generating Documentation

### Automatic Generation

The server includes embedded documentation that's always available. However, you can also generate updated documentation from code annotations:

```bash
# Generate swagger docs from code annotations
make swagger

# Or manually with swag
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

### Manual Updates

You can manually edit the `swagger.yaml` file to:
- Add new endpoints
- Update descriptions
- Modify examples
- Add security schemes

## 🧪 Testing with Swagger UI

Swagger UI provides an interactive interface for testing the API:

1. **Select an endpoint** from the list
2. **Click "Try it out"** to enable the form
3. **Fill in parameters** and request body
4. **Click "Execute"** to make the request
5. **View the response** including status code, headers, and body

### Testing Examples

#### Create a Publisher
1. Navigate to `POST /api/editoras`
2. Click "Try it out"
3. Enter request body:
   ```json
   {
     "nome": "Test Publisher",
     "endereco": "Test Address"
   }
   ```
4. Click "Execute"
5. Check the response for the created publisher with assigned ID

#### Search Books
1. Navigate to `GET /api/livros`
2. Click "Try it out"
3. Add query parameters:
   - `title`: "Solidão"
   - `limit`: 10
4. Click "Execute"
5. View matching books in the response

## 📊 Data Models

### Core Entities

- **Book**: ISBN, title, publication date, publisher ID
- **Author**: RG, name, address
- **Publisher**: ID, name, address
- **Printing Company**: ID, name, type (particular/contracted), address
- **Contract**: ID, value, responsible person, company ID
- **Printing Job**: ISBN, company ID, copies, delivery date

### Relationships

- Books ↔ Authors (many-to-many via `Escreve`)
- Books → Publishers (many-to-one)
- Printing Jobs → Books (many-to-one)
- Printing Jobs → Companies (many-to-one)
- Contracts → Contracted Companies (many-to-one)

## 🚀 Advanced Features

### Query Parameters

Many endpoints support query parameters for filtering and pagination:

- **Pagination**: `limit`, `offset`
- **Search**: `title`, `name`
- **Filtering**: `status` (for printing jobs)

### Date Handling

Dates should be provided in ISO 8601 format:
- **Date only**: `2024-01-15`
- **Date and time**: `2024-01-15T10:30:00Z`

### Status Tracking

Printing jobs automatically calculate status based on delivery dates:
- **pending**: Future delivery date
- **overdue**: Past delivery date
- **in_progress**: Implementation-dependent
- **completed**: Implementation-dependent

## 🔍 Troubleshooting

### Common Issues

1. **Swagger UI not loading**
   - Check that the server is running on the correct port
   - Verify the `/docs/` endpoint is accessible
   - Check browser console for JavaScript errors

2. **API calls failing**
   - Verify the server is running and healthy (`/health`)
   - Check request format matches the schema
   - Ensure required fields are provided

3. **Documentation not updating**
   - Regenerate docs with `make swagger`
   - Restart the server to load new documentation
   - Clear browser cache

### Debug Endpoints

- **Health Check**: `GET /health` - Overall system status
- **Database Health**: `GET /health/db` - Database connectivity
- **Application Health**: `GET /health/app` - Application services
- **Server Status**: `GET /status` - Server runtime information

## 📚 Additional Resources

- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [Go Swag Documentation](https://github.com/swaggo/swag)

## 🤝 Contributing

When adding new endpoints or modifying existing ones:

1. Update the `swagger.yaml` file
2. Add appropriate examples and descriptions
3. Test the changes in Swagger UI
4. Ensure all required fields are documented
5. Add proper error response documentation

## 📞 Support

For API documentation issues:
- Check the GitHub repository for the latest updates
- Use the health endpoints to verify system status
- Review the server logs for detailed error information

---

**Happy coding!** 🚀📚