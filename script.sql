-- ================================
-- SCRIPT SQL COMPLETO - EDNA BAR
-- ================================
-- Este arquivo contém todas as queries SQL das migrations e services
-- incluindo variações com filtros aplicados conforme definido em filter.go

-- ================================
-- MIGRATIONS - CRIAÇÃO DE ESTRUTURA
-- ================================

-- Migration 000001: Criação da tabela Fornecedor
CREATE TABLE IF NOT EXISTS Fornecedor (
    CNPJ char(14),
    nome varchar(50) NOT NULL,
    id_fornecedor serial PRIMARY KEY
);

-- Migration 000002: Criação das tabelas Produto
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

ALTER TABLE ProdutoComercial ADD PRIMARY KEY (id_produto);

-- Migration 000003: Criação da tabela Lote
CREATE TABLE IF NOT EXISTS Lote (
    id_lote serial PRIMARY KEY,
    id_fornecedor int NOT NULL,
    id_produto int NOT NULL,

    data_fornecimento date NOT NULL,
    validade date,
    preco_unitario decimal(6, 2) NOT NULL CHECK (preco_unitario > 0),
    estragados int CHECK (estragados >= 0) DEFAULT 0,
    quantidade_inicial int CHECK (quantidade_inicial > 0),

    FOREIGN KEY (id_fornecedor) REFERENCES Fornecedor(id_fornecedor),
    FOREIGN KEY (id_produto) REFERENCES Produto(id_produto)
);

-- Migration 000004: Criação da tabela Cliente
CREATE TABLE IF NOT EXISTS Cliente (
    id_cliente serial PRIMARY KEY,
    CPF char(11),
    data_nascimento date,
    nome varchar(50) NOT NULL
);

-- Migration 000005: Criação da tabela Funcionario
DROP TYPE IF EXISTS tipo_de_funcionario;
CREATE TYPE tipo_de_funcionario AS ENUM ('garcom', 'seguranca',
'caixa', 'faxineiro', 'balconista');

DROP TYPE IF EXISTS tipo_de_expediente;
CREATE TYPE tipo_de_expediente AS ENUM ('manha', 'tarde', 'noite', 'madrugada');

CREATE TABLE IF NOT EXISTS Funcionario (
    id_funcionario serial PRIMARY KEY,
    nome varchar(50) NOT NULL,
    CPF char(11) NOT NULL,
    tipo tipo_de_funcionario NOT NULL,
    expediente tipo_de_expediente NOT NULL,
    data_contratacao date NOT NULL,
    salario decimal(8, 2) NOT NULL
);

-- Migration 000006: Criação da tabela Venda
DROP TYPE IF EXISTS tipo_de_pagamento;
CREATE TYPE tipo_de_pagamento AS ENUM ('credito', 'debito', 'pix', 'dinheiro', 'VA/VR');

CREATE TABLE IF NOT EXISTS Venda (
    id_venda serial PRIMARY KEY,
    data_hora_venda timestamp NOT NULL DEFAULT now(),
    data_hora_pagamento timestamp,
    tipo_pagamento tipo_de_pagamento,

    id_cliente int NOT NULL,
    id_funcionario int NOT NULL,

    FOREIGN KEY (id_cliente) REFERENCES Cliente(id_cliente) ON DELETE RESTRICT,
    FOREIGN KEY (id_funcionario) REFERENCES Funcionario(id_funcionario) ON DELETE SET NULL
);

-- Migration 000007: Criação da tabela Oferta
CREATE TABLE IF NOT EXISTS Oferta (
    id_oferta serial PRIMARY KEY,
    nome varchar(50) NOT NULL,
    data_criacao date NOT NULL DEFAULT CURRENT_DATE,
    data_inicio date,
    data_fim date,
    valor_fixo decimal(6, 2),
    percentual_desconto int
);

-- Migration 000008: Criação da tabela contem_item_oferta
CREATE TABLE IF NOT EXISTS contem_item_oferta (
    quantidade int NOT NULL CHECK (quantidade > 0),
    id_produto int,
    id_oferta int,
    PRIMARY KEY (id_produto, id_oferta),

    FOREIGN KEY (id_produto) REFERENCES ProdutoComercial(id_produto),
    FOREIGN KEY (id_oferta) REFERENCES Oferta(id_oferta) ON DELETE CASCADE
);

-- Migration 000011: Criação da tabela item_venda
CREATE TABLE IF NOT EXISTS item_venda (
    id_item_venda SERIAL PRIMARY KEY,
    id_venda int NOT NULL REFERENCES Venda(id_venda) ON DELETE CASCADE,
    id_lote int REFERENCES Lote(id_lote) ON DELETE RESTRICT,
    quantidade int CHECK (quantidade > 0),
    valor_unitario decimal(6, 2) NOT NULL
);

-- Migration 000012: Criação da tabela aplica_oferta
CREATE TABLE IF NOT EXISTS aplica_oferta (
    id_aplica_oferta SERIAL PRIMARY KEY,
    id_oferta int REFERENCES Oferta(id_oferta) ON DELETE CASCADE,
    id_venda int REFERENCES Venda(id_venda) ON DELETE CASCADE,
    id_item_venda int,

    FOREIGN KEY (id_item_venda) REFERENCES item_venda(id_item_venda)
);

-- ================================
-- DADOS INICIAIS (Migration 000010)
-- ================================

-- Insere Fornecedores
INSERT INTO Fornecedor (CNPJ, nome) VALUES
('11223344000155', 'Distribuidora de Bebidas Central'),
('22334455000166', 'Frios & Carnes Mart'),
('33445566000177', 'Hortifruti Verdão');

-- Insere Clientes
INSERT INTO Cliente (nome, CPF, data_nascimento) VALUES
('Serginho Groisman', '11122233344', '1990-05-15'),
('Neymar Jr.', '22233344455', '1985-11-01'),
('Martinho da Vila', '33344455566', '2001-02-20');

-- Insere Funcionarios
INSERT INTO Funcionario (nome, CPF, tipo, expediente, data_contratacao, salario) VALUES
('Joelma Calypso', '99988877766', 'garcom', 'noite', '2023-01-10', 2200.00),
('Derick Green Sepultura', '88877766655', 'caixa', 'noite', '2022-11-05', 2350.00);

-- Insere Produtos Estruturais
INSERT INTO Produto (nome, categoria, marca) VALUES
('Guardanapo de Papel', 'Descartáveis', 'Tio Mané'),
('Copo Descartável 300ml', 'Descartáveis', 'PlastCop');

-- Insere Produtos Comerciais
INSERT INTO Produto (nome, categoria, marca) VALUES
('Cerveja 600ml', 'Bebidas', 'Skoll'),
('Porção de Fritas Grande', 'Porções', 'Casa'),
('Refrigerante Lata 350ml', 'Bebidas', 'Guarana Mineiro');

INSERT INTO ProdutoComercial (id_produto, preco_venda) VALUES
(
    (SELECT id_produto FROM Produto WHERE nome = 'Cerveja 600ml'),
    12.50
),
(
    (SELECT id_produto FROM Produto WHERE nome = 'Porção de Fritas Grande'),
    25.00
),
(
    (SELECT id_produto FROM Produto WHERE nome = 'Refrigerante Lata 350ml'),
    6.00
);

-- Insere Lotes
INSERT INTO Lote (id_fornecedor, id_produto, data_fornecimento, preco_unitario, quantidade_inicial, validade) VALUES
(
    (SELECT id_fornecedor FROM Fornecedor WHERE nome = 'Distribuidora de Bebidas Central'),
    (SELECT id_produto FROM Produto WHERE nome = 'Cerveja 600ml'),
    '2025-10-01', 5.50, 100, '2026-04-01'
),
(
    (SELECT id_fornecedor FROM Fornecedor WHERE nome = 'Frios & Carnes Mart'),
    (SELECT id_produto FROM Produto WHERE nome = 'Porção de Fritas Grande'),
    '2025-11-01', 10.00, 50, '2025-11-15'
),
(
    (SELECT id_fornecedor FROM Fornecedor WHERE nome = 'Distribuidora de Bebidas Central'),
    (SELECT id_produto FROM Produto WHERE nome = 'Refrigerante Lata 350ml'),
    '2025-10-20', 2.50, 200, '2026-10-20'
),
(
    (SELECT id_fornecedor FROM Fornecedor WHERE nome = 'Hortifruti Verdão'),
    (SELECT id_produto FROM Produto WHERE nome = 'Guardanapo de Papel'),
    '2025-10-01', 0.50, 500, NULL
),
(
    (SELECT id_fornecedor FROM Fornecedor WHERE nome = 'Distribuidora de Bebidas Central'),
    (SELECT id_produto FROM Produto WHERE nome = 'Copo Descartável 300ml'),
    '2025-10-01', 0.50, 500, NULL
);

-- Insere Ofertas
INSERT INTO Oferta (nome, data_inicio, data_fim, percentual_desconto) VALUES
('Happy Hour Cerveja', '2025-11-01', '2025-11-30', 20);

INSERT INTO Oferta (nome, data_inicio, data_fim, valor_fixo) VALUES
('Combo Fritas + Refri', '2025-11-01', '2025-11-30', 28.00);

-- Insere Itens da Oferta
INSERT INTO contem_item_oferta (id_oferta, id_produto, quantidade) VALUES
(
    (SELECT id_oferta FROM Oferta WHERE nome = 'Happy Hour Cerveja'),
    (SELECT id_produto FROM Produto WHERE nome = 'Cerveja 600ml'),
    1
),
(
    (SELECT id_oferta FROM Oferta WHERE nome = 'Combo Fritas + Refri'),
    (SELECT id_produto FROM Produto WHERE nome = 'Porção de Fritas Grande'),
    1
),
(
    (SELECT id_oferta FROM Oferta WHERE nome = 'Combo Fritas + Refri'),
    (SELECT id_produto FROM Produto WHERE nome = 'Refrigerante Lata 350ml'),
    1
);

-- Insere Vendas
INSERT INTO Venda (id_cliente, id_funcionario, data_hora_pagamento, tipo_pagamento) VALUES
(
    (SELECT id_cliente FROM Cliente WHERE nome = 'Serginho Groisman'),
    (SELECT id_funcionario FROM Funcionario WHERE nome = 'Derick Green Sepultura'),
    NOW(),
    'pix'
);
INSERT INTO Venda (id_cliente, id_funcionario, data_hora_pagamento, tipo_pagamento) VALUES
(
    (SELECT id_cliente FROM Cliente WHERE nome = 'Neymar Jr.'),
    (SELECT id_funcionario FROM Funcionario WHERE nome = 'Joelma Calypso'),
    NOW(),
    'debito'
);

-- Mais dados (migration 13 & 14)
-- Insere itens na venda para o Serginho (Cerveja)
INSERT INTO item_venda (id_venda, id_lote, quantidade, valor_unitario)
VALUES (
    (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Serginho Groisman')),
    (SELECT id_lote FROM Lote WHERE id_produto = (SELECT id_produto FROM Produto WHERE nome = 'Cerveja 600ml') LIMIT 1),
    2, -- Quantidade
    12.50 -- Preço unitário original
);

-- Aplica a oferta "Happy Hour" na venda do Serginho
INSERT INTO aplica_oferta (id_venda, id_oferta, id_item_venda)
VALUES (
    (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Serginho Groisman')),
    (SELECT id_oferta FROM Oferta WHERE nome = 'Happy Hour Cerveja'),
    (SELECT id_item_venda FROM item_venda WHERE id_venda = (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Serginho Groisman')) LIMIT 1)
);

-- Insere itens na venda para o Neymar (Fritas)
INSERT INTO item_venda (id_venda, id_lote, quantidade, valor_unitario)
VALUES (
    (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Neymar Jr.')),
    (SELECT id_lote FROM Lote WHERE id_produto = (SELECT id_produto FROM Produto WHERE nome = 'Porção de Fritas Grande') LIMIT 1),
    1,
    25.00
);

-- Insere itens na venda para o Neymar (Refri)
INSERT INTO item_venda (id_venda, id_lote, quantidade, valor_unitario)
VALUES (
    (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Neymar Jr.')),
    (SELECT id_lote FROM Lote WHERE id_produto = (SELECT id_produto FROM Produto WHERE nome = 'Refrigerante Lata 350ml') LIMIT 1),
    1,
    6.00
);

-- Aplicando a oferta "Combo" nos itens do Neymar
-- Nota: Dependendo da sua lógica de negócio, a oferta pode ser ligada à venda geral ou a cada item.
-- Aqui estou ligando aos itens para consistência com a tabela 'aplica_oferta' que tem FK para 'item_venda'.

-- Liga oferta à Frita
INSERT INTO aplica_oferta (id_venda, id_oferta, id_item_venda)
VALUES (
    (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Neymar Jr.')),
    (SELECT id_oferta FROM Oferta WHERE nome = 'Combo Fritas + Refri'),
    (SELECT id_item_venda FROM item_venda
     WHERE id_venda = (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Neymar Jr.'))
     AND id_lote = (SELECT id_lote FROM Lote WHERE id_produto = (SELECT id_produto FROM Produto WHERE nome = 'Porção de Fritas Grande') LIMIT 1)
    )
);

-- Liga oferta ao Refri
INSERT INTO aplica_oferta (id_venda, id_oferta, id_item_venda)
VALUES (
    (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Neymar Jr.')),
    (SELECT id_oferta FROM Oferta WHERE nome = 'Combo Fritas + Refri'),
    (SELECT id_item_venda FROM item_venda
     WHERE id_venda = (SELECT id_venda FROM Venda WHERE id_cliente = (SELECT id_cliente FROM Cliente WHERE nome = 'Neymar Jr.'))
     AND id_lote = (SELECT id_lote FROM Lote WHERE id_produto = (SELECT id_produto FROM Produto WHERE nome = 'Refrigerante Lata 350ml') LIMIT 1)
    )
);
-- Popula a relação da oferta "Happy Hour Cerveja"
-- Define que esta oferta está associada ao produto "Cerveja 600ml"
INSERT INTO contem_item_oferta (id_oferta, id_produto, quantidade)
VALUES (
    (SELECT id_oferta FROM Oferta WHERE nome = 'Happy Hour Cerveja'),
    (SELECT id_produto FROM Produto WHERE nome = 'Cerveja 600ml'),
    3 -- Quantidade base para a regra da oferta
)
ON CONFLICT DO NOTHING;

-- Popula a relação da oferta "Combo Fritas + Refri" (Item 1: Fritas)
INSERT INTO contem_item_oferta (id_oferta, id_produto, quantidade)
VALUES (
    (SELECT id_oferta FROM Oferta WHERE nome = 'Combo Fritas + Refri'),
    (SELECT id_produto FROM Produto WHERE nome = 'Porção de Fritas Grande'),
    1
)
ON CONFLICT DO NOTHING;

-- Popula a relação da oferta "Combo Fritas + Refri" (Item 2: Refri)
INSERT INTO contem_item_oferta (id_oferta, id_produto, quantidade)
VALUES (
    (SELECT id_oferta FROM Oferta WHERE nome = 'Combo Fritas + Refri'),
    (SELECT id_produto FROM Produto WHERE nome = 'Refrigerante Lata 350ml'),
    1
)
ON CONFLICT DO NOTHING;

-- ================================
-- QUERIES DOS SERVICES
-- ================================

-- ================================
-- FORNECEDOR SERVICE
-- ================================

-- GetAll - Query básica
SELECT id_fornecedor, nome, CNPJ FROM Fornecedor AS f;

-- GetAll com filtros (campos: nome, cnpj)
-- Exemplo com filtros aplicados:
SELECT id_fornecedor, nome, CNPJ FROM Fornecedor AS f
 WHERE f.nome LIKE '%' || $1 || '%' AND f.CNPJ = $2
 ORDER BY nome DESC OFFSET $3 LIMIT $4;

-- Create
INSERT INTO Fornecedor (nome, CNPJ) VALUES ($1, $2) RETURNING id_fornecedor;

-- GetByID
SELECT id_fornecedor, nome, CNPJ FROM Fornecedor WHERE id_fornecedor = $1;

-- Update
UPDATE Fornecedor SET nome = $1, CNPJ = $2 WHERE id_fornecedor = $3;

-- Delete
DELETE FROM Fornecedor WHERE id_fornecedor = $1 RETURNING id_fornecedor, nome, CNPJ;

-- ================================
-- CLIENTE SERVICE
-- ================================

-- GetAll - Query básica
SELECT id_cliente, nome, cpf, data_nascimento FROM Cliente AS c;

-- GetAll com filtros (campos: nome, cpf, data_nascimento)
-- Exemplo com filtros aplicados:
SELECT id_cliente, nome, cpf, data_nascimento FROM Cliente AS c
 WHERE c.nome ILIKE '%' || $1 || '%' AND c.cpf = $2 AND c.data_nascimento >= $3
 ORDER BY nome DESC OFFSET $4 LIMIT $5;

-- GetAllWithSaldo - Query básica
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

-- GetAllWithSaldo com filtros (campos: nome, cpf, data_nascimento, saldo_devedor)
-- Exemplo com filtros aplicados:
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
 WHERE nome LIKE '%' || $1 || '%' AND COALESCE(saldo_devedor, 0)::numeric(12, 2) > $2
 ORDER BY nome, saldo_devedor DESC OFFSET $3 LIMIT $4;

-- GetByID
SELECT id_cliente, nome, cpf, data_nascimento FROM Cliente WHERE id_cliente = $1;

-- GetByIDWithSaldo
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

-- Create
INSERT INTO Cliente (nome, cpf, data_nascimento) VALUES ($1, $2, $3) RETURNING id_cliente;

-- Update
UPDATE Cliente SET nome = $1, cpf = $2, data_nascimento = $3 WHERE id_cliente = $4;

-- Delete
DELETE FROM Cliente WHERE id_cliente = $1 RETURNING id_cliente, nome, cpf, data_nascimento;

-- ================================
-- PRODUTO SERVICE
-- ================================

-- GetAll - Query básica
SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda FROM Produto p LEFT JOIN ProdutoComercial AS c using (id_produto);

-- GetAll com filtros (campos: nome, categoria, marca)
-- Exemplo com filtros aplicados:
SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda FROM Produto p LEFT JOIN ProdutoComercial AS c using (id_produto)
 WHERE p.nome LIKE '%' || $1 || '%' AND p.categoria = $2 AND p.marca ILIKE '%' || $3 || '%'
 ORDER BY nome DESC OFFSET $4 LIMIT $5;

-- GetAllComercial - Query básica
SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda
FROM Produto p
INNER JOIN ProdutoComercial c ON p.id_produto = c.id_produto;

-- GetAllComercial com filtros (campos: nome, categoria, marca, preco_venda)
-- Exemplo com filtros aplicados:
SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda
FROM Produto p
INNER JOIN ProdutoComercial c ON p.id_produto = c.id_produto
 WHERE p.nome LIKE '%' || $1 || '%' AND p.categoria = $2 AND c.preco_venda >= $3
 ORDER BY nome, preco_venda DESC OFFSET $4 LIMIT $5;

-- GetAllEstrutural - Query básica
SELECT p.id_produto, p.nome, p.categoria, p.marca
FROM Produto p
LEFT JOIN ProdutoComercial c ON p.id_produto = c.id_produto
WHERE c.id_produto IS NULL;

-- GetAllEstrutural com filtros (campos: nome, categoria, marca)
-- Exemplo com filtros aplicados:
SELECT p.id_produto, p.nome, p.categoria, p.marca
FROM Produto p
LEFT JOIN ProdutoComercial c ON p.id_produto = c.id_produto
WHERE c.id_produto IS NULL AND p.nome LIKE '%' || $1 || '%' AND p.categoria = $2
 ORDER BY nome DESC OFFSET $3 LIMIT $4;

-- CreateComercial (transação)
INSERT INTO Produto (nome, categoria, marca) VALUES ($1, $2, $3) RETURNING id_produto;
INSERT INTO ProdutoComercial (id_produto, preco_venda) VALUES ($1, $2);

-- Create
INSERT INTO Produto (nome, categoria, marca) VALUES ($1, $2, $3) RETURNING id_produto;

-- UpdateComercial (transação)
UPDATE Produto SET nome = $1, categoria = $2, marca = $3 WHERE id_produto = $4;
UPDATE ProdutoComercial SET preco_venda = $1 WHERE id_produto = $2;

-- Update
UPDATE Produto SET nome = $1, categoria = $2, marca = $3, quantidade_disponivel = $4, quantidade_total = $5 WHERE id_produto = $6;

-- GetComercialByID
SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda
FROM Produto p
INNER JOIN ProdutoComercial c ON p.id_produto = c.id_produto
WHERE p.id_produto = $1;

-- GetByID
SELECT id_produto, nome, categoria, marca FROM Produto WHERE id_produto = $1;

-- GetQntByID
SELECT p.id_produto, p.nome, p.categoria, p.marca,
    COALESCE(SUM(quantidade_inicial) - SUM(estragados), 0) - COALESCE(SUM(quantidade), 0) AS quantidade_disponivel
    FROM Produto p
    LEFT JOIN lote USING (id_produto)
    LEFT JOIN item_venda USING (id_lote)
    WHERE p.id_produto = $1
    GROUP BY p.id_produto;

-- Delete
DELETE FROM Produto WHERE id_produto = $1;

-- ================================
-- FUNCIONARIO SERVICE
-- ================================

-- GetAll - Query básica
SELECT id_funcionario, nome, CPF, tipo, expediente, salario, data_contratacao FROM Funcionario AS fc;

-- GetAll com filtros (campos: data_contratacao, salario, expediente, tipo, CPF, nome, id_funcionario)
-- Exemplo com filtros aplicados:
SELECT id_funcionario, nome, CPF, tipo, expediente, salario, data_contratacao FROM Funcionario AS fc
 WHERE fc.nome LIKE '%' || $1 || '%' AND fc.tipo = $2 AND fc.data_contratacao >= $3
 ORDER BY nome DESC OFFSET $4 LIMIT $5;

-- Create
INSERT INTO Funcionario (nome, CPF, tipo, expediente, salario, data_contratacao) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_funcionario;

-- GetByID
SELECT id_funcionario, nome, CPF, tipo, expediente, salario, data_contratacao FROM Funcionario WHERE id_funcionario = $1;

-- Update
UPDATE Funcionario SET nome = $1, CPF = $2, tipo = $3, expediente = $4, salario = $5, data_contratacao = $6 WHERE id_funcionario = $7;

-- Delete
DELETE FROM Funcionario WHERE id_funcionario = $1 RETURNING id_funcionario, nome, CPF, tipo, expediente, salario, data_contratacao;

-- ================================
-- LOTE SERVICE
-- ================================

-- GetAll - Query básica
SELECT id_lote, id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial FROM Lote AS l;

-- GetAll com filtros (campos: id_lote, id_fornecedor, id_produto, preco_unitario, estragados, quantidade_inicial, validade)
-- Exemplo com filtros aplicados:
SELECT id_lote, id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial FROM Lote AS l
 WHERE l.id_fornecedor = $1 AND l.preco_unitario >= $2 AND l.validade <= $3
 ORDER BY validade, preco_unitario DESC OFFSET $4 LIMIT $5;

-- GetByID
SELECT id_lote, id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial FROM Lote WHERE id_lote = $1;

-- GetAllByIDProduto
SELECT * FROM Lote WHERE id_produto = $1;

-- Create
INSERT INTO Lote (id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id_lote;

-- Update
UPDATE Lote SET
id_fornecedor = $1, id_produto = $2, data_fornecimento = $3, validade = $4,
preco_unitario = $5, estragados = $6, quantidade_inicial = $7
WHERE id_lote = $8;

-- Delete
DELETE FROM Lote WHERE id_lote = $1 RETURNING id_lote, id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial;

-- GetRelatorio
SELECT
    EXTRACT(YEAR FROM data_fornecimento)::int AS ano,
    EXTRACT(MONTH FROM data_fornecimento)::int AS mes,
    COALESCE(SUM(preco_unitario * quantidade_inicial), 0) AS total_gasto,
    COUNT(*)::int AS quantidade
FROM Lote
GROUP BY ano, mes
ORDER BY ano, mes;

-- ================================
-- VENDA SERVICE
-- ================================

-- GetAll - Query básica
SELECT id_venda, id_cliente, id_funcionario, data_hora_venda, data_hora_pagamento, tipo_pagamento FROM Venda AS v;

-- GetAll com filtros (campos: data_hora_venda, data_hora_pagamento, tipo_pagamento, id_cliente, id_funcionario)
-- Exemplo com filtros aplicados:
SELECT id_venda, id_cliente, id_funcionario, data_hora_venda, data_hora_pagamento, tipo_pagamento FROM Venda AS v
 WHERE v.tipo_pagamento = $1 AND v.id_cliente = $2 AND v.data_hora_venda >= $3
 ORDER BY data_hora_venda DESC OFFSET $4 LIMIT $5;

-- Create
INSERT INTO Venda (id_cliente, id_funcionario, data_hora_venda, data_hora_pagamento, tipo_pagamento) VALUES ($1, $2, $3, $4, $5) RETURNING id_venda;

-- GetByID
SELECT id_venda, id_cliente, id_funcionario, data_hora_venda, data_hora_pagamento, tipo_pagamento FROM Venda WHERE id_venda = $1;

-- Update
UPDATE Venda SET id_cliente = $1, id_funcionario = $2, data_hora_venda = $3, data_hora_pagamento = $4, tipo_pagamento = $5 WHERE id_venda = $6;

-- Delete
DELETE FROM Venda WHERE id_venda = $1 RETURNING id_venda, id_cliente, id_funcionario, data_hora_venda, data_hora_pagamento, tipo_pagamento;

-- ================================
-- OFERTA SERVICE
-- ================================

-- GetAll - Query básica
SELECT id_oferta, nome, data_criacao, data_inicio, data_fim, valor_fixo, percentual_desconto FROM Oferta AS o;

-- GetAll com filtros (campos: nome, valor_fixo, percentual_desconto, data_criacao, data_inicio, data_fim)
-- Exemplo com filtros aplicados:
SELECT id_oferta, nome, data_criacao, data_inicio, data_fim, valor_fixo, percentual_desconto FROM Oferta AS o
 WHERE o.nome LIKE '%' || $1 || '%' AND o.valor_fixo >= $2 AND o.data_inicio >= $3
 ORDER BY nome, data_criacao DESC OFFSET $4 LIMIT $5;

-- GetByID
SELECT id_oferta, nome, data_criacao, data_inicio, data_fim, valor_fixo, percentual_desconto FROM Oferta WHERE id_oferta = $1;

-- Create
INSERT INTO Oferta (nome, data_inicio, data_fim, valor_fixo, percentual_desconto)
VALUES ($1, $2, $3, $4, $5)
RETURNING id_oferta, data_criacao;

-- Update
UPDATE Oferta SET
nome = $1, data_inicio = $2, data_fim = $3, valor_fixo = $4, percentual_desconto = $5
WHERE id_oferta = $6;

-- Delete
DELETE FROM Oferta WHERE id_oferta = $1 RETURNING id_oferta, nome, data_criacao, data_inicio, data_fim, valor_fixo, percentual_desconto;

-- ================================
-- ITEM_OFERTA SERVICE
-- ================================

-- GetAll - Query básica
SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta as io;

-- GetAll com filtros (campos: quantidade, id_produto, id_oferta)
-- Exemplo com filtros aplicados:
SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta as io
 WHERE io.quantidade >= $1 AND io.id_produto = $2 AND io.id_oferta = $3
 ORDER BY quantidade DESC OFFSET $4 LIMIT $5;

-- GetByComposedID
SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta WHERE id_produto = $1 AND id_oferta = $2;

-- GetAllByItemID
SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta WHERE id_produto = $1;

-- GetAllByOfertaID
SELECT quantidade, id_produto, id_oferta FROM contem_item_oferta WHERE id_oferta = $1;

-- Create
INSERT INTO contem_item_oferta (quantidade, id_produto, id_oferta) VALUES ($1, $2, $3);

-- Update
UPDATE contem_item_oferta SET quantidade = $1 WHERE id_produto = $2 AND id_oferta = $3;

-- Delete
DELETE FROM contem_item_oferta WHERE id_produto = $1 AND id_oferta = $2;

-- ================================
-- ITEM_VENDA SERVICE
-- ================================

-- FindAvailableLote
SELECT
    l.id_lote
FROM
    Lote l
LEFT JOIN (
    SELECT id_lote, SUM(quantidade) as total_vendido
    FROM item_venda
    GROUP BY id_lote
) iv ON l.id_lote = iv.id_lote
WHERE
    l.id_produto = $1
    AND (l.validade IS NULL OR l.validade > CURRENT_DATE)
    AND (l.quantidade_inicial - l.estragados - COALESCE(iv.total_vendido, 0)) >= $2
ORDER BY
    l.validade ASC
LIMIT 1;

-- GetItemsByVendaID
SELECT
    iv.id_item_venda, iv.id_venda, iv.id_lote, iv.quantidade, iv.valor_unitario,
    p.nome, p.marca, l.validade
FROM
    item_venda iv
JOIN
    Lote l ON iv.id_lote = l.id_lote
JOIN
    Produto p ON l.id_produto = p.id_produto
WHERE
    iv.id_venda = $1;

-- GetAll - Query básica
SELECT id_item_venda, id_venda, id_lote, quantidade, valor_unitario FROM item_venda AS IV;

-- GetAll com filtros (campos: id_venda, id_produto, quantidade, valor_unitario)
-- Exemplo com filtros aplicados:
SELECT id_item_venda, id_venda, id_lote, quantidade, valor_unitario FROM item_venda AS IV
 WHERE IV.id_venda = $1 AND IV.quantidade >= $2 AND IV.valor_unitario >= $3
 ORDER BY quantidade DESC OFFSET $4 LIMIT $5;

-- GetByID
SELECT id_item_venda, id_venda, id_lote, quantidade, valor_unitario FROM item_venda WHERE id_item_venda = $1;

-- Create
INSERT INTO item_venda (id_venda, id_lote, quantidade, valor_unitario) VALUES ($1, $2, $3, $4) RETURNING id_item_venda;

-- Update
UPDATE item_venda SET id_venda = $1, id_lote = $2, quantidade = $3, valor_unitario = $4 WHERE id_item_venda = $5;

-- Delete
DELETE FROM item_venda WHERE id_item_venda = $1 RETURNING id_item_venda, id_venda, id_lote, quantidade, valor_unitario;

-- ================================
-- APLICA_OFERTA SERVICE
-- ================================

-- GetByVendaID
SELECT
    ao.id_aplica_oferta, ao.id_oferta, ao.id_venda, ao.id_item_venda,
    o.nome, o.valor_fixo, o.percentual_desconto
FROM
    aplica_oferta ao
JOIN
    Oferta o ON ao.id_oferta = o.id_oferta
WHERE
    ao.id_venda = $1;

-- DeleteByItemVendaID
DELETE FROM aplica_oferta WHERE id_item_venda = $1;

-- GetAll - Query básica
SELECT id_aplica_oferta, id_oferta, id_venda, id_item_venda
FROM aplica_oferta;

-- GetAll com filtros (campos: id_oferta, id_venda, id_item_venda)
-- Exemplo com filtros aplicados:
SELECT id_aplica_oferta, id_oferta, id_venda, id_item_venda
FROM aplica_oferta AS c
 WHERE c.id_oferta = $1 AND c.id_venda = $2 AND c.id_item_venda = $3
 ORDER BY id_oferta DESC OFFSET $4 LIMIT $5;

-- GetByID
SELECT id_aplica_oferta, id_oferta, id_venda, id_item_venda
FROM aplica_oferta
WHERE id_aplica_oferta = $1;

-- Create
INSERT INTO aplica_oferta (id_oferta, id_venda, id_item_venda)
VALUES ($1, $2, $3)
RETURNING id_aplica_oferta;

-- Update
UPDATE aplica_oferta
SET id_oferta = $2, id_venda = $3, id_item_venda = $4
WHERE id_aplica_oferta = $1;

-- Delete
DELETE FROM aplica_oferta
WHERE id_aplica_oferta = $1
RETURNING id_aplica_oferta, id_oferta, id_venda, id_item_venda;

-- ================================
-- RELATORIO SERVICE
-- ================================

-- GetPayrollReport - Query para funcionários no período
SELECT
    id_funcionario,
    nome,
    CPF,
    tipo::text,
    expediente::text,
    salario,
    data_contratacao
FROM Funcionario
WHERE data_contratacao <= $1::date
ORDER BY nome;

-- GetPayrollReport com filtro por tipo
SELECT
    id_funcionario,
    nome,
    CPF,
    tipo::text,
    expediente::text,
    salario,
    data_contratacao
FROM Funcionario
WHERE data_contratacao <= $1::date AND tipo = $2
ORDER BY nome;

-- GetFinancialReport - Query para receita
SELECT date_trunc('day', v.data_hora_venda) AS period,
       COALESCE(SUM(iv.quantidade * iv.valor_unitario), 0) AS receita
FROM Venda v
JOIN item_venda iv ON iv.id_venda = v.id_venda
WHERE v.data_hora_venda::date BETWEEN $1::date AND $2::date
GROUP BY period
ORDER BY period;

-- GetFinancialReport - Query para despesa
SELECT date_trunc('day', l.data_fornecimento) AS period,
       COALESCE(SUM(l.preco_unitario * l.quantidade_inicial), 0) AS despesa
FROM Lote l
WHERE l.data_fornecimento::date BETWEEN $1::date AND $2::date
GROUP BY period
ORDER BY period;

-- ================================
-- EXEMPLOS COM TODOS OS FILTROS APLICADOS
-- ================================

-- Exemplo FORNECEDOR com todos os filtros possíveis:
SELECT id_fornecedor, nome, CNPJ FROM Fornecedor AS f
 WHERE f.nome LIKE '%' || $1 || '%' AND f.CNPJ ILIKE '%' || $2 || '%'
 ORDER BY nome DESC, CNPJ OFFSET $3 LIMIT $4;

-- Exemplo CLIENTE com todos os filtros possíveis:
SELECT id_cliente, nome, cpf, data_nascimento FROM Cliente AS c
 WHERE c.nome LIKE '%' || $1 || '%' AND c.cpf = $2 AND c.data_nascimento >= $3
 ORDER BY nome, data_nascimento DESC OFFSET $4 LIMIT $5;

-- Exemplo PRODUTO COMERCIAL com todos os filtros possíveis:
SELECT p.id_produto, p.nome, p.categoria, p.marca, c.preco_venda
FROM Produto p
INNER JOIN ProdutoComercial c ON p.id_produto = c.id_produto
 WHERE p.nome LIKE '%' || $1 || '%' AND p.categoria ILIKE '%' || $2 || '%'
   AND p.marca = $3 AND c.preco_venda >= $4
 ORDER BY nome, categoria, marca, preco_venda DESC OFFSET $5 LIMIT $6;

-- Exemplo FUNCIONARIO com todos os filtros possíveis:
SELECT id_funcionario, nome, CPF, tipo, expediente, salario, data_contratacao FROM Funcionario AS fc
 WHERE fc.nome LIKE '%' || $1 || '%' AND fc.CPF = $2 AND fc.tipo = $3
   AND fc.expediente = $4 AND fc.salario >= $5 AND fc.data_contratacao <= $6
 ORDER BY nome, data_contratacao DESC OFFSET $7 LIMIT $8;

-- Exemplo LOTE com todos os filtros possíveis:
SELECT id_lote, id_fornecedor, id_produto, data_fornecimento, validade, preco_unitario, estragados, quantidade_inicial FROM Lote AS l
 WHERE l.id_lote = $1 AND l.id_fornecedor = $2 AND l.id_produto = $3
   AND l.preco_unitario >= $4 AND l.estragados <= $5 AND l.quantidade_inicial >= $6
   AND l.validade <= $7
 ORDER BY id_lote, preco_unitario, validade DESC OFFSET $8 LIMIT $9;

-- Exemplo VENDA com todos os filtros possíveis:
SELECT id_venda, id_cliente, id_funcionario, data_hora_venda, data_hora_pagamento, tipo_pagamento FROM Venda AS v
 WHERE v.tipo_pagamento = $1 AND v.id_cliente = $2 AND v.id_funcionario = $3
   AND v.data_hora_venda >= $4 AND v.data_hora_pagamento <= $5
 ORDER BY data_hora_venda, data_hora_pagamento DESC OFFSET $6 LIMIT $7;

-- Exemplo OFERTA com todos os filtros possíveis:
SELECT id_oferta, nome, data_criacao, data_inicio, data_fim, valor_fixo, percentual_desconto FROM Oferta AS o
 WHERE o.nome LIKE '%' || $1 || '%' AND o.valor_fixo >= $2 AND o.percentual_desconto <= $3
   AND o.data_criacao >= $4 AND o.data_inicio >= $5 AND o.data_fim <= $6
 ORDER BY nome, data_criacao DESC OFFSET $7 LIMIT $8;
