DROP TYPE IF EXISTS tipo_de_produto;
CREATE TYPE tipo_de_produto AS ENUM ('estrutural', 'comercial');

CREATE TABLE IF NOT EXISTS Produto (
    id_produto serial PRIMARY KEY,
    nome varchar(50) NOT NULL,
    categoria text,
    marca text
);

CREATE TABLE IF NOT EXISTS ProdutoComercial (
    id_produto int NOT NULL,
    preco_venda decimal(6, 2) CHECK (preco_venda > 0) NOT NULL,

    CONSTRAINT fk_produto
        FOREIGN KEY(id_produto) 
        REFERENCES Produto(id_produto)
        ON DELETE CASCADE
);

ALTER TABLE ProdutoComercial
ADD PRIMARY KEY (id_produto);
