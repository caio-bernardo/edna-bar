package aplica_oferta

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

// NOVO TIPO: Para o resultado da consulta com JOIN.
type AplicaOfertaDetail struct {
	model.AplicaOferta
	NomeOferta         string   `json:"nome_oferta"`
	ValorFixo          *float64 `json:"valor_fixo"`
	PercentualDesconto *int     `json:"percentual_desconto"`
}

// Busca todas as ofertas aplicadas a uma venda específica.
func (s *Store) GetByVendaID(ctx context.Context, idVenda int64) ([]AplicaOfertaDetail, error) {
	query := `
		SELECT
			ao.id_aplica_oferta, ao.id_oferta, ao.id_venda, ao.id_item_venda,
			o.nome, o.valor_fixo, o.percentual_desconto
		FROM
			aplica_oferta ao
		JOIN
			Oferta o ON ao.id_oferta = o.id_oferta
		WHERE
			ao.id_venda = $1;
	`
	rows, err := s.db.QueryContext(ctx, query, idVenda)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ofertas []AplicaOfertaDetail
	for rows.Next() {
		var o AplicaOfertaDetail
		err := rows.Scan(
			&o.IDAplicaOferta, &o.IDOferta, &o.IDVenda, &o.IDItemVenda,
			&o.NomeOferta, &o.ValorFixo, &o.PercentualDesconto,
		)
		if err != nil {
			return nil, err
		}
		ofertas = append(ofertas, o)
	}
	return ofertas, nil
}

// Deleta entradas com base no id_item_venda.
// Será usado pelo serviço de item_venda dentro de uma transação ao deletar um item.
func (s *Store) DeleteByItemVendaID(ctx context.Context, idItemVenda int64) error {
	query := `DELETE FROM aplica_oferta WHERE id_item_venda = $1;`
	_, err := s.db.ExecContext(ctx, query, idItemVenda)
	return err
}

func (s *Store) GetAll(ctx context.Context, filter util.Filter) ([]model.AplicaOferta, error) {

	query := `
		SELECT id_aplica_oferta, id_oferta, id_venda, id_item_venda
		FROM aplica_oferta
	`

	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, &filter, "c")
	if err != nil {
		return nil, err
	}

	aplicaOfertas := make([]model.AplicaOferta, 0)

	for rows.Next() {
		var c model.AplicaOferta

		err := rows.Scan(&c.IDAplicaOferta, &c.IDOferta, &c.IDVenda, &c.IDItemVenda)

		if err != nil {
			return nil, err
		}

		aplicaOfertas = append(aplicaOfertas, c)
	}

	return aplicaOfertas, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.AplicaOferta, error) {
	query := `
		SELECT id_aplica_oferta, id_oferta, id_venda, id_item_venda
		FROM aplica_oferta
		WHERE id_aplica_oferta = $1
	`

	row := s.db.QueryRowContext(ctx, query, id)

	var c model.AplicaOferta

	err := row.Scan(&c.IDAplicaOferta, &c.IDOferta, &c.IDVenda, &c.IDItemVenda)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}

func (s *Store) Create(ctx context.Context, c *model.AplicaOferta) error {
	query := `
		INSERT INTO aplica_oferta (id_oferta, id_venda, id_item_venda)
		VALUES ($1, $2, $3)
		RETURNING id_aplica_oferta
	`

	res := s.db.QueryRowContext(ctx, query, c.IDOferta, c.IDVenda, c.IDItemVenda)
	return res.Scan(&c.IDAplicaOferta)
}

func (s *Store) Update(ctx context.Context, c *model.AplicaOferta) error {
	query := `
		UPDATE aplica_oferta
		SET id_oferta = $2, id_venda = $3, id_item_venda = $4
		WHERE id_aplica_oferta = $1
	`

	res, err := s.db.ExecContext(ctx, query, c.IDAplicaOferta, c.IDOferta, c.IDVenda, c.IDItemVenda)

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

func (s *Store) Delete(ctx context.Context, id int64) (*model.AplicaOferta, error) {
	query := `
		DELETE FROM aplica_oferta
		WHERE id_aplica_oferta = $1
		RETURNING id_aplica_oferta, id_oferta, id_venda, id_item_venda
	`

	var a model.AplicaOferta
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&a.IDAplicaOferta, &a.IDOferta, &a.IDVenda, &a.IDItemVenda)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &a, nil
}
