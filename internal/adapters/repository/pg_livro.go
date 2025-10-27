package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgLivroRepository struct {
	db *sql.DB
}

// NewPgLivroRepository creates a new PostgreSQL-based Livro repository
func NewPgLivroRepository(dbService database.Service) domain.LivroRepository {
	return &pgLivroRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new livro into the database
func (r *pgLivroRepository) Save(ctx context.Context, livro *domain.Livro) error {
	query := `
		INSERT INTO Livro (ISBN, titulo, data_de_publicacao, editora_id)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, 
		livro.ISBN, 
		livro.Titulo, 
		livro.DataDePublicacao, 
		livro.EditoraID,
	)
	if err != nil {
		return fmt.Errorf("failed to save livro: %w", err)
	}
	return nil
}

// FindByISBN retrieves a livro by its ISBN
func (r *pgLivroRepository) FindByISBN(ctx context.Context, isbn string) (*domain.Livro, error) {
	query := `
		SELECT ISBN, titulo, data_de_publicacao, editora_id
		FROM Livro
		WHERE ISBN = $1
	`
	row := r.db.QueryRowContext(ctx, query, isbn)
	
	var livro domain.Livro
	err := row.Scan(
		&livro.ISBN,
		&livro.Titulo,
		&livro.DataDePublicacao,
		&livro.EditoraID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find livro by ISBN: %w", err)
	}
	
	return &livro, nil
}

// FindByTitle retrieves livros by title (partial match)
func (r *pgLivroRepository) FindByTitle(ctx context.Context, title string) ([]*domain.Livro, error) {
	query := `
		SELECT ISBN, titulo, data_de_publicacao, editora_id
		FROM Livro
		WHERE titulo ILIKE '%' || $1 || '%'
		ORDER BY titulo
	`
	rows, err := r.db.QueryContext(ctx, query, title)
	if err != nil {
		return nil, fmt.Errorf("failed to find livros by title: %w", err)
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

// FindByEditora retrieves livros by editora ID
func (r *pgLivroRepository) FindByEditora(ctx context.Context, editoraID int) ([]*domain.Livro, error) {
	query := `
		SELECT ISBN, titulo, data_de_publicacao, editora_id
		FROM Livro
		WHERE editora_id = $1
		ORDER BY data_de_publicacao DESC
	`
	rows, err := r.db.QueryContext(ctx, query, editoraID)
	if err != nil {
		return nil, fmt.Errorf("failed to find livros by editora: %w", err)
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

// FindByPublicationDateRange retrieves livros within a date range
func (r *pgLivroRepository) FindByPublicationDateRange(ctx context.Context, start, end time.Time) ([]*domain.Livro, error) {
	query := `
		SELECT ISBN, titulo, data_de_publicacao, editora_id
		FROM Livro
		WHERE data_de_publicacao BETWEEN $1 AND $2
		ORDER BY data_de_publicacao DESC
	`
	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to find livros by date range: %w", err)
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

// FindAll retrieves all livros
func (r *pgLivroRepository) FindAll(ctx context.Context) ([]*domain.Livro, error) {
	query := `
		SELECT ISBN, titulo, data_de_publicacao, editora_id
		FROM Livro
		ORDER BY titulo
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all livros: %w", err)
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

// Update modifies an existing livro
func (r *pgLivroRepository) Update(ctx context.Context, livro *domain.Livro) error {
	query := `
		UPDATE Livro
		SET titulo = $2, data_de_publicacao = $3, editora_id = $4
		WHERE ISBN = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		livro.ISBN,
		livro.Titulo,
		livro.DataDePublicacao,
		livro.EditoraID,
	)
	if err != nil {
		return fmt.Errorf("failed to update livro: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("livro with ISBN %s not found", livro.ISBN)
	}
	
	return nil
}

// Delete removes a livro from the database
func (r *pgLivroRepository) Delete(ctx context.Context, isbn string) error {
	query := `DELETE FROM Livro WHERE ISBN = $1`
	result, err := r.db.ExecContext(ctx, query, isbn)
	if err != nil {
		return fmt.Errorf("failed to delete livro: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("livro with ISBN %s not found", isbn)
	}
	
	return nil
}