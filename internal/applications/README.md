# Application Layer

The application layer contains the use cases and application services that orchestrate the business logic defined in the domain layer. This layer implements the clean architecture pattern and serves as the entry point for external interfaces.

## Overview

The application layer is responsible for:
- Coordinating business operations through use cases
- Converting between domain entities and DTOs
- Validating input data
- Orchestrating domain services
- Managing transactions and error handling

## Structure

```
applications/
├── README.md                  # This documentation
├── application_service.go     # Main application service coordinator
├── dtos.go                   # Common DTOs and response types
├── autor_usecase.go          # Author use cases
├── editora_usecase.go        # Publisher use cases
├── livro_usecase.go          # Book use cases
├── grafica_usecase.go        # Printing company use cases
├── contrato_usecase.go       # Contract use cases
└── imprime_usecase.go        # Printing job use cases
```

## Use Cases

### AutorUsecase (Author Management)
- **Create**: Create new authors with validation
- **Read**: Get authors by RG, name, or list all
- **Update**: Update author information
- **Delete**: Delete authors (with business rule validation)
- **Relationships**: Manage author-book relationships

**Key Operations:**
- `Create(ctx, CreateAutorRequest) -> AutorResponse`
- `GetWithBooks(ctx, rg) -> AutorResponse` (includes authored books)
- `AddToBook(ctx, rg, isbn)` - Add author to a book
- `RemoveFromBook(ctx, rg, isbn)` - Remove author from book

### EditoraUsecase (Publisher Management)
- **Create**: Create new publishers
- **Read**: Get publishers by ID, name, or list all
- **Update**: Update publisher information
- **Delete**: Delete publishers (validates no associated books)
- **Catalog**: View publisher's book catalog

**Key Operations:**
- `Create(ctx, CreateEditoraRequest) -> EditoraResponse`
- `GetWithBooks(ctx, id) -> EditoraResponse` (includes published books)
- `GetByName(ctx, name) -> []EditoraResponse`

### LivroUsecase (Book Management)
- **Create**: Create new books with publisher validation
- **Read**: Get books by ISBN, title, publisher, or date range
- **Update**: Update book information
- **Delete**: Delete books (removes author relationships)
- **Authors**: Manage book-author relationships

**Key Operations:**
- `Create(ctx, CreateLivroRequest) -> LivroResponse`
- `GetWithAuthors(ctx, isbn) -> LivroResponse` (includes authors)
- `GetByEditora(ctx, editoraID) -> []LivroResponse`
- `GetByPublicationDateRange(ctx, start, end) -> []LivroResponse`
- `AddAuthor(ctx, isbn, authorRG)` - Add author to book
- `RemoveAuthor(ctx, isbn, authorRG)` - Remove author from book

### GraficaUsecase (Printing Company Management)
- **Create**: Create printing companies (particular or contracted)
- **Read**: Get companies by ID, name, or type
- **Update**: Update company information
- **Delete**: Delete companies (validates no active jobs)
- **Jobs**: View printing jobs and contracts

**Key Operations:**
- `Create(ctx, CreateGraficaRequest) -> GraficaResponse`
- `GetWithContracts(ctx, id) -> GraficaResponse` (for contracted companies)
- `GetWithPrintingJobs(ctx, id) -> GraficaResponse`
- `GetByType(ctx, type) -> []GraficaResponse` (particular/contracted)

### ContratoUsecase (Contract Management)
- **Create**: Create contracts for contracted printing companies
- **Read**: Get contracts by ID, company, or responsible person
- **Update**: Update contract terms
- **Delete**: Delete contracts
- **Analysis**: Contract value analysis and reporting

**Key Operations:**
- `Create(ctx, CreateContratoRequest) -> ContratoResponse`
- `GetByGraficaContID(ctx, graficaContID) -> []ContratoResponse`
- `GetByValueRange(ctx, min, max) -> []ContratoResponse`
- `GetContractAnalysis(ctx) -> ContractAnalysisResponse`

### ImprimeUsecase (Printing Job Management)
- **Create**: Schedule new printing jobs
- **Read**: Get jobs by book, company, or delivery date
- **Update**: Update job details and delivery dates
- **Delete**: Cancel printing jobs
- **Tracking**: Monitor overdue and pending jobs
- **Analytics**: Printing statistics and reporting

**Key Operations:**
- `Create(ctx, CreateImprimeRequest) -> ImprimeResponse`
- `GetOverdueJobs(ctx) -> []ImprimeResponse`
- `GetPendingJobs(ctx) -> []ImprimeResponse`
- `GetPrintingStatistics(ctx, start, end) -> PrintingStatisticsResponse`
- `MarkAsCompleted(ctx, isbn, graficaID)`

## DTOs (Data Transfer Objects)

### Request DTOs
- Input validation and sanitization
- Field-level validation with descriptive errors
- Business rule validation

### Response DTOs
- Consistent JSON serialization
- Optional nested data (relationships)
- Status and metadata information

### Common DTOs
- `PaginationRequest/Response` - For paginated queries
- `ErrorResponse` - Standardized error format
- `SuccessResponse` - Standardized success format
- `SearchRequest` - Search parameters
- `DateRangeRequest` - Date range queries

## Error Handling

The application layer provides comprehensive error handling:

```go
// Domain errors are wrapped with context
if err := usecase.Create(ctx, request); err != nil {
    return fmt.Errorf("failed to create entity: %w", err)
}

// Validation errors with field information
if request.Name == "" {
    return domain.NewFieldRequiredError("name")
}

// Business rule violations
if len(authors) <= 1 {
    return domain.NewDomainError("CANNOT_REMOVE_LAST_AUTHOR", 
        "cannot remove the last author from a book")
}
```

## Validation

### Input Validation
- Required field validation
- Format validation (email, dates, etc.)
- Length constraints
- Business rule validation

### Business Rule Validation
- Entity relationships (foreign key validation)
- Domain invariants (book must have at least one author)
- State transitions (cannot delete entities with dependencies)

## Usage Example

```go
// Initialize application service
config := ApplicationServiceConfig{
    LivroRepo:       livroRepo,
    AutorRepo:       autorRepo,
    EditoraRepo:     editoraRepo,
    // ... other repositories
}

appService := NewApplicationService(config)

// Create a new book
request := CreateLivroRequest{
    ISBN:             "978-3-16-148410-0",
    Titulo:           "Clean Architecture",
    DataDePublicacao: time.Date(2017, 9, 20, 0, 0, 0, 0, time.UTC),
    EditoraID:        1,
}

book, err := appService.LivroUsecase.Create(ctx, request)
if err != nil {
    // Handle error
    return err
}

// Add an author to the book
err = appService.LivroUsecase.AddAuthor(ctx, book.ISBN, "12345678901")
if err != nil {
    // Handle error
    return err
}
```

## Testing

Each use case should be thoroughly tested:

```go
func TestLivroUsecase_Create(t *testing.T) {
    // Setup mocks
    mockLivroRepo := &MockLivroRepository{}
    mockEditoraRepo := &MockEditoraRepository{}
    
    // Create use case
    usecase := NewLivroUsecase(mockLivroRepo, mockEditoraRepo, ...)
    
    // Test cases
    t.Run("success", func(t *testing.T) {
        // Test successful creation
    })
    
    t.Run("validation_error", func(t *testing.T) {
        // Test validation errors
    })
    
    t.Run("business_rule_violation", func(t *testing.T) {
        // Test business rule violations
    })
}
```

## Architecture Compliance

This application layer follows clean architecture principles:

1. **Dependency Inversion**: Depends only on domain interfaces
2. **Single Responsibility**: Each use case has a single responsibility
3. **Interface Segregation**: Fine-grained repository interfaces
4. **Open/Closed**: Extensible through new use cases
5. **Separation of Concerns**: Clear separation between validation, orchestration, and business logic

## Integration

The application layer integrates with:
- **Domain Layer**: Uses domain entities, services, and repositories
- **Infrastructure Layer**: Through repository implementations
- **Presentation Layer**: Through HTTP handlers, gRPC services, etc.

## Performance Considerations

- Repository operations are kept minimal
- Lazy loading for optional relationships
- Pagination for large result sets
- Caching at the infrastructure layer
- Database transaction management through UnitOfWork pattern