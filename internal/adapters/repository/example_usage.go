package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

// ExampleUsage demonstrates how to use the repository implementations
func ExampleUsage() {
	// Initialize database service
	dbService := database.New()
	defer dbService.Close()

	// Create repository registry
	repos := NewRepositoryRegistry(dbService)

	ctx := context.Background()

	// Example 1: Basic CRUD operations
	fmt.Println("=== Basic CRUD Operations ===")
	exampleBasicCRUD(ctx, repos)

	// Example 2: Relationship operations
	fmt.Println("\n=== Relationship Operations ===")
	exampleRelationshipOperations(ctx, repos)

	// Example 3: Advanced queries
	fmt.Println("\n=== Advanced Queries ===")
	exampleAdvancedQueries(ctx, repos)

	// Example 4: Printing job management
	fmt.Println("\n=== Printing Job Management ===")
	examplePrintingJobManagement(ctx, repos)
}

func exampleBasicCRUD(ctx context.Context, repos *RepositoryRegistry) {
	// Create a new publisher
	editora := domain.NewEditora(0, "Penguin Random House", "New York, NY")
	err := repos.Editora().Save(ctx, editora)
	if err != nil {
		log.Printf("Error saving editora: %v", err)
		return
	}
	fmt.Printf("Created editora with ID: %d\n", editora.ID)

	// Create a new author
	autor := domain.NewAutor("12345678901", "Gabriel García Márquez", "Colombia")
	err = repos.Autor().Save(ctx, autor)
	if err != nil {
		log.Printf("Error saving autor: %v", err)
		return
	}
	fmt.Printf("Created autor: %s\n", autor.Nome)

	// Create a new book
	publicationDate := time.Date(1967, 6, 5, 0, 0, 0, 0, time.UTC)
	livro := domain.NewLivro("978-0060883287", "One Hundred Years of Solitude", publicationDate, editora.ID)
	err = repos.Livro().Save(ctx, livro)
	if err != nil {
		log.Printf("Error saving livro: %v", err)
		return
	}
	fmt.Printf("Created livro: %s\n", livro.Titulo)

	// Read operations
	foundEditora, err := repos.Editora().FindByID(ctx, editora.ID)
	if err != nil {
		log.Printf("Error finding editora: %v", err)
		return
	}
	fmt.Printf("Found editora: %s\n", foundEditora.Nome)

	// Update operation
	foundEditora.Endereco = "London, UK"
	err = repos.Editora().Update(ctx, foundEditora)
	if err != nil {
		log.Printf("Error updating editora: %v", err)
		return
	}
	fmt.Printf("Updated editora address to: %s\n", foundEditora.Endereco)

	// Search operations
	autoresFound, err := repos.Autor().FindByName(ctx, "García")
	if err != nil {
		log.Printf("Error searching autores: %v", err)
		return
	}
	fmt.Printf("Found %d autores with 'García' in name\n", len(autoresFound))
}

func exampleRelationshipOperations(ctx context.Context, repos *RepositoryRegistry) {
	// Assuming we have a book and author from the previous example
	isbn := "978-0060883287"
	rg := "12345678901"

	// Add author to book relationship
	err := repos.BookAuthors().AddAuthorToBook(ctx, isbn, rg)
	if err != nil {
		log.Printf("Error adding author to book: %v", err)
		return
	}
	fmt.Printf("Added author %s to book %s\n", rg, isbn)

	// Find all authors for a book
	autores, err := repos.BookAuthors().FindAuthorsByBook(ctx, isbn)
	if err != nil {
		log.Printf("Error finding authors by book: %v", err)
		return
	}
	fmt.Printf("Book %s has %d authors\n", isbn, len(autores))

	// Find all books by an author
	livros, err := repos.BookAuthors().FindBooksByAuthor(ctx, rg)
	if err != nil {
		log.Printf("Error finding books by author: %v", err)
		return
	}
	fmt.Printf("Author %s has written %d books\n", rg, len(livros))
}

func exampleAdvancedQueries(ctx context.Context, repos *RepositoryRegistry) {
	// Date range query for books
	startDate := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(1970, 12, 31, 0, 0, 0, 0, time.UTC)
	
	livros, err := repos.Livro().FindByPublicationDateRange(ctx, startDate, endDate)
	if err != nil {
		log.Printf("Error finding books by date range: %v", err)
		return
	}
	fmt.Printf("Found %d books published between 1960-1970\n", len(livros))

	// Find overdue printing jobs
	overdueJobs, err := repos.Imprime().FindOverdueDeliveries(ctx)
	if err != nil {
		log.Printf("Error finding overdue deliveries: %v", err)
		return
	}
	fmt.Printf("Found %d overdue printing jobs\n", len(overdueJobs))

	// Find pending deliveries
	pendingJobs, err := repos.Imprime().FindPendingDeliveries(ctx)
	if err != nil {
		log.Printf("Error finding pending deliveries: %v", err)
		return
	}
	fmt.Printf("Found %d pending printing jobs\n", len(pendingJobs))

	// Contract value range query
	contracts, err := repos.Contrato().FindByValueRange(ctx, 1000.00, 10000.00)
	if err != nil {
		log.Printf("Error finding contracts by value range: %v", err)
		return
	}
	fmt.Printf("Found %d contracts with value between $1,000-$10,000\n", len(contracts))
}

func examplePrintingJobManagement(ctx context.Context, repos *RepositoryRegistry) {
	// Create a printing company
	grafica := domain.NewGrafica(0, "PrintTech Solutions")
	err := repos.Grafica().Save(ctx, grafica)
	if err != nil {
		log.Printf("Error saving grafica: %v", err)
		return
	}
	fmt.Printf("Created grafica with ID: %d\n", grafica.ID)

	// Create a contracted printing company
	contratada := domain.NewContratada(grafica.ID, "Industrial District, São Paulo")
	err = repos.Contratada().Save(ctx, contratada)
	if err != nil {
		log.Printf("Error saving contratada: %v", err)
		return
	}
	fmt.Printf("Created contratada for grafica ID: %d\n", contratada.GraficaID)

	// Create a contract
	contrato := &domain.Contrato{
		Valor:           5000.00,
		NomeResponsavel: "João Silva",
		GraficaContID:   grafica.ID,
	}
	err = repos.Contrato().Save(ctx, contrato)
	if err != nil {
		log.Printf("Error saving contrato: %v", err)
		return
	}
	fmt.Printf("Created contract with ID: %d, value: $%.2f\n", contrato.ID, contrato.Valor)

	// Create a printing job
	isbn := "978-0060883287"
	deliveryDate := time.Now().Add(30 * 24 * time.Hour) // 30 days from now
	
	err = repos.PrintingJob().CreatePrintingJob(ctx, isbn, grafica.ID, 1000, deliveryDate)
	if err != nil {
		log.Printf("Error creating printing job: %v", err)
		return
	}
	fmt.Printf("Created printing job: %s -> Grafica %d (1000 copies)\n", isbn, grafica.ID)

	// Get total copies for the book
	totalCopies, err := repos.PrintingJob().GetTotalCopiesByBook(ctx, isbn)
	if err != nil {
		log.Printf("Error getting total copies: %v", err)
		return
	}
	fmt.Printf("Total copies of book %s: %d\n", isbn, totalCopies)

	// Get total copies printed by the grafica
	graficaCopies, err := repos.PrintingJob().GetTotalCopiesByGrafica(ctx, grafica.ID)
	if err != nil {
		log.Printf("Error getting grafica copies: %v", err)
		return
	}
	fmt.Printf("Total copies printed by grafica %d: %d\n", grafica.ID, graficaCopies)

	// Find all books printed by this grafica
	livrosByGrafica, err := repos.PrintingJob().FindBooksByGrafica(ctx, grafica.ID)
	if err != nil {
		log.Printf("Error finding books by grafica: %v", err)
		return
	}
	fmt.Printf("Grafica %d prints %d different books\n", grafica.ID, len(livrosByGrafica))

	// Find all graficas that print a specific book
	graficasByBook, err := repos.PrintingJob().FindGraficasByBook(ctx, isbn)
	if err != nil {
		log.Printf("Error finding graficas by book: %v", err)
		return
	}
	fmt.Printf("Book %s is printed by %d graficas\n", isbn, len(graficasByBook))
}

// ExampleTransactionPattern demonstrates how you might implement transaction-like operations
// Note: This is a simplified example since we don't have UnitOfWork implemented
func ExampleTransactionPattern(ctx context.Context, repos *RepositoryRegistry) error {
	// When you need multiple operations to succeed or fail together,
	// you would typically use a transaction. For now, we show the pattern
	// that would be used with a proper UnitOfWork implementation.

	// Step 1: Create a book and its author relationship
	livro := &domain.Livro{
		ISBN:             "978-1234567890",
		Titulo:           "Example Book",
		DataDePublicacao: time.Now(),
		EditoraID:        1, // Assuming editora exists
	}

	// In a transaction, all these operations would succeed or fail together
	err := repos.Livro().Save(ctx, livro)
	if err != nil {
		return fmt.Errorf("failed to save book: %w", err)
	}

	// Add author relationship
	err = repos.BookAuthors().AddAuthorToBook(ctx, livro.ISBN, "12345678901")
	if err != nil {
		// In a real transaction, we would rollback the book creation
		repos.Livro().Delete(ctx, livro.ISBN) // Cleanup
		return fmt.Errorf("failed to add author to book: %w", err)
	}

	// Create printing job
	err = repos.PrintingJob().CreatePrintingJob(ctx, livro.ISBN, 1, 500, time.Now().Add(15*24*time.Hour))
	if err != nil {
		// In a real transaction, we would rollback everything
		repos.BookAuthors().RemoveAuthorFromBook(ctx, livro.ISBN, "12345678901")
		repos.Livro().Delete(ctx, livro.ISBN)
		return fmt.Errorf("failed to create printing job: %w", err)
	}

	fmt.Println("Successfully created book with author and printing job")
	return nil
}