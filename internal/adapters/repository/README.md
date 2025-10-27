# Repository Layer

This directory contains the PostgreSQL repository implementations for the book printing domain.

## Overview

The repository layer provides concrete implementations of the domain repository interfaces defined in `internal/domain/repositories.go`. All repositories follow clean architecture principles and implement the domain interfaces using PostgreSQL as the persistence layer.

## Repository Implementations

### Core Entity Repositories

- **`pg_editora.go`** - Publisher repository with CRUD operations
- **`pg_autor.go`** - Author repository with CRUD operations  
- **`pg_livro.go`** - Book repository with CRUD operations and date range queries
- **`pg_grafica.go`** - Printing company repository with CRUD operations

### Specialized Entity Repositories

- **`pg_particular.go`** - Private printing company repository
- **`pg_contratada.go`** - Contracted printing company repository with address search
- **`pg_contrato.go`** - Contract repository with value range and responsible person queries

### Relationship Repositories

- **`pg_escreve.go`** - Author-Book many-to-many relationship repository
- **`pg_imprime.go`** - Book-Printing company relationship with delivery tracking

### Aggregate Repositories

- **`pg_book_authors.go`** - Aggregate operations for book-author relationships
- **`pg_printing_job.go`** - Aggregate operations for printing jobs and copy tracking

### Registry

- **`registry.go`** - Central registry for managing all repository instances with lazy loading

## Database Schema

The repositories work with the following PostgreSQL tables:

### Core Tables
- `Editora` (id, nome, endereco)
- `Autor` (RG, nome, endereco) 
- `Livro` (ISBN, titulo, data_de_publicacao, editora_id)
- `Grafica` (id, nome)

### Specialized Tables
- `Particular` (grafica_id) - Extends Grafica
- `Contratada` (grafica_id, endereco) - Extends Grafica
- `Contrato` (id, valor, nome_responsavel, grafica_cont_id)

### Relationship Tables
- `Escreve` (ISBN, RG) - Many-to-many: Livro ↔ Autor
- `Imprime` (lisbn, grafica_id, nto_copias, data_entrega) - Many-to-many: Livro ↔ Grafica

## Usage

### Using the Registry (Recommended)

```go
// Initialize database service
dbService := database.New()

// Create repository registry
repos := repository.NewRepositoryRegistry(dbService)

// Access any repository
editoraRepo := repos.Editora()
autorRepo := repos.Autor()
livroRepo := repos.Livro()

// Use repositories
editora, err := editoraRepo.FindByID(ctx, 1)
autores, err := autorRepo.FindByName(ctx, "Garcia")
```

### Using Individual Repositories

```go
dbService := database.New()

// Create specific repositories
editoraRepo := repository.NewPgEditoraRepository(dbService)
autorRepo := repository.NewPgAutorRepository(dbService)

// Use directly
editora := &domain.Editora{Nome: "Penguin", Endereco: "New York"}
err := editoraRepo.Save(ctx, editora)
```

## Repository Features

### Error Handling
- All repositories return descriptive errors with context
- SQL errors are wrapped with operation context
- Not found cases return `nil` without error (following Go conventions)

### Query Capabilities
- **Text Search**: Name/title searches use `ILIKE` for case-insensitive partial matching
- **Date Ranges**: Publication and delivery date range queries
- **Value Ranges**: Contract value range queries
- **Aggregate Functions**: Total copy counts, overdue deliveries

### Performance Considerations
- Proper indexing recommended on frequently queried columns
- Lazy loading in repository registry
- Connection pooling handled by database service
- Prepared statements via QueryContext/ExecContext

## Example Operations

### Basic CRUD
```go
// Create
editora := domain.NewEditora(0, "Nova Editora", "São Paulo, SP")
err := repos.Editora().Save(ctx, editora) // ID populated after save

// Read
editora, err := repos.Editora().FindByID(ctx, editora.ID)
editoras, err := repos.Editora().FindByName(ctx, "Nova")

// Update
editora.Endereco = "Rio de Janeiro, RJ"
err := repos.Editora().Update(ctx, editora)

// Delete
err := repos.Editora().Delete(ctx, editora.ID)
```

### Relationship Operations
```go
// Add author to book
err := repos.BookAuthors().AddAuthorToBook(ctx, "978-0123456789", "12345678901")

// Find all books by author
livros, err := repos.BookAuthors().FindBooksByAuthor(ctx, "12345678901")

// Create printing job
err := repos.PrintingJob().CreatePrintingJob(ctx, "978-0123456789", 1, 1000, deliveryDate)

// Get total copies for a book
totalCopies, err := repos.PrintingJob().GetTotalCopiesByBook(ctx, "978-0123456789")
```

### Advanced Queries
```go
// Find books by publication date range
startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
livros, err := repos.Livro().FindByPublicationDateRange(ctx, startDate, endDate)

// Find overdue printing jobs
overdueJobs, err := repos.Imprime().FindOverdueDeliveries(ctx)

// Find contracts by value range
contracts, err := repos.Contrato().FindByValueRange(ctx, 1000.00, 5000.00)
```

## Testing

Repository tests should:
1. Use a test database with the same schema
2. Run migrations before tests
3. Clean up data between tests
4. Test both success and error scenarios
5. Verify SQL query correctness

Example test structure:
```go
func TestEditoraRepository_Save(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    repo := repository.NewPgEditoraRepository(db)
    
    // Test save operation
    editora := domain.NewEditora(0, "Test Publisher", "Test Address")
    err := repo.Save(context.Background(), editora)
    
    assert.NoError(t, err)
    assert.NotZero(t, editora.ID)
}
```

## Dependencies

- `database/sql` - Standard Go SQL interface
- `github.com/jackc/pgx/v5/stdlib` - PostgreSQL driver
- `edna/internal/domain` - Domain entities and interfaces
- `edna/internal/infra/database` - Database service

## Future Enhancements

- [ ] Connection pooling optimization
- [ ] Query result caching
- [ ] Bulk operations support
- [ ] Database migration support within repositories
- [ ] Comprehensive UnitOfWork implementation for transactions
- [ ] Repository metrics and monitoring
- [ ] Query performance logging