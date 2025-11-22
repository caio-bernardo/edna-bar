package item_venda

import (
	"context"
	"database/sql"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

// Encontra um ID de Lote adequado para um produto.
func (s *Store) FindAvailableLote(ctx context.Context, idProduto int64, quantidade int64) (int64, error) {
	query := `
		SELECT
			l.id_lote
		FROM
			Lote l
		LEFT JOIN (
			-- Subconsulta para calcular a quantidade já vendida por lote
			SELECT id_lote, SUM(quantidade) as total_vendido
			FROM item_venda
			GROUP BY id_lote
		) iv ON l.id_lote = iv.id_lote
		WHERE
			l.id_produto = $1
			AND (l.validade IS NULL OR l.validade > CURRENT_DATE)
			-- Calcula o estoque restante e verifica se é suficiente
			AND (l.quantidade_inicial - l.estragados - COALESCE(iv.total_vendido, 0)) >= $2
		ORDER BY
			l.validade ASC -- Estratégia FIFO: Pega o lote que vence primeiro
		LIMIT 1;
	`
	var idLote int64
	err := s.db.QueryRowContext(ctx, query, idProduto, quantidade).Scan(&idLote)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, types.ErrNotFound
		}
		return 0, err
	}
	return idLote, nil
}

// NOVO TIPO: Para o resultado da consulta com JOIN.
type ItemVendaDetail struct {
	model.ItemVenda
	NomeProduto string     `json:"nome_produto"`
	Marca       string     `json:"marca"`
	Validade    *time.Time `json:"validade"`
}

// Busca todos os itens de uma venda específica com detalhes do produto.
func (s *Store) GetItemsByVendaID(ctx context.Context, idVenda int64) ([]ItemVendaDetail, error) {
	query := `
		SELECT
			iv.id_item_venda, iv.id_venda, iv.id_lote, iv.quantidade, iv.valor_unitario,
			p.nome, p.marca, l.validade
		FROM
			item_venda iv
		JOIN
			Lote l ON iv.id_lote = l.id_lote
		JOIN
			Produto p ON l.id_produto = p.id_produto
		WHERE
			iv.id_venda = $1;
	`
	rows, err := s.db.QueryContext(ctx, query, idVenda)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ItemVendaDetail
	for rows.Next() {
		var i ItemVendaDetail
		err := rows.Scan(
			&i.IDItemVenda, &i.IDVenda, &i.IDLote, &i.Quantidade, &i.ValorUnitario,
			&i.NomeProduto, &i.Marca, &i.Validade,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (s *Store) GetAll(ctx context.Context, filter util.Filter) ([]model.ItemVenda, error) {
	query := "SELECT id_item_venda, id_venda, id_lote, quantidade, valor_unitario FROM item_venda AS IV"

	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, &filter, "IV")
	if err != nil {
		return nil, err
	}

	itensVenda := make([]model.ItemVenda, 0)
	for rows.Next() {
		var iv model.ItemVenda
		err = rows.Scan(&iv.IDItemVenda, &iv.IDVenda, &iv.IDLote, &iv.Quantidade, &iv.ValorUnitario)
		if err != nil {
			return nil, err
		}
		itensVenda = append(itensVenda, iv)
	}
	return itensVenda, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.ItemVenda, error) {
	query := "SELECT id_item_venda, id_venda, id_lote, quantidade, valor_unitario FROM item_venda WHERE id_item_venda = $1;"
	row := s.db.QueryRowContext(ctx, query, id)

	var iv model.ItemVenda
	err := row.Scan(&iv.IDItemVenda, &iv.IDVenda, &iv.IDLote, &iv.Quantidade, &iv.ValorUnitario)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &iv, nil
}

func (s *Store) Create(ctx context.Context, props *model.ItemVenda) error {
	query := "INSERT INTO item_venda (id_venda, id_lote, quantidade, valor_unitario) VALUES ($1, $2, $3, $4) RETURNING id_item_venda;"
	res := s.db.QueryRowContext(ctx, query, props.IDVenda, props.IDLote, props.Quantidade, props.ValorUnitario)
	return res.Scan(&props.IDItemVenda)
}

func (s *Store) Update(ctx context.Context, props *model.ItemVenda) error {
	query := "UPDATE item_venda SET id_venda = $1, id_lote = $2, quantidade = $3, valor_unitario = $4 WHERE id_item_venda = $5;"
	res, err := s.db.ExecContext(ctx, query, props.IDVenda, props.IDLote, props.Quantidade, props.ValorUnitario, props.IDItemVenda)
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

func (s *Store) Delete(ctx context.Context, id int64) (*model.ItemVenda, error) {
	query := "DELETE FROM item_venda WHERE id_item_venda = $1 RETURNING id_item_venda, id_venda, id_lote, quantidade, valor_unitario;"
	var iv model.ItemVenda
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&iv.IDItemVenda, &iv.IDVenda, &iv.IDLote, &iv.Quantidade, &iv.ValorUnitario)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &iv, nil
}
