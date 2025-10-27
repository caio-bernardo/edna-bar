package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgPrintingJobRepository struct {
	db *sql.DB
}

// NewPgPrintingJobRepository creates a new PostgreSQL-based PrintingJob repository
func NewPgPrintingJobRepository(dbService database.Service) domain.PrintingJobRepository {
	return &pgPrintingJobRepository{
		db: dbService.Conn(),
	}
}

// FindBooksByGrafica retrieves all books printed by a specific grafica
func (r *pgPrintingJobRepository) FindBooksByGrafica(ctx context.Context, graficaID int) ([]*domain.Livro, error) {
	query := `
		SELECT DISTINCT l.ISBN, l.titulo, l.data_de_publicacao, l.editora_id
		FROM Livro l
		INNER JOIN Imprime i ON l.ISBN = i.lisbn
		WHERE i.grafica_id = $1
		ORDER BY l.titulo
	`
	rows, err := r.db.QueryContext(ctx, query, graficaID)
	if err != nil {
		return nil, fmt.Errorf("failed to find books by grafica: %w", err)
	}
	defer rows.Close()
	
	var livros []*domain.Livro
	for rows.Next() {
		var livro domain.Livro
		err := rows.Scan(
			&livro.ISBN,
			&livro.Titulo,
			&livro.DataDePublicacao,
			&livro.EditoraID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan livro: %w", err)
		}
		livros = append(livros, &livro)
	}
	
	return livros, nil
}

// FindGraficasByBook retrieves all graficas that print a specific book
func (r *pgPrintingJobRepository) FindGraficasByBook(ctx context.Context, isbn string) ([]*domain.Grafica, error) {
	query := `
		SELECT DISTINCT g.id, g.nome
		FROM Grafica g
		INNER JOIN Imprime i ON g.id = i.grafica_id
		WHERE i.lisbn = $1
		ORDER BY g.nome
	`
	rows, err := r.db.QueryContext(ctx, query, isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to find graficas by book: %w", err)
	}
	defer rows.Close()
	
	var graficas []*domain.Grafica
	for rows.Next() {
		var grafica domain.Grafica
		err := rows.Scan(
			&grafica.ID,
			&grafica.Nome,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan grafica: %w", err)
		}
		graficas = append(graficas, &grafica)
	}
	
	return graficas, nil
}

// GetTotalCopiesByBook retrieves the total number of copies for a specific book across all graficas
func (r *pgPrintingJobRepository) GetTotalCopiesByBook(ctx context.Context, isbn string) (int, error) {
	query := `
		SELECT COALESCE(SUM(nto_copias), 0) as total_copies
		FROM Imprime
		WHERE lisbn = $1
	`
	row := r.db.QueryRowContext(ctx, query, isbn)
	
	var totalCopies int
	err := row.Scan(&totalCopies)
	if err != nil {
		return 0, fmt.Errorf("failed to get total copies by book: %w", err)
	}
	
	return totalCopies, nil
}

// GetTotalCopiesByGrafica retrieves the total number of copies for a specific grafica across all books
func (r *pgPrintingJobRepository) GetTotalCopiesByGrafica(ctx context.Context, graficaID int) (int, error) {
	query := `
		SELECT COALESCE(SUM(nto_copias), 0) as total_copies
		FROM Imprime
		WHERE grafica_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, graficaID)
	
	var totalCopies int
	err := row.Scan(&totalCopies)
	if err != nil {
		return 0, fmt.Errorf("failed to get total copies by grafica: %w", err)
	}
	
	return totalCopies, nil
}

// CreatePrintingJob creates a new printing job relationship
func (r *pgPrintingJobRepository) CreatePrintingJob(ctx context.Context, isbn string, graficaID int, copies int, deliveryDate time.Time) error {
	// First check if the relationship already exists
	checkQuery := `SELECT 1 FROM Imprime WHERE lisbn = $1 AND grafica_id = $2`
	var exists int
	err := r.db.QueryRowContext(ctx, checkQuery, isbn, graficaID).Scan(&exists)
	if err == nil {
		return fmt.Errorf("printing job between book %s and grafica %d already exists", isbn, graficaID)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing printing job: %w", err)
	}
	
	// Check if book exists
	bookQuery := `SELECT 1 FROM Livro WHERE ISBN = $1`
	err = r.db.QueryRowContext(ctx, bookQuery, isbn).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("book with ISBN %s not found", isbn)
		}
		return fmt.Errorf("failed to check book existence: %w", err)
	}
	
	// Check if grafica exists
	graficaQuery := `SELECT 1 FROM Grafica WHERE id = $1`
	err = r.db.QueryRowContext(ctx, graficaQuery, graficaID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("grafica with ID %d not found", graficaID)
		}
		return fmt.Errorf("failed to check grafica existence: %w", err)
	}
	
	// Validate input parameters
	if copies <= 0 {
		return fmt.Errorf("number of copies must be greater than zero")
	}
	
	if deliveryDate.IsZero() {
		return fmt.Errorf("delivery date cannot be zero")
	}
	
	// Create the printing job
	insertQuery := `
		INSERT INTO Imprime (lisbn, grafica_id, nto_copias, data_entrega)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.db.ExecContext(ctx, insertQuery, isbn, graficaID, copies, deliveryDate)
	if err != nil {
		return fmt.Errorf("failed to create printing job: %w", err)
	}
	
	return nil
}