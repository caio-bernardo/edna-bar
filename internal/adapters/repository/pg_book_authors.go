package repository

import (
	"context"
	"database/sql"
	"fmt"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgBookAuthorsRepository struct {
	db *sql.DB
}

// NewPgBookAuthorsRepository creates a new PostgreSQL-based BookAuthors repository
func NewPgBookAuthorsRepository(dbService database.Service) domain.BookAuthorsRepository {
	return &pgBookAuthorsRepository{
		db: dbService.Conn(),
	}
}

// FindAuthorsByBook retrieves all authors for a specific book
func (r *pgBookAuthorsRepository) FindAuthorsByBook(ctx context.Context, isbn string) ([]*domain.Autor, error) {
	query := `
		SELECT a.RG, a.nome, a.endereco
		FROM Autor a
		INNER JOIN Escreve e ON a.RG = e.RG
		WHERE e.ISBN = $1
		ORDER BY a.nome
	`
	rows, err := r.db.QueryContext(ctx, query, isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to find authors by book: %w", err)
	}
	defer rows.Close()
	
	var autores []*domain.Autor
	for rows.Next() {
		var autor domain.Autor
		err := rows.Scan(
			&autor.RG,
			&autor.Nome,
			&autor.Endereco,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan autor: %w", err)
		}
		autores = append(autores, &autor)
	}
	
	return autores, nil
}

// FindBooksByAuthor retrieves all books for a specific author
func (r *pgBookAuthorsRepository) FindBooksByAuthor(ctx context.Context, rg string) ([]*domain.Livro, error) {
	query := `
		SELECT l.ISBN, l.titulo, l.data_de_publicacao, l.editora_id
		FROM Livro l
		INNER JOIN Escreve e ON l.ISBN = e.ISBN
		WHERE e.RG = $1
		ORDER BY l.data_de_publicacao DESC
	`
	rows, err := r.db.QueryContext(ctx, query, rg)
	if err != nil {
		return nil, fmt.Errorf("failed to find books by author: %w", err)
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

// AddAuthorToBook creates a relationship between an author and a book
func (r *pgBookAuthorsRepository) AddAuthorToBook(ctx context.Context, isbn, rg string) error {
	// First check if the relationship already exists
	checkQuery := `SELECT 1 FROM Escreve WHERE ISBN = $1 AND RG = $2`
	var exists int
	err := r.db.QueryRowContext(ctx, checkQuery, isbn, rg).Scan(&exists)
	if err == nil {
		return fmt.Errorf("relationship between book %s and author %s already exists", isbn, rg)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing relationship: %w", err)
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
	
	// Check if author exists
	authorQuery := `SELECT 1 FROM Autor WHERE RG = $1`
	err = r.db.QueryRowContext(ctx, authorQuery, rg).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("author with RG %s not found", rg)
		}
		return fmt.Errorf("failed to check author existence: %w", err)
	}
	
	// Create the relationship
	insertQuery := `INSERT INTO Escreve (ISBN, RG) VALUES ($1, $2)`
	_, err = r.db.ExecContext(ctx, insertQuery, isbn, rg)
	if err != nil {
		return fmt.Errorf("failed to add author to book: %w", err)
	}
	
	return nil
}

// RemoveAuthorFromBook removes a relationship between an author and a book
func (r *pgBookAuthorsRepository) RemoveAuthorFromBook(ctx context.Context, isbn, rg string) error {
	query := `DELETE FROM Escreve WHERE ISBN = $1 AND RG = $2`
	result, err := r.db.ExecContext(ctx, query, isbn, rg)
	if err != nil {
		return fmt.Errorf("failed to remove author from book: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("relationship between book %s and author %s not found", isbn, rg)
	}
	
	return nil
}