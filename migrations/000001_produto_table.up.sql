DROP TYPE IF EXISTS tipo_de_produto;
CREATE TYPE tipo_de_produto AS ENUM ('estrutural', 'comercial');

CREATE TABLE IF NOT EXISTS Produto (
    id_produto serial PRIMARY KEY,
    nome varchar(50) NOT NULL,
    categoria text,
    marca text,
    quantidade_disponivel int NOT NULL,
    quantidade_total int NOT NULL
);

CREATE TABLE IF NOT EXISTS ProdutoComercial (
    id_produto int PRIMARY KEY,
    preco_venda decimal(6, 2) NOT NULL,

    FOREIGN KEY (id_produto) REFERENCES Produto(id_produto) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ProdutoEstrutural (
    id_produto int PRIMARY KEY,
    FOREIGN KEY (id_produto) REFERENCES Produto(id_produto) ON DELETE CASCADE
);
