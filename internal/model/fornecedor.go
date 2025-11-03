package model

type Fornecedor struct {
	Id int64 `json:"id"`
	Nome string `json:"nome"`
	CNPJ string `json:"cpnj"`
}

type FornecedorPayload struct {
	Nome string `json:"nome"`
	CNPJ string `json:"cpnj"`
}

func FromPayload(payload FornecedorPayload) Fornecedor {
	return Fornecedor{
		Nome: payload.Nome,
		CNPJ: payload.CNPJ,
	}
}
