package repository

import (
	"context"
	"database/sql"
	"fmt"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgAutorRepository struct {
	db *sql.DB
}

// NewPgAutorRepository creates a new PostgreSQL-based Autor repository
func NewPgAutorRepository(dbService database.Service) domain.AutorRepository {
	return &pgAutorRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new autor into the database
func (r *pgAutorRepository) Save(ctx context.Context, autor *domain.Autor) error {
	query := `
		INSERT INTO Autor (RG, nome, endereco)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, 
		autor.RG, 
		autor.Nome, 
		autor.Endereco,
	)
	if err != nil {
		return fmt.Errorf("failed to save autor: %w", err)
	}
	return nil
}

// FindByRG retrieves an autor by their RG
func (r *pgAutorRepository) FindByRG(ctx context.Context, rg string) (*domain.Autor, error) {
	query := `
		SELECT RG, nome, endereco
		FROM Autor
		WHERE RG = $1
	`
	row := r.db.QueryRowContext(ctx, query, rg)
	
	var autor domain.Autor
	err := row.Scan(
		&autor.RG,
		&autor.Nome,
		&autor.Endereco,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find autor by RG: %w", err)
	}
	
	return &autor, nil
}

// FindByName retrieves autores by name (partial match)
func (r *pgAutorRepository) FindByName(ctx context.Context, name string) ([]*domain.Autor, error) {
	query := `
		SELECT RG, nome, endereco
		FROM Autor
		WHERE nome ILIKE '%' || $1 || '%'
		ORDER BY nome
	`
	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find autores by name: %w", err)
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

// FindAll retrieves all autores
func (r *pgAutorRepository) FindAll(ctx context.Context) ([]*domain.Autor, error) {
	query := `
		SELECT RG, nome, endereco
		FROM Autor
		ORDER BY nome
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all autores: %w", err)
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

// Update modifies an existing autor
func (r *pgAutorRepository) Update(ctx context.Context, autor *domain.Autor) error {
	query := `
		UPDATE Autor
		SET nome = $2, endereco = $3
		WHERE RG = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		autor.RG,
		autor.Nome,
		autor.Endereco,
	)
	if err != nil {
		return fmt.Errorf("failed to update autor: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("autor with RG %s not found", autor.RG)
	}
	
	return nil
}

// Delete removes an autor from the database
func (r *pgAutorRepository) Delete(ctx context.Context, rg string) error {
	query := `DELETE FROM Autor WHERE RG = $1`
	result, err := r.db.ExecContext(ctx, query, rg)
	if err != nil {
		return fmt.Errorf("failed to delete autor: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("autor with RG %s not found", rg)
	}
	
	return nil
}