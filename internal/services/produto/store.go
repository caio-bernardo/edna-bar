package produto

import (
	"context"
	"database/sql"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
	"log"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetAll(ctx context.Context, filter *util.Filter) ([]model.UnionProduto, error) {
	query := "SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda FROM Produto p LEFT JOIN ProdutoComercial AS c using (id_produto)"
	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, filter, "p")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	produtos := make([]model.UnionProduto, 0)
	for rows.Next() {
		c := model.UnionProduto{}
		err = rows.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca, &c.PrecoVenda)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, types.ErrNotFound
			}
			log.Printf("Error scanning row: %v", err)
			return nil, types.ErrInternalServer
		}
		produtos = append(produtos, c)
	}

	return produtos, nil
}

func (s *Store) GetAllComercial(ctx context.Context, filter *util.Filter) ([]model.Comercial, error) {
	query := `
		SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda
		FROM Produto p
		INNER JOIN ProdutoComercial c ON p.id_produto = c.id_produto`
	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, filter, "p")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	produtos := make([]model.Comercial, 0)
	for rows.Next() {
		c := model.Comercial{}
		err = rows.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca, &c.PrecoVenda)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, types.ErrNotFound
			}
			log.Printf("Error scanning row: %v", err)
			return nil, types.ErrInternalServer
		}
		produtos = append(produtos, c)
	}

	return produtos, nil
}

func (s *Store) GetAllEstrutural(ctx context.Context, filter *util.Filter) ([]model.Produto, error) {
	query := `
		SELECT p.id_produto, p.nome, p.categoria, p.marca
		FROM Produto p
		LEFT JOIN ProdutoComercial c ON p.id_produto = c.id_produto
		WHERE c.id_produto IS NULL`

	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, filter, "p")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	produtos := make([]model.Produto, 0)
	for rows.Next() {
		c := model.Produto{}
		err = rows.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, types.ErrNotFound
			}
			log.Printf("Error scanning row: %v", err)
			return nil, types.ErrInternalServer
		}
		produtos = append(produtos, c)
	}

	return produtos, nil
}

func (s *Store) CreateComercial(ctx context.Context, props *model.Comercial) error {
	// Inicia a transação
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insere na tabela Produto
	queryProduto := "INSERT INTO Produto (nome, categoria, marca) VALUES ($1, $2, $3) RETURNING id_produto;"
	row := tx.QueryRowContext(ctx, queryProduto, props.Nome, props.Categoria, props.Marca)
	err = row.Scan(&props.Id)
	if err != nil {
		return err
	}

	// Insere na tabela ProdutoComercial
	queryComercial := "INSERT INTO ProdutoComercial (id_produto, preco_venda) VALUES ($1, $2);"
	_, err = tx.ExecContext(ctx, queryComercial, props.Id, props.PrecoVenda)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) Create(ctx context.Context, props *model.Produto) error {
	query := "INSERT INTO Produto (nome, categoria, marca) VALUES ($1, $2, $3) RETURNING id_produto;"

	row := s.db.QueryRowContext(ctx, query, props.Nome, props.Categoria, props.Marca)
	err := row.Scan(&props.Id)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) UpdateComercial(ctx context.Context, props *model.Comercial) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Atualiza a tabela Produto
	queryProduto := "UPDATE Produto SET nome = $1, categoria = $2, marca = $3 WHERE id_produto = $4;"
	res, err := tx.ExecContext(ctx, queryProduto, props.Nome, props.Categoria, props.Marca, props.Id)
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

	// Atualiza a tabela ProdutoComercial
	queryComercial := "UPDATE ProdutoComercial SET preco_venda = $1 WHERE id_produto = $2;"
	_, err = tx.ExecContext(ctx, queryComercial, props.PrecoVenda, props.Id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) Update(ctx context.Context, props *model.Produto) error {
	query := "UPDATE Produto SET nome = $1, categoria = $2, marca = $3 WHERE id_produto = $4;"

	res, err := s.db.ExecContext(ctx, query, props.Nome, props.Categoria, props.Marca, props.Id)
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
	query := `
		SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda
		FROM Produto p
		INNER JOIN ProdutoComercial c ON p.id_produto = c.id_produto
		WHERE p.id_produto = $1`

	row := s.db.QueryRowContext(ctx, query, id)
	c := model.Comercial{}
	err := row.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca, &c.PrecoVenda)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, types.ErrNotFound
		case sql.ErrConnDone:
			log.Println("[WARN] Connection closed.")
			return nil, types.ErrInternalServer
		default:
			log.Println("[ERROR] Unexpected error:", err)
			return nil, err
		}
	}
	return &c, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.Produto, error) {
	query := "SELECT id_produto, nome, categoria, marca FROM Produto WHERE id_produto = $1"
	row := s.db.QueryRowContext(ctx, query, id)
	c := model.Produto{}
	err := row.Scan(&c.Id, &c.Nome, &c.Categoria, &c.Marca)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) GetQntByID(ctx context.Context, id int64) (*model.ProdutoWithQnt, error) {

	// Quantidade de produtos disponiveis
	// Resultado = inicias - estrados - vendidos
	// coalesce converte o null da soma em zero.
	// Assim possiveis valores nulos resultam em zero
	query := `
	SELECT p.id_produto, p.nome, p.categoria, p.marca,
		COALESCE(SUM(quantidade_inicial) - SUM(estragados), 0) - COALESCE(SUM(quantidade), 0) AS quantidade_disponivel
		FROM Produto p
		LEFT JOIN lote USING (id_produto)
	 	LEFT JOIN item_venda USING (id_lote)
		WHERE p.id_produto = $1
		GROUP BY p.id_produto;`

	row := s.db.QueryRowContext(ctx, query, id)

	var model model.ProdutoWithQnt
	err := row.Scan(&model.Id, &model.Nome, &model.Categoria, &model.Marca, &model.Qnt)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	// Derivadas do produto serão apagadas automaticamente por conta da herança
	query := "DELETE FROM Produto WHERE id_produto = $1"
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
