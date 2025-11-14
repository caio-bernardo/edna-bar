package model

type Funcionario struct {
	Id              int64   `json:"id"`
	Nome            string  `json:"nome"`
	CPF             string  `json:"CPF"`
	Tipo            string  `json:"tipo"`
	Expediente      string  `json:"expediente"`
	Salario         float64 `json:"salario`
	DataContratacao string  `json:"data_contratacao"`
}

type FuncionarioCreate struct {
	Nome            string  `json:"nome"`
	CPF             string  `json:"CPF"`
	Tipo            string  `json:"tipo"`
	Expediente      string  `json:"expediente"`
	Salario         float64 `json:"salario`
	DataContratacao string  `json:"data_contratacao"`
}

func (fc FuncionarioCreate) ToFuncionario() Funcionario {
	return Funcionario{
		Nome:            fc.CPF,
		CPF:             fc.CPF,
		Tipo:            fc.Tipo,
		Expediente:      fc.Expediente,
		Salario:         fc.Salario,
		DataContratacao: fc.DataContratacao,
	}
}
