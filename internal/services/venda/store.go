package venda

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

func (s *Store) GetAll(ctx context.Context, filter util.Filter) ([]model.Venda, error) {

	query := "SELECT id_venda, data_hora_venda, data_hora_pagamento, tipo_pagamento FROM Venda AS v"
	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, &filter, "v")
	if err != nil {
		return nil, err
	}

	vendas := make([]model.Venda, 0)
	for rows.Next() {
		var venda model.Venda
		if err := rows.Scan(&venda.Id, &venda.DataHoraVenda, &venda.DataHoraPagamento, &venda.TipoPagamento); err != nil {
			return nil, err
		}
		vendas = append(vendas, venda)
	}
	return vendas, nil
}

func (s *Store) Create(ctx context.Context, venda *model.Venda) error {
	query := "INSERT INTO Venda (data_hora_venda, data_hora_pagamento, tipo_pagamento) VALUES ($1, $2, $3) RETURNING id_venda"
	res := s.db.QueryRowContext(ctx, query, venda.DataHoraVenda, venda.DataHoraPagamento, venda.TipoPagamento)
	return res.Scan(&venda.Id)
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.Venda, error) {
	query := "SELECT id_venda, data_hora_venda, data_hora_pagamento, tipo_pagamento FROM Venda WHERE id_venda = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	var venda model.Venda
	err := row.Scan(&venda.Id, &venda.DataHoraVenda, &venda.DataHoraPagamento, &venda.TipoPagamento)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &venda, nil
}

func (s *Store) Update(ctx context.Context, props *model.Venda) error {
	query := "UPDATE Venda SET data_hora_venda = $1, data_hora_pagamento = $2, tipo_pagamento = $3 WHERE id_venda = $4;"
	res, err := s.db.ExecContext(ctx, query, props.DataHoraVenda, props.DataHoraPagamento, props.TipoPagamento, props.Id)
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

func (s *Store) Delete(ctx context.Context, id int64) (*model.Venda, error) {
	query := "DELETE FROM Venda WHERE id_venda = $1 RETURNING id_venda, data_hora_venda, data_hora_pagamento, tipo_pagamento;"

	var venda model.Venda
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&venda.Id, &venda.DataHoraVenda, &venda.DataHoraPagamento, &venda.TipoPagamento)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &venda, nil
}
