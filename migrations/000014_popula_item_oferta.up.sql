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
