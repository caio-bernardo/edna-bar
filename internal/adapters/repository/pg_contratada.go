package repository

import (
	"context"
	"database/sql"
	"fmt"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgContratadaRepository struct {
	db *sql.DB
}

// NewPgContratadaRepository creates a new PostgreSQL-based Contratada repository
func NewPgContratadaRepository(dbService database.Service) domain.ContratadaRepository {
	return &pgContratadaRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new contratada into the database
func (r *pgContratadaRepository) Save(ctx context.Context, contratada *domain.Contratada) error {
	query := `
		INSERT INTO Contratada (grafica_id, endereco)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, 
		contratada.GraficaID, 
		contratada.Endereco,
	)
	if err != nil {
		return fmt.Errorf("failed to save contratada: %w", err)
	}
	return nil
}

// FindByGraficaID retrieves a contratada by grafica ID
func (r *pgContratadaRepository) FindByGraficaID(ctx context.Context, graficaID int) (*domain.Contratada, error) {
	query := `
		SELECT grafica_id, endereco
		FROM Contratada
		WHERE grafica_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, graficaID)
	
	var contratada domain.Contratada
	err := row.Scan(
		&contratada.GraficaID,
		&contratada.Endereco,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find contratada by grafica ID: %w", err)
	}
	
	return &contratada, nil
}

// FindByAddress retrieves contratadas by address (partial match)
func (r *pgContratadaRepository) FindByAddress(ctx context.Context, endereco string) ([]*domain.Contratada, error) {
	query := `
		SELECT grafica_id, endereco
		FROM Contratada
		WHERE endereco ILIKE '%' || $1 || '%'
		ORDER BY endereco
	`
	rows, err := r.db.QueryContext(ctx, query, endereco)
	if err != nil {
		return nil, fmt.Errorf("failed to find contratadas by address: %w", err)
	}
	defer rows.Close()
	
	var contratadas []*domain.Contratada
	for rows.Next() {
		var contratada domain.Contratada
		err := rows.Scan(
			&contratada.GraficaID,
			&contratada.Endereco,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contratada: %w", err)
		}
		contratadas = append(contratadas, &contratada)
	}
	
	return contratadas, nil
}

// FindAll retrieves all contratadas
func (r *pgContratadaRepository) FindAll(ctx context.Context) ([]*domain.Contratada, error) {
	query := `
		SELECT grafica_id, endereco
		FROM Contratada
		ORDER BY grafica_id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all contratadas: %w", err)
	}
	defer rows.Close()
	
	var contratadas []*domain.Contratada
	for rows.Next() {
		var contratada domain.Contratada
		err := rows.Scan(
			&contratada.GraficaID,
			&contratada.Endereco,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contratada: %w", err)
		}
		contratadas = append(contratadas, &contratada)
	}
	
	return contratadas, nil
}

// Update modifies an existing contratada
func (r *pgContratadaRepository) Update(ctx context.Context, contratada *domain.Contratada) error {
	query := `
		UPDATE Contratada
		SET endereco = $2
		WHERE grafica_id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		contratada.GraficaID,
		contratada.Endereco,
	)
	if err != nil {
		return fmt.Errorf("failed to update contratada: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("contratada with grafica ID %d not found", contratada.GraficaID)
	}
	
	return nil
}

// Delete removes a contratada from the database
func (r *pgContratadaRepository) Delete(ctx context.Context, graficaID int) error {
	query := `DELETE FROM Contratada WHERE grafica_id = $1`
	result, err := r.db.ExecContext(ctx, query, graficaID)
	if err != nil {
		return fmt.Errorf("failed to delete contratada: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("contratada with grafica ID %d not found", graficaID)
	}
	
	return nil
}