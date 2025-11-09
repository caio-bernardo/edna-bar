package model

import (
	"time"
)

type Lote struct {
	Id                int64      `json:"id_lote"`
	IdFornecedor      int64      `json:"id_fornecedor"`
	IdProduto         int64      `json:"id_produto"`
	DataFornecimento  time.Time  `json:"data_fornecimento"`
	Validade          *time.Time `json:"validade"`
	PrecoUnitario     float64    `json:"preco_unitario"`
	Estragados        *int       `json:"estragados"`
	QuantidadeInicial *int       `json:"quantidade_inicial"`
}

type LoteCreate struct {
	IdFornecedor      int64      `json:"id_fornecedor"`
	IdProduto         int64      `json:"id_produto"`
	DataFornecimento  time.Time  `json:"data_fornecimento"`
	Validade          *time.Time `json:"validade"`
	PrecoUnitario     float64    `json:"preco_unitario"`
	Estragados        *int       `json:"estragados"`
	QuantidadeInicial *int       `json:"quantidade_inicial"`
}

func (lc LoteCreate) ToLote() Lote {
	return Lote{
		IdFornecedor:      lc.IdFornecedor,
		IdProduto:         lc.IdProduto,
		DataFornecimento:  lc.DataFornecimento,
		Validade:          lc.Validade,
		PrecoUnitario:     lc.PrecoUnitario,
		Estragados:        lc.Estragados,
		QuantidadeInicial: lc.QuantidadeInicial,
	}
}
