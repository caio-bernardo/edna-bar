package oferta

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

func (s *Store) GetAll(ctx context.Context, filter util.Filter) ([]model.Oferta, error) {
	query := "SELECT id_oferta, nome, data_criacao, data_inicio, data_fim, valor_fixo, percentual_desconto FROM Oferta AS o"
	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, &filter, "o")
	if err != nil {
		return nil, err
	}

	ofertas := make([]model.Oferta, 0)
	for rows.Next() {
		var o model.Oferta
		err = rows.Scan(&o.Id, &o.Nome, &o.DataCriacao, &o.DataInicio, &o.DataFim, &o.ValorFixo, &o.PercentualDesconto)
		if err != nil {
			return nil, err
		}
		ofertas = append(ofertas, o)
	}
	return ofertas, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.Oferta, error) {
	query := "SELECT id_oferta, nome, data_criacao, data_inicio, data_fim, valor_fixo, percentual_desconto FROM Oferta WHERE id_oferta = $1;"
	row := s.db.QueryRowContext(ctx, query, id)

	var o model.Oferta
	err := row.Scan(&o.Id, &o.Nome, &o.DataCriacao, &o.DataInicio, &o.DataFim, &o.ValorFixo, &o.PercentualDesconto)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (s *Store) Create(ctx context.Context, props *model.Oferta) error {
	query := `
		INSERT INTO Oferta (nome, data_inicio, data_fim, valor_fixo, percentual_desconto) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id_oferta, data_criacao;`
	res := s.db.QueryRowContext(ctx, query, props.Nome, props.DataInicio, props.DataFim, props.ValorFixo, props.PercentualDesconto)
	// Escaneamos de volta o ID e a data de criação (definida por DEFAULT)
	return res.Scan(&props.Id, &props.DataCriacao)
}

func (s *Store) Update(ctx context.Context, props *model.Oferta) error {
	query := `
		UPDATE Oferta SET 
		nome = $1, data_inicio = $2, data_fim = $3, valor_fixo = $4, percentual_desconto = $5
		WHERE id_oferta = $6;`
	res, err := s.db.ExecContext(ctx, query, props.Nome, props.DataInicio, props.DataFim, props.ValorFixo, props.PercentualDesconto, props.Id)
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

func (s *Store) Delete(ctx context.Context, id int64) (*model.Oferta, error) {
	query := "DELETE FROM Oferta WHERE id_oferta = $1 RETURNING id_oferta, nome, data_criacao, data_inicio, data_fim, valor_fixo, percentual_desconto;"
	var o model.Oferta
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&o.Id, &o.Nome, &o.DataCriacao, &o.DataInicio, &o.DataFim, &o.ValorFixo, &o.PercentualDesconto)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}
