package repository

import (
	"context"
	"database/sql"
	"fmt"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgContratoRepository struct {
	db *sql.DB
}

// NewPgContratoRepository creates a new PostgreSQL-based Contrato repository
func NewPgContratoRepository(dbService database.Service) domain.ContratoRepository {
	return &pgContratoRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new contrato into the database
func (r *pgContratoRepository) Save(ctx context.Context, contrato *domain.Contrato) error {
	query := `
		INSERT INTO Contrato (valor, nome_responsavel, grafica_cont_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, 
		contrato.Valor, 
		contrato.NomeResponsavel,
		contrato.GraficaContID,
	).Scan(&contrato.ID)
	if err != nil {
		return fmt.Errorf("failed to save contrato: %w", err)
	}
	return nil
}

// FindByID retrieves a contrato by its ID
func (r *pgContratoRepository) FindByID(ctx context.Context, id int) (*domain.Contrato, error) {
	query := `
		SELECT id, valor, nome_responsavel, grafica_cont_id
		FROM Contrato
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	
	var contrato domain.Contrato
	err := row.Scan(
		&contrato.ID,
		&contrato.Valor,
		&contrato.NomeResponsavel,
		&contrato.GraficaContID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find contrato by ID: %w", err)
	}
	
	return &contrato, nil
}

// FindByGraficaContID retrieves contratos by grafica contract ID
func (r *pgContratoRepository) FindByGraficaContID(ctx context.Context, graficaContID int) ([]*domain.Contrato, error) {
	query := `
		SELECT id, valor, nome_responsavel, grafica_cont_id
		FROM Contrato
		WHERE grafica_cont_id = $1
		ORDER BY id
	`
	rows, err := r.db.QueryContext(ctx, query, graficaContID)
	if err != nil {
		return nil, fmt.Errorf("failed to find contratos by grafica cont ID: %w", err)
	}
	defer rows.Close()
	
	var contratos []*domain.Contrato
	for rows.Next() {
		var contrato domain.Contrato
		err := rows.Scan(
			&contrato.ID,
			&contrato.Valor,
			&contrato.NomeResponsavel,
			&contrato.GraficaContID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contrato: %w", err)
		}
		contratos = append(contratos, &contrato)
	}
	
	return contratos, nil
}

// FindByResponsavel retrieves contratos by responsible person name (partial match)
func (r *pgContratoRepository) FindByResponsavel(ctx context.Context, nomeResponsavel string) ([]*domain.Contrato, error) {
	query := `
		SELECT id, valor, nome_responsavel, grafica_cont_id
		FROM Contrato
		WHERE nome_responsavel ILIKE '%' || $1 || '%'
		ORDER BY nome_responsavel
	`
	rows, err := r.db.QueryContext(ctx, query, nomeResponsavel)
	if err != nil {
		return nil, fmt.Errorf("failed to find contratos by responsavel: %w", err)
	}
	defer rows.Close()
	
	var contratos []*domain.Contrato
	for rows.Next() {
		var contrato domain.Contrato
		err := rows.Scan(
			&contrato.ID,
			&contrato.Valor,
			&contrato.NomeResponsavel,
			&contrato.GraficaContID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contrato: %w", err)
		}
		contratos = append(contratos, &contrato)
	}
	
	return contratos, nil
}

// FindByValueRange retrieves contratos within a value range
func (r *pgContratoRepository) FindByValueRange(ctx context.Context, minValue, maxValue float64) ([]*domain.Contrato, error) {
	query := `
		SELECT id, valor, nome_responsavel, grafica_cont_id
		FROM Contrato
		WHERE valor BETWEEN $1 AND $2
		ORDER BY valor
	`
	rows, err := r.db.QueryContext(ctx, query, minValue, maxValue)
	if err != nil {
		return nil, fmt.Errorf("failed to find contratos by value range: %w", err)
	}
	defer rows.Close()
	
	var contratos []*domain.Contrato
	for rows.Next() {
		var contrato domain.Contrato
		err := rows.Scan(
			&contrato.ID,
			&contrato.Valor,
			&contrato.NomeResponsavel,
			&contrato.GraficaContID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contrato: %w", err)
		}
		contratos = append(contratos, &contrato)
	}
	
	return contratos, nil
}

// FindAll retrieves all contratos
func (r *pgContratoRepository) FindAll(ctx context.Context) ([]*domain.Contrato, error) {
	query := `
		SELECT id, valor, nome_responsavel, grafica_cont_id
		FROM Contrato
		ORDER BY id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all contratos: %w", err)
	}
	defer rows.Close()
	
	var contratos []*domain.Contrato
	for rows.Next() {
		var contrato domain.Contrato
		err := rows.Scan(
			&contrato.ID,
			&contrato.Valor,
			&contrato.NomeResponsavel,
			&contrato.GraficaContID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contrato: %w", err)
		}
		contratos = append(contratos, &contrato)
	}
	
	return contratos, nil
}

// Update modifies an existing contrato
func (r *pgContratoRepository) Update(ctx context.Context, contrato *domain.Contrato) error {
	query := `
		UPDATE Contrato
		SET valor = $2, nome_responsavel = $3, grafica_cont_id = $4
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		contrato.ID,
		contrato.Valor,
		contrato.NomeResponsavel,
		contrato.GraficaContID,
	)
	if err != nil {
		return fmt.Errorf("failed to update contrato: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("contrato with ID %d not found", contrato.ID)
	}
	
	return nil
}

// Delete removes a contrato from the database
func (r *pgContratoRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM Contrato WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contrato: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("contrato with ID %d not found", id)
	}
	
	return nil
}