-- migrations/000009_seed_initial_data.up.sql

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

-- Insere Produtos Comerciais (em DUAS etapas)

-- Inserir a parte base em 'Produto'
INSERT INTO Produto (nome, categoria, marca) VALUES
('Cerveja 600ml', 'Bebidas', 'Skoll'),
('Porção de Fritas Grande', 'Porções', 'Casa'),
('Refrigerante Lata 350ml', 'Bebidas', 'Guarana Mineiro');

-- Inserir a parte comercial em 'ProdutoComercial' referenciando a Etapa 1
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

-- Insere Lotes (depende de Fornecedor e Produto)
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
('Happy Hour Cerveja', '2025-11-01', '2025-11-30', 20); -- 20% de desconto

INSERT INTO Oferta (nome, data_inicio, data_fim, valor_fixo) VALUES
('Combo Fritas + Refri', '2025-11-01', '2025-11-30', 28.00); -- Preço fixo 28 reais

-- Insere Itens da Oferta (depende de Oferta e ProdutoComercial)
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

-- Insere Vendas (depende de Cliente e Funcionario)
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
