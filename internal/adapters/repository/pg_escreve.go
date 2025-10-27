package repository

import (
	"context"
	"database/sql"
	"fmt"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgEscreveRepository struct {
	db *sql.DB
}

// NewPgEscreveRepository creates a new PostgreSQL-based Escreve repository
func NewPgEscreveRepository(dbService database.Service) domain.EscreveRepository {
	return &pgEscreveRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new escreve relationship into the database
func (r *pgEscreveRepository) Save(ctx context.Context, escreve *domain.Escreve) error {
	query := `
		INSERT INTO Escreve (ISBN, RG)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, 
		escreve.ISBN, 
		escreve.RG,
	)
	if err != nil {
		return fmt.Errorf("failed to save escreve relationship: %w", err)
	}
	return nil
}

// FindByISBN retrieves all escreve relationships for a specific book
func (r *pgEscreveRepository) FindByISBN(ctx context.Context, isbn string) ([]*domain.Escreve, error) {
	query := `
		SELECT ISBN, RG
		FROM Escreve
		WHERE ISBN = $1
		ORDER BY RG
	`
	rows, err := r.db.QueryContext(ctx, query, isbn)
	if err != nil {
		return nil, fmt.Errorf("failed to find escreve by ISBN: %w", err)
	}
	defer rows.Close()
	
	var escreves []*domain.Escreve
	for rows.Next() {
		var escreve domain.Escreve
		err := rows.Scan(
			&escreve.ISBN,
			&escreve.RG,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan escreve: %w", err)
		}
		escreves = append(escreves, &escreve)
	}
	
	return escreves, nil
}

// FindByRG retrieves all escreve relationships for a specific author
func (r *pgEscreveRepository) FindByRG(ctx context.Context, rg string) ([]*domain.Escreve, error) {
	query := `
		SELECT ISBN, RG
		FROM Escreve
		WHERE RG = $1
		ORDER BY ISBN
	`
	rows, err := r.db.QueryContext(ctx, query, rg)
	if err != nil {
		return nil, fmt.Errorf("failed to find escreve by RG: %w", err)
	}
	defer rows.Close()
	
	var escreves []*domain.Escreve
	for rows.Next() {
		var escreve domain.Escreve
		err := rows.Scan(
			&escreve.ISBN,
			&escreve.RG,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan escreve: %w", err)
		}
		escreves = append(escreves, &escreve)
	}
	
	return escreves, nil
}

// FindByISBNAndRG retrieves a specific escreve relationship
func (r *pgEscreveRepository) FindByISBNAndRG(ctx context.Context, isbn, rg string) (*domain.Escreve, error) {
	query := `
		SELECT ISBN, RG
		FROM Escreve
		WHERE ISBN = $1 AND RG = $2
	`
	row := r.db.QueryRowContext(ctx, query, isbn, rg)
	
	var escreve domain.Escreve
	err := row.Scan(
		&escreve.ISBN,
		&escreve.RG,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find escreve by ISBN and RG: %w", err)
	}
	
	return &escreve, nil
}

// FindAll retrieves all escreve relationships
func (r *pgEscreveRepository) FindAll(ctx context.Context) ([]*domain.Escreve, error) {
	query := `
		SELECT ISBN, RG
		FROM Escreve
		ORDER BY ISBN, RG
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all escreve relationships: %w", err)
	}
	defer rows.Close()
	
	var escreves []*domain.Escreve
	for rows.Next() {
		var escreve domain.Escreve
		err := rows.Scan(
			&escreve.ISBN,
			&escreve.RG,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan escreve: %w", err)
		}
		escreves = append(escreves, &escreve)
	}
	
	return escreves, nil
}

// Delete removes a specific escreve relationship
func (r *pgEscreveRepository) Delete(ctx context.Context, isbn, rg string) error {
	query := `DELETE FROM Escreve WHERE ISBN = $1 AND RG = $2`
	result, err := r.db.ExecContext(ctx, query, isbn, rg)
	if err != nil {
		return fmt.Errorf("failed to delete escreve relationship: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("escreve relationship with ISBN %s and RG %s not found", isbn, rg)
	}
	
	return nil
}

// DeleteByISBN removes all escreve relationships for a specific book
func (r *pgEscreveRepository) DeleteByISBN(ctx context.Context, isbn string) error {
	query := `DELETE FROM Escreve WHERE ISBN = $1`
	_, err := r.db.ExecContext(ctx, query, isbn)
	if err != nil {
		return fmt.Errorf("failed to delete escreve relationships by ISBN: %w", err)
	}
	
	return nil
}

// DeleteByRG removes all escreve relationships for a specific author
func (r *pgEscreveRepository) DeleteByRG(ctx context.Context, rg string) error {
	query := `DELETE FROM Escreve WHERE RG = $1`
	_, err := r.db.ExecContext(ctx, query, rg)
	if err != nil {
		return fmt.Errorf("failed to delete escreve relationships by RG: %w", err)
	}
	
	return nil
}