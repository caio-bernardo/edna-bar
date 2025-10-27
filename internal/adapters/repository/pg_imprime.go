package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"edna/internal/domain"
	"edna/internal/infra/database"
)

type pgImprimeRepository struct {
	db *sql.DB
}

// NewPgImprimeRepository creates a new PostgreSQL-based Imprime repository
func NewPgImprimeRepository(dbService database.Service) domain.ImprimeRepository {
	return &pgImprimeRepository{
		db: dbService.Conn(),
	}
}

// Save inserts a new imprime relationship into the database
func (r *pgImprimeRepository) Save(ctx context.Context, imprime *domain.Imprime) error {
	query := `
		INSERT INTO Imprime (lisbn, grafica_id, nto_copias, data_entrega)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, 
		imprime.LISBN, 
		imprime.GraficaID,
		imprime.NtoCopias,
		imprime.DataEntrega,
	)
	if err != nil {
		return fmt.Errorf("failed to save imprime relationship: %w", err)
	}
	return nil
}

// FindByISBN retrieves all imprime relationships for a specific book
func (r *pgImprimeRepository) FindByISBN(ctx context.Context, lisbn string) ([]*domain.Imprime, error) {
	query := `
		SELECT lisbn, grafica_id, nto_copias, data_entrega
		FROM Imprime
		WHERE lisbn = $1
		ORDER BY data_entrega
	`
	rows, err := r.db.QueryContext(ctx, query, lisbn)
	if err != nil {
		return nil, fmt.Errorf("failed to find imprime by ISBN: %w", err)
	}
	defer rows.Close()
	
	var imprimes []*domain.Imprime
	for rows.Next() {
		var imprime domain.Imprime
		err := rows.Scan(
			&imprime.LISBN,
			&imprime.GraficaID,
			&imprime.NtoCopias,
			&imprime.DataEntrega,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan imprime: %w", err)
		}
		imprimes = append(imprimes, &imprime)
	}
	
	return imprimes, nil
}

// FindByGraficaID retrieves all imprime relationships for a specific grafica
func (r *pgImprimeRepository) FindByGraficaID(ctx context.Context, graficaID int) ([]*domain.Imprime, error) {
	query := `
		SELECT lisbn, grafica_id, nto_copias, data_entrega
		FROM Imprime
		WHERE grafica_id = $1
		ORDER BY data_entrega
	`
	rows, err := r.db.QueryContext(ctx, query, graficaID)
	if err != nil {
		return nil, fmt.Errorf("failed to find imprime by grafica ID: %w", err)
	}
	defer rows.Close()
	
	var imprimes []*domain.Imprime
	for rows.Next() {
		var imprime domain.Imprime
		err := rows.Scan(
			&imprime.LISBN,
			&imprime.GraficaID,
			&imprime.NtoCopias,
			&imprime.DataEntrega,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan imprime: %w", err)
		}
		imprimes = append(imprimes, &imprime)
	}
	
	return imprimes, nil
}

// FindByISBNAndGraficaID retrieves a specific imprime relationship
func (r *pgImprimeRepository) FindByISBNAndGraficaID(ctx context.Context, lisbn string, graficaID int) (*domain.Imprime, error) {
	query := `
		SELECT lisbn, grafica_id, nto_copias, data_entrega
		FROM Imprime
		WHERE lisbn = $1 AND grafica_id = $2
	`
	row := r.db.QueryRowContext(ctx, query, lisbn, graficaID)
	
	var imprime domain.Imprime
	err := row.Scan(
		&imprime.LISBN,
		&imprime.GraficaID,
		&imprime.NtoCopias,
		&imprime.DataEntrega,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find imprime by ISBN and grafica ID: %w", err)
	}
	
	return &imprime, nil
}

// FindByDeliveryDateRange retrieves imprime relationships within a delivery date range
func (r *pgImprimeRepository) FindByDeliveryDateRange(ctx context.Context, start, end time.Time) ([]*domain.Imprime, error) {
	query := `
		SELECT lisbn, grafica_id, nto_copias, data_entrega
		FROM Imprime
		WHERE data_entrega BETWEEN $1 AND $2
		ORDER BY data_entrega
	`
	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to find imprime by delivery date range: %w", err)
	}
	defer rows.Close()
	
	var imprimes []*domain.Imprime
	for rows.Next() {
		var imprime domain.Imprime
		err := rows.Scan(
			&imprime.LISBN,
			&imprime.GraficaID,
			&imprime.NtoCopias,
			&imprime.DataEntrega,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan imprime: %w", err)
		}
		imprimes = append(imprimes, &imprime)
	}
	
	return imprimes, nil
}

// FindOverdueDeliveries retrieves imprime relationships with overdue deliveries
func (r *pgImprimeRepository) FindOverdueDeliveries(ctx context.Context) ([]*domain.Imprime, error) {
	query := `
		SELECT lisbn, grafica_id, nto_copias, data_entrega
		FROM Imprime
		WHERE data_entrega < CURRENT_DATE
		ORDER BY data_entrega
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find overdue deliveries: %w", err)
	}
	defer rows.Close()
	
	var imprimes []*domain.Imprime
	for rows.Next() {
		var imprime domain.Imprime
		err := rows.Scan(
			&imprime.LISBN,
			&imprime.GraficaID,
			&imprime.NtoCopias,
			&imprime.DataEntrega,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan imprime: %w", err)
		}
		imprimes = append(imprimes, &imprime)
	}
	
	return imprimes, nil
}

// FindPendingDeliveries retrieves imprime relationships with pending deliveries
func (r *pgImprimeRepository) FindPendingDeliveries(ctx context.Context) ([]*domain.Imprime, error) {
	query := `
		SELECT lisbn, grafica_id, nto_copias, data_entrega
		FROM Imprime
		WHERE data_entrega >= CURRENT_DATE
		ORDER BY data_entrega
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending deliveries: %w", err)
	}
	defer rows.Close()
	
	var imprimes []*domain.Imprime
	for rows.Next() {
		var imprime domain.Imprime
		err := rows.Scan(
			&imprime.LISBN,
			&imprime.GraficaID,
			&imprime.NtoCopias,
			&imprime.DataEntrega,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan imprime: %w", err)
		}
		imprimes = append(imprimes, &imprime)
	}
	
	return imprimes, nil
}

// FindAll retrieves all imprime relationships
func (r *pgImprimeRepository) FindAll(ctx context.Context) ([]*domain.Imprime, error) {
	query := `
		SELECT lisbn, grafica_id, nto_copias, data_entrega
		FROM Imprime
		ORDER BY lisbn, grafica_id
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all imprime relationships: %w", err)
	}
	defer rows.Close()
	
	var imprimes []*domain.Imprime
	for rows.Next() {
		var imprime domain.Imprime
		err := rows.Scan(
			&imprime.LISBN,
			&imprime.GraficaID,
			&imprime.NtoCopias,
			&imprime.DataEntrega,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan imprime: %w", err)
		}
		imprimes = append(imprimes, &imprime)
	}
	
	return imprimes, nil
}

// Update modifies an existing imprime relationship
func (r *pgImprimeRepository) Update(ctx context.Context, imprime *domain.Imprime) error {
	query := `
		UPDATE Imprime
		SET nto_copias = $3, data_entrega = $4
		WHERE lisbn = $1 AND grafica_id = $2
	`
	result, err := r.db.ExecContext(ctx, query,
		imprime.LISBN,
		imprime.GraficaID,
		imprime.NtoCopias,
		imprime.DataEntrega,
	)
	if err != nil {
		return fmt.Errorf("failed to update imprime: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("imprime with ISBN %s and grafica ID %d not found", imprime.LISBN, imprime.GraficaID)
	}
	
	return nil
}

// Delete removes a specific imprime relationship
func (r *pgImprimeRepository) Delete(ctx context.Context, lisbn string, graficaID int) error {
	query := `DELETE FROM Imprime WHERE lisbn = $1 AND grafica_id = $2`
	result, err := r.db.ExecContext(ctx, query, lisbn, graficaID)
	if err != nil {
		return fmt.Errorf("failed to delete imprime relationship: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("imprime relationship with ISBN %s and grafica ID %d not found", lisbn, graficaID)
	}
	
	return nil
}

// DeleteByISBN removes all imprime relationships for a specific book
func (r *pgImprimeRepository) DeleteByISBN(ctx context.Context, lisbn string) error {
	query := `DELETE FROM Imprime WHERE lisbn = $1`
	_, err := r.db.ExecContext(ctx, query, lisbn)
	if err != nil {
		return fmt.Errorf("failed to delete imprime relationships by ISBN: %w", err)
	}
	
	return nil
}

// DeleteByGraficaID removes all imprime relationships for a specific grafica
func (r *pgImprimeRepository) DeleteByGraficaID(ctx context.Context, graficaID int) error {
	query := `DELETE FROM Imprime WHERE grafica_id = $1`
	_, err := r.db.ExecContext(ctx, query, graficaID)
	if err != nil {
		return fmt.Errorf("failed to delete imprime relationships by grafica ID: %w", err)
	}
	
	return nil
}