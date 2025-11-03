package produto

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
	return &Store{
		db: db,
	}
}

func (s *Store) GetAllComercial(ctx context.Context) ([]model.Comercial, error) {
	query := "SELECT id_produto, nome, categoria, marca, quantidade_disponivel, quantidade_total, preco_venda FROM ProdutoComercial c INNER JOIN Produto p ON c.id_produto = p.id_produto;"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []model.Comercial
	for rows.Next() {
		c := model.Comercial{}
		err = rows.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca, &c.QntDisponivel, &c.QntTotal, &c.PrecoVenda)
		if err != nil {
			return nil, err
		}
		produtos = append(produtos, c)
	}

	return produtos, nil
}
func (s *Store) GetAllEstrutural(ctx context.Context) ([]model.Estrutural, error) {
	query := "SELECT id_produto, nome, categoria, marca, quantidade_disponivel, quantidade_total, preco_venda FROM ProdutoEstrutural e INNER JOIN Produto p ON e.id_produto = p.id_produto;"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []model.Estrutural
	for rows.Next() {
		c := model.Estrutural{}
		err = rows.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca, &c.QntDisponivel, &c.QntTotal)
		if err != nil {
			return nil, err
		}
		produtos = append(produtos, c)
	}

	return produtos, nil
}
func (s *Store) CreateComercial(ctx context.Context, props *model.Comercial) error {
	// Executa uma transação (desfaz insercao em caso de erro) pois precisamos fazer 2 insercoes
	t, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()

	query1 := "INSERT INTO Produto (nome, categoria, marca, quantidade_disponivel, quantidade_total) VALUES ($1, $2, $3, $4, $5) RETURNING id_produto;"
	query2 := "INSERT INTO ProdutoComercial (id_produto, preco_venda) VALUES ($1, $2);"

	r1, err := t.ExecContext(ctx, query1, props.Nome, props.Categoria, props.Marca, props.QntDisponivel, props.QntTotal)
	if err != nil {
		return err
	}

	idProduto, err := r1.LastInsertId()
	if err != nil {
		return err
	}

	_, err = t.ExecContext(ctx, query2, idProduto, props.PrecoVenda)
	if err != nil {
		return err
	}

	err = t.Commit()
	if err != nil {
		return err
	}
	// Atualiza o id do produto
	props.Id = idProduto
	return nil
}

func (s *Store) CreateEstrutural(ctx context.Context, props *model.Estrutural) error {
	// Executa uma transação (desfaz insercao em caso de erro) pois precisamos fazer 2 insercoes
	t, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()

	query1 := "INSERT INTO Produto (nome, categoria, marca, quantidade_disponivel, quantidade_total) VALUES ($1, $2, $3, $4, $5) RETURNING id_produto;"
	query2 := "INSERT INTO ProdutoEstrutural (id_produto) VALUES ($1);"

	r1, err := t.ExecContext(ctx, query1, props.Nome, props.Categoria, props.Marca, props.QntDisponivel, props.QntTotal)
	if err != nil {
		return err
	}

	idProduto, err := r1.LastInsertId()
	if err != nil {
		return err
	}

	_, err = t.ExecContext(ctx, query2, idProduto)
	if err != nil {
		return err
	}

	err = t.Commit()
	if err != nil {
		return err
	}

	// Atualiza o id do produto
	props.Id = idProduto
	return nil
}
func (s *Store) UpdateComercial(ctx context.Context, id int64, props *model.Comercial) error {
	t, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()
	query := "UPDATE Produto SET nome = $1, categoria = $2, marca = $3, quantidade_disponivel = $4, quantidade_total = $5 WHERE id_produto = $6;"
	res, err := t.ExecContext(ctx, query, props.Nome, props.Categoria, props.Marca, props.QntDisponivel, props.QntTotal, id)
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

	query2 := "UPDATE Produto SET preco_venda = $1 WHERE id_produto = $2;"
	res, err = t.ExecContext(ctx, query2, props.PrecoVenda, id)
	if err != nil {
		return err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return types.ErrNotFound
	}

	t.Commit()
	return nil
}

func (s *Store) UpdateEstrutural(ctx context.Context, id int64, props *model.Estrutural) error {
	query := "UPDATE Produto SET nome = $1, categoria = $2, marca = $3, quantidade_disponivel = $4, quantidade_total = $5 WHERE id_produto = $6;"

	res, err := s.db.ExecContext(ctx, query, props.Nome, props.Categoria, props.Marca, props.QntDisponivel, props.QntTotal, id)
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

func (s *Store) GetComercialByID(ctx context.Context, id int64) (*model.Comercial, error) {
	query := "SELECT id_produto, nome, categoria, marca, quantidade_disponivel, quantidade_total, preco_venda FROM ProdutoComercial c INNER JOIN Produto p ON c.id_produto = p.id_produto WHERE c.id_produto = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	c := model.Comercial{}
	err := row.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca, &c.QntDisponivel, &c.QntTotal, &c.PrecoVenda)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) GetEstruturalByID(ctx context.Context, id int64) (*model.Estrutural, error) {
	query := "SELECT id_produto, nome, categoria, marca, quantidade_disponivel, quantidade_total FROM ProdutoEstrutural e INNER JOIN Produto p ON e.id_produto = p.id_produto WHERE e.id_produto = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	c := model.Estrutural{}
	err := row.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca, &c.QntDisponivel, &c.QntTotal)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	// Derivadas do produto serão apagadas automaticamente por conta do condicional CASCADE
	query := "DELETE FROM Produto WHERE id_produto = $1"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
