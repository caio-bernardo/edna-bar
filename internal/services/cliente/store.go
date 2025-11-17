package cliente

import (
	"context"
	"database/sql"
	"edna/internal/model"
	"edna/internal/types"
	"edna/internal/util"
	"fmt"
	"strconv"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetAll(ctx context.Context, filter util.Filter) ([]model.Cliente, error) {
	query := "SELECT id_cliente, nome, cpf, data_nascimento FROM Cliente AS c"

	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, &filter, "c")
	if err != nil {
		return nil, err
	}

	clientes := make([]model.Cliente, 0)
	for rows.Next() {
		var c model.Cliente
		err = rows.Scan(&c.Id, &c.Nome, &c.CPF, &c.DataNascimento)
		if err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}
	return clientes, nil
}

func (s *Store) GetAllWithSaldo(ctx context.Context, filter util.Filter) ([]model.ClienteWithSaldo, error) {
	// Criamos uma lista de ids de clientes que estão devendo dinheiro
	// Juntamos com clientes e substituimos por zero valores nulos.
	query := `
	WITH ClienteDevedor AS (
		SELECT id_cliente, COALESCE(SUM(quantidade * valor_unitario), 0)::numeric(12, 2) as saldo_devedor
		FROM Venda
		LEFT JOIN item_venda USING(id_venda)
	 	WHERE data_hora_pagamento IS NULL
		GROUP BY id_cliente
	) SELECT id_cliente, nome, cpf, data_nascimento,
		COALESCE(saldo_devedor, 0)::numeric(12, 2)
		FROM Cliente
		LEFT JOIN ClienteDevedor USING(id_cliente);
	`
	// Constroi a query de filtros manualmente
	var values []any
	i := 0
	for k, v := range filter.Filters {
		// reescrever saldo_devedor como zero caso seja nulo
		if k == "saldo_devedor" {
			k = "COALESCE(saldo_devedor, 0)::numeric(12, 2)"
		}
		if i == 0 {
			query += " WHERE"
		} else {
			query += " AND"
		}
		switch v.Operator {
		case "lt":
			values = append(values, v.Value)
			query += fmt.Sprintf(" %s < $%d", k, len(values))
		case "gt":
			values = append(values, v.Value)
			query += fmt.Sprintf(" %s > $%d", k, len(values))
		case "eq":
			values = append(values, v.Value)
			query += fmt.Sprintf(" %s = $%d", k, len(values))
		case "le":
			values = append(values, v.Value)
			query += fmt.Sprintf(" %s <= $%d", k, len(values))
		case "ge":
			values = append(values, v.Value)
			query += fmt.Sprintf(" %s >= $%d", k, len(values))
		case "ne":
			values = append(values, v.Value)
			query += fmt.Sprintf(" %s != $%d", k, len(values))
		case "like":
			values = append(values, v.Value)
			query += fmt.Sprintf(" %s LIKE '%%' || $%d || '%%'", k, len(values))
		case "ilike":
			values = append(values, v.Value)
			query += fmt.Sprintf(" %s ILIKE '%%' || $%d || '%%'", k, len(values))
		default:
		}
		i += 1
	}

	// ordenação
	for i, v := range filter.Sorts {
		if i == 0 {
			query += " ORDER BY"
		} else {
			query += ","
		}

		str, fminus := strings.CutPrefix(v, "-")
		query += " " + str
		if fminus {
			query += " DESC"
		}
	}

	// paginação
	if filter.Offset > 0 {
		values = append(values, filter.Offset)
		query += " OFFSET $" + strconv.Itoa(len(values))
	}
	if filter.Limit > 0 {
		values = append(values, filter.Limit)
		query += " LIMIT $" + strconv.Itoa(len(values))
	}

	rows, err := s.db.QueryContext(ctx, query, values)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clientes := make([]model.ClienteWithSaldo, 0)
	for rows.Next() {
		var c model.ClienteWithSaldo
		err = rows.Scan(&c.Id, &c.Nome, &c.CPF, &c.DataNascimento, &c.SaldoDevedor)
		if err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}
	return clientes, nil
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.Cliente, error) {
	query := "SELECT id_cliente, nome, cpf, data_nascimento FROM Cliente WHERE id_cliente = $1;"
	row := s.db.QueryRowContext(ctx, query, id)

	var c model.Cliente
	err := row.Scan(&c.Id, &c.Nome, &c.CPF, &c.DataNascimento)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (s *Store) GetByIDWithSaldo(ctx context.Context, id int64) (*model.ClienteWithSaldo, error) {
	query := `
	WITH ClienteDevedor AS (
		SELECT id_cliente, COALESCE(SUM(quantidade * valor_unitario), 0)::numeric(12, 2) as saldo_devedor
		FROM Venda
		LEFT JOIN item_venda USING(id_venda)
	 	WHERE data_hora_pagamento IS NULL
		GROUP BY id_cliente
	) SELECT id_cliente, nome, cpf, data_nascimento,
		COALESCE(saldo_devedor, 0)::numeric(12, 2)
		FROM Cliente
	 	LEFT JOIN ClienteDevedor USING(id_cliente)
		WHERE id_cliente = $1;
	`
	row := s.db.QueryRowContext(ctx, query, id)

	var c model.ClienteWithSaldo
	err := row.Scan(&c.Id, &c.Nome, &c.CPF, &c.DataNascimento, &c.SaldoDevedor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (s *Store) Create(ctx context.Context, props *model.Cliente) error {
	query := "INSERT INTO Cliente (nome, cpf, data_nascimento) VALUES ($1, $2, $3) RETURNING id_cliente;"
	res := s.db.QueryRowContext(ctx, query, props.Nome, props.CPF, props.DataNascimento)
	return res.Scan(&props.Id)
}

func (s *Store) Update(ctx context.Context, props *model.Cliente) error {
	query := "UPDATE Cliente SET nome = $1, cpf = $2, data_nascimento = $3 WHERE id_cliente = $4;"
	res, err := s.db.ExecContext(ctx, query, props.Nome, props.CPF, props.DataNascimento, props.Id)
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

func (s *Store) Delete(ctx context.Context, id int64) (*model.Cliente, error) {
	query := "DELETE FROM Cliente WHERE id_cliente = $1 RETURNING id_cliente, nome, cpf, data_nascimento;"
	var m model.Cliente
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&m.Id, &m.Nome, &m.CPF, &m.DataNascimento)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}
