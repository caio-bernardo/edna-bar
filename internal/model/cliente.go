package model

import (
	"time"
)

type Cliente struct {
	Id             int64      `json:"id"`
	Nome           string     `json:"nome"`
	CPF            *string    `json:"cpf"`
	DataNascimento *time.Time `json:"data_nascimento"`
}

type ClienteCreate struct {
	Nome           string     `json:"nome"`
	CPF            *string    `json:"cpf"`
	DataNascimento *time.Time `json:"data_nascimento"` // Espera-se "YYYY-MM-DD" ou formato RFC3339
}

func (cc ClienteCreate) ToCliente() Cliente {
	return Cliente{
		Nome:           cc.Nome,
		CPF:            cc.CPF,
		DataNascimento: cc.DataNascimento,
	}
}
