package fornecedor

import (
	"context"
	"database/sql"
	"edna/internal/model"
	"edna/internal/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetAll(ctx context.Context) ([]model.Fornecedor, error) {
	query := "SELECT id_fornecedor, nome, CNPJ FROM Fornecedor;"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	fornecedores := make([]model.Fornecedor, 0)
	for rows.Next() {
		var fornecedor model.Fornecedor
		err = rows.Scan(&fornecedor.Id, &fornecedor.Nome, &fornecedor.CNPJ)
		if err != nil {
			return nil, err
		}
		fornecedores = append(fornecedores, fornecedor)
	}

	return fornecedores, nil
}

func (s *Store) Create(ctx context.Context, props *model.Fornecedor) error {
	query := "INSERT INTO Fornecedor (nome, CNPJ) VALUES ($1, $2) RETURNING id_fornecedor;"

	res := s.db.QueryRowContext(ctx, query, props.Nome, props.CNPJ)
	return res.Scan(&props.Id)
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.Fornecedor, error) {
	query := "SELECT id_fornecedor, nome, CNPJ FROM Fornecedor WHERE id_fornecedor = $1;"

	row := s.db.QueryRowContext(ctx, query, id)

	var fornecedor model.Fornecedor
	err := row.Scan(&fornecedor.Id, &fornecedor.Nome, &fornecedor.CNPJ)
	if err != nil {
		return nil, err
	}

	return &fornecedor, nil
}

func (s *Store) Update(ctx context.Context, props *model.Fornecedor) error {
	query := "UPDATE Fornecedor SET nome = $1, CNPJ = $2 WHERE id_fornecedor = $3;"

	res, err := s.db.ExecContext(ctx, query, props.Nome, props.CNPJ, props.Id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return types.ErrNotFound
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM Fornecedor WHERE id_fornecedor = $1;"

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return types.ErrNotFound
	}
	return nil
}
