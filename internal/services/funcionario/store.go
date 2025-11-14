package funcionario

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

func (s *Store) GetAll(ctx context.Context, filter util.Filter) ([]model.Funcionario, error) {

	query := "SELECT id_funcionario, nome, CPF, tipo, expediente, salario, data_contratacao FROM Funcionario AS fc"
	rows, err := util.QueryRowsWithFilter(s.db, ctx, query, &filter, "fc")
	if err != nil {
		return nil, err
	}

	funcionarios := make([]model.Funcionario, 0)

	for rows.Next() {
		var funcionario model.Funcionario
		err = rows.Scan(&funcionario.Id, &funcionario.Nome, &funcionario.CPF, &funcionario.Tipo, &funcionario.Expediente, &funcionario.Salario, &funcionario.DataContratacao)
		if err != nil {
			return nil, err
		}
		funcionarios = append(funcionarios, funcionario)
	}

	return funcionarios, nil
}

func (s *Store) Create(ctx context.Context, props *model.Funcionario) error {
	query := "INSERT INTO Funcionario (nome, CPF, tipo, expediente, salario, data_contratacao) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_funcionario"
	res := s.db.QueryRowContext(ctx, query, props.Nome, props.CPF, props.Tipo, props.Expediente, props.Salario, props.DataContratacao)
	return res.Scan(&props.Id)
}

func (s *Store) GetByID(ctx context.Context, id int64) (*model.Funcionario, error) {
	query := "SELECT id_funcionario, nome, CPF, tipo, expediente, salario, data_contratacao FROM Funcionario WHERE id_funcionario = $1;"

	row := s.db.QueryRowContext(ctx, query, id)

	var funcionario model.Funcionario
	err := row.Scan(&funcionario.Id, &funcionario.Nome, &funcionario.CPF, &funcionario.Tipo, &funcionario.Expediente, &funcionario.Salario, &funcionario.DataContratacao)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &funcionario, nil
}

func (s *Store) Update(ctx context.Context, props *model.Funcionario) error {
	query := "UPDATE Funcionario SET nome = $1, CPF = $2, tipo = $3, expediente = $4, salario = $5, data_contratacao = $6 WHERE id_funcionario = $7;"

	res, err := s.db.ExecContext(ctx, query, props.Nome, props.CPF, props.Tipo, props.Expediente, props.Salario, props.DataContratacao, props.Id)
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

func (s *Store) Delete(ctx context.Context, id int64) (*model.Funcionario, error) {
	query := "DELETE FROM Funcionario WHERE id_funcionario = $1 RETURNING id_funcionario, nome, CPF, tipo, expediente, salario, data_contratacao;"

	var model model.Funcionario
	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&model.Id, &model.Nome, &model.CPF, &model.Tipo, &model.Expediente, &model.Salario, &model.DataContratacao)
	if err != nil {
		return nil, err
	}
	return &model, nil
}
