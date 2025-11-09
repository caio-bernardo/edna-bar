package lote

import (
	"context"
	"database/sql"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetAll(ctx context.Context, filter util.Filter) ([]model.Lote, error) {
	query := "SELECT id_lote, id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial FROM Lote AS l"
	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, &filter, "l")
	if err != nil {
		return nil, err
	}

	lotes := make([]model.Lote, 0)
	for rows.Next() {
		var l model.Lote
		err = rows.Scan(&l.Id, &l.IdFornecedor, &l.IdProduto, &l.DataFornecimento, &l.Validade, &l.PrecoUnitario, &l.Estragados, &l.QuantidadeInicial)
		if err != nil {
			return nil, err
		}
		lotes = append(lotes, l)
	}
	return lotes, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.Lote, error) {
	query := "SELECT id_lote, id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial FROM Lote WHERE id_lote = $1;"
	row := s.db.QueryRowContext(ctx, query, id)

	var l model.Lote
	err := row.Scan(&l.Id, &l.IdFornecedor, &l.IdProduto, &l.DataFornecimento, &l.Validade, &l.PrecoUnitario, &l.Estragados, &l.QuantidadeInicial)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (s *Store) Create(ctx context.Context, props *model.Lote) error {
	query := `
		INSERT INTO Lote (id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id_lote;`
	res := s.db.QueryRowContext(ctx, query, props.IdFornecedor, props.IdProduto, props.DataFornecimento, props.Validade, props.PrecoUnitario, props.Estragados, props.QuantidadeInicial)
	return res.Scan(&props.Id)
}

func (s *Store) Update(ctx context.Context, props *model.Lote) error {
	query := `
		UPDATE Lote SET 
		id_fornecedor = $1, id_produto = $2, data_fornecimento = $3, validade = $4, 
		preco_unitario = $5, estragados = $6, quantidade_inicial = $7
		WHERE id_lote = $8;`
	res, err := s.db.ExecContext(ctx, query, props.IdFornecedor, props.IdProduto, props.DataFornecimento, props.Validade, props.PrecoUnitario, props.Estragados, props.QuantidadeInicial, props.Id)
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

func (s *Store) Delete(ctx context.Context, id int64) (*model.Lote, error) {
	query := "DELETE FROM Lote WHERE id_lote = $1 RETURNING id_lote, id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial;"
	var l model.Lote
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&l.Id, &l.IdFornecedor, &l.IdProduto, &l.DataFornecimento, &l.Validade, &l.PrecoUnitario, &l.Estragados, &l.QuantidadeInicial)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}
