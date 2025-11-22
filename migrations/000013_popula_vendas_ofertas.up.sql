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
