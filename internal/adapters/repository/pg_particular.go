package repository

import (
	"context"
	"database/sql"
	"fmt"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgParticularRepository struct {
	db *sql.DB
}

// NewPgParticularRepository creates a new PostgreSQL-based Particular repository
func NewPgParticularRepository(dbService database.Service) domain.ParticularRepository {
	return &pgParticularRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new particular into the database
func (r *pgParticularRepository) Save(ctx context.Context, particular *domain.Particular) error {
	query := `
		INSERT INTO Particular (grafica_id)
		VALUES ($1)
	`
	_, err := r.db.ExecContext(ctx, query, particular.GraficaID)
	if err != nil {
		return fmt.Errorf("failed to save particular: %w", err)
	}
	return nil
}

// FindByGraficaID retrieves a particular by grafica ID
func (r *pgParticularRepository) FindByGraficaID(ctx context.Context, graficaID int) (*domain.Particular, error) {
	query := `
		SELECT grafica_id
		FROM Particular
		WHERE grafica_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, graficaID)
	
	var particular domain.Particular
	err := row.Scan(&particular.GraficaID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find particular by grafica ID: %w", err)
	}
	
	return &particular, nil
}

// FindAll retrieves all particulares
func (r *pgParticularRepository) FindAll(ctx context.Context) ([]*domain.Particular, error) {
	query := `
		SELECT grafica_id
		FROM Particular
		ORDER BY grafica_id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all particulares: %w", err)
	}
	defer rows.Close()
	
	var particulares []*domain.Particular
	for rows.Next() {
		var particular domain.Particular
		err := rows.Scan(&particular.GraficaID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan particular: %w", err)
		}
		particulares = append(particulares, &particular)
	}
	
	return particulares, nil
}

// Delete removes a particular from the database
func (r *pgParticularRepository) Delete(ctx context.Context, graficaID int) error {
	query := `DELETE FROM Particular WHERE grafica_id = $1`
	result, err := r.db.ExecContext(ctx, query, graficaID)
	if err != nil {
		return fmt.Errorf("failed to delete particular: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("particular with grafica ID %d not found", graficaID)
	}
	
	return nil
}