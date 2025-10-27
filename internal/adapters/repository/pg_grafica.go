package repository

import (
	"context"
	"database/sql"
	"fmt"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgGraficaRepository struct {
	db *sql.DB
}

// NewPgGraficaRepository creates a new PostgreSQL-based Grafica repository
func NewPgGraficaRepository(dbService database.Service) domain.GraficaRepository {
	return &pgGraficaRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new grafica into the database
func (r *pgGraficaRepository) Save(ctx context.Context, grafica *domain.Grafica) error {
	query := `
		INSERT INTO Grafica (nome)
		VALUES ($1)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, 
		grafica.Nome,
	).Scan(&grafica.ID)
	if err != nil {
		return fmt.Errorf("failed to save grafica: %w", err)
	}
	return nil
}

// FindByID retrieves a grafica by its ID
func (r *pgGraficaRepository) FindByID(ctx context.Context, id int) (*domain.Grafica, error) {
	query := `
		SELECT id, nome
		FROM Grafica
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	
	var grafica domain.Grafica
	err := row.Scan(
		&grafica.ID,
		&grafica.Nome,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find grafica by ID: %w", err)
	}
	
	return &grafica, nil
}

// FindByName retrieves graficas by name (partial match)
func (r *pgGraficaRepository) FindByName(ctx context.Context, name string) ([]*domain.Grafica, error) {
	query := `
		SELECT id, nome
		FROM Grafica
		WHERE nome ILIKE '%' || $1 || '%'
		ORDER BY nome
	`
	rows, err := r.db.QueryContext(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to find graficas by name: %w", err)
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

// FindAll retrieves all graficas
func (r *pgGraficaRepository) FindAll(ctx context.Context) ([]*domain.Grafica, error) {
	query := `
		SELECT id, nome
		FROM Grafica
		ORDER BY nome
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all graficas: %w", err)
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

// Update modifies an existing grafica
func (r *pgGraficaRepository) Update(ctx context.Context, grafica *domain.Grafica) error {
	query := `
		UPDATE Grafica
		SET nome = $2
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		grafica.ID,
		grafica.Nome,
	)
	if err != nil {
		return fmt.Errorf("failed to update grafica: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("grafica with ID %d not found", grafica.ID)
	}
	
	return nil
}

// Delete removes a grafica from the database
func (r *pgGraficaRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM Grafica WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete grafica: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("grafica with ID %d not found", id)
	}
	
	return nil
}