package repository

import (
	"context"
	"database/sql"
	"fmt"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgEditoraRepository struct {
	db *sql.DB
}

// NewPgEditoraRepository creates a new PostgreSQL-based Editora repository
func NewPgEditoraRepository(dbService database.Service) domain.EditoraRepository {
	return &pgEditoraRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new editora into the database
func (r *pgEditoraRepository) Save(ctx context.Context, editora *domain.Editora) error {
	query := `
		INSERT INTO Editora (nome, endereco)
		VALUES ($1, $2)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, 
		editora.Nome, 
		editora.Endereco,
	).Scan(&editora.ID)
	if err != nil {
		return fmt.Errorf("failed to save editora: %w", err)
	}
	return nil
}

// FindByID retrieves an editora by its ID
func (r *pgEditoraRepository) FindByID(ctx context.Context, id int) (*domain.Editora, error) {
	query := `
		SELECT id, nome, endereco
		FROM Editora
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	
	var editora domain.Editora
	err := row.Scan(
		&editora.ID,
		&editora.Nome,
		&editora.Endereco,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find editora by ID: %w", err)
	}
	
	return &editora, nil
}

// FindByName retrieves editoras by name (partial match)
func (r *pgEditoraRepository) FindByName(ctx context.Context, name string) ([]*domain.Editora, error) {
	query := `
		SELECT id, nome, endereco
		FROM Editora
		WHERE nome ILIKE '%' || $1 || '%'
		ORDER BY nome
	`
	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find editoras by name: %w", err)
	}
	defer rows.Close()
	
	var editoras []*domain.Editora
	for rows.Next() {
		var editora domain.Editora
		err := rows.Scan(
			&editora.ID,
			&editora.Nome,
			&editora.Endereco,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan editora: %w", err)
		}
		editoras = append(editoras, &editora)
	}
	
	return editoras, nil
}

// FindAll retrieves all editoras
func (r *pgEditoraRepository) FindAll(ctx context.Context) ([]*domain.Editora, error) {
	query := `
		SELECT id, nome, endereco
		FROM Editora
		ORDER BY nome
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all editoras: %w", err)
	}
	defer rows.Close()
	
	var editoras []*domain.Editora
	for rows.Next() {
		var editora domain.Editora
		err := rows.Scan(
			&editora.ID,
			&editora.Nome,
			&editora.Endereco,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan editora: %w", err)
		}
		editoras = append(editoras, &editora)
	}
	
	return editoras, nil
}

// Update modifies an existing editora
func (r *pgEditoraRepository) Update(ctx context.Context, editora *domain.Editora) error {
	query := `
		UPDATE Editora
		SET nome = $2, endereco = $3
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		editora.ID,
		editora.Nome,
		editora.Endereco,
	)
	if err != nil {
		return fmt.Errorf("failed to update editora: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("editora with ID %d not found", editora.ID)
	}
	
	return nil
}

// Delete removes an editora from the database
func (r *pgEditoraRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM Editora WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete editora: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("editora with ID %d not found", id)
	}
	
	return nil
}