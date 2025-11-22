DELETE FROM contem_item_oferta
WHERE id_oferta = (SELECT id_oferta FROM Oferta WHERE nome = 'Happy Hour Cerveja')
  AND id_produto = (SELECT id_produto FROM Produto WHERE nome = 'Cerveja 600ml');

DELETE FROM contem_item_oferta
WHERE id_oferta = (SELECT id_oferta FROM Oferta WHERE nome = 'Combo Fritas + Refri')
  AND id_produto = (SELECT id_produto FROM Produto WHERE nome = 'Porção de Fritas Grande');

DELETE FROM contem_item_oferta
WHERE id_oferta = (SELECT id_oferta FROM Oferta WHERE nome = 'Combo Fritas + Refri')
  AND id_produto = (SELECT id_produto FROM Produto WHERE nome = 'Refrigerante Lata 350ml');
