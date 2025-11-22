package item_oferta

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

func (s *Store) GetAll(ctx context.Context, filter util.Filter) ([]model.ItemOferta, error) {
	query := "SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta as io"

	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, &filter, "io")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itensOferta := make([]model.ItemOferta, 0)
	for rows.Next() {
		var io model.ItemOferta
		err = rows.Scan(&io.Quantidade, &io.IDProduto, &io.IDOferta)
		if err != nil {
			return nil, err
		}
		itensOferta = append(itensOferta, io)
	}
	return itensOferta, nil
}

// GetByComposedID busca uma entrada específica de ItemOferta pela sua chave primária composta.
func (s *Store) GetByComposedID(ctx context.Context, id_produto int64, id_oferta int64) (*model.ItemOferta, error) {
	query := "SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta WHERE id_produto = $1 AND id_oferta = $2"
	row := s.db.QueryRowContext(ctx, query, id_produto, id_oferta)

	var io model.ItemOferta
	err := row.Scan(&io.Quantidade, &io.IDProduto, &io.IDOferta)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &io, nil
}

// GetAllByItemID busca todas as entradas de ItemOferta para um determinado produto.
func (s *Store) GetAllByItemID(ctx context.Context, id_produto int64) ([]model.ItemOferta, error) {
	query := "SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta WHERE id_produto = $1"
	rows, err := s.db.QueryContext(ctx, query, id_produto)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itensOferta := make([]model.ItemOferta, 0)
	for rows.Next() {
		var io model.ItemOferta
		err = rows.Scan(&io.Quantidade, &io.IDProduto, &io.IDOferta)
		if err != nil {
			return nil, err
		}
		itensOferta = append(itensOferta, io)
	}
	return itensOferta, nil
}

// GetAllByOfertaID busca todas as entradas de ItemOferta para uma determinada oferta.
func (s *Store) GetAllByOfertaID(ctx context.Context, id_oferta int64) ([]model.ItemOferta, error) {
	query := "SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta WHERE id_oferta = $1"
	rows, err := s.db.QueryContext(ctx, query, id_oferta)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ofertas := make([]model.ItemOferta, 0)
	for rows.Next() {
		var io model.ItemOferta
		err = rows.Scan(&io.Quantidade, &io.IDProduto, &io.IDOferta)
		if err != nil {
			return nil, err
		}
		ofertas = append(ofertas, io)
	}
	return ofertas, nil
}

func (s *Store) Create(ctx context.Context, props *model.ItemOferta) error {
	query := "INSERT INTO contem_item_oferta (quantidade, id_produto, id_oferta) VALUES ($1, $2, $3);"
	_, err := s.db.ExecContext(ctx, query, props.Quantidade, props.IDProduto, props.IDOferta)
	return err
}

func (s *Store) Update(ctx context.Context, props *model.ItemOferta) error {
	query := "UPDATE contem_item_oferta SET quantidade = $1 WHERE id_produto = $2 AND id_oferta = $3"
	res, err := s.db.ExecContext(ctx, query, props.Quantidade, props.IDProduto, props.IDOferta)
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

func (s *Store) Delete(ctx context.Context, id_produto int64, id_oferta int64) (*model.ItemOferta, error) {
	item, err := s.GetByComposedID(ctx, id_produto, id_oferta)
	if err != nil {
		return nil, err
	}

	query := "DELETE FROM contem_item_oferta WHERE id_produto = $1 AND id_oferta = $2"
	res, err := s.db.ExecContext(ctx, query, id_produto, id_oferta)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, types.ErrNotFound
	}

	return item, nil
}
