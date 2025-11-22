DELETE FROM aplica_oferta 
WHERE id_venda IN (
    SELECT id_venda FROM Venda WHERE id_cliente IN (
        SELECT id_cliente FROM Cliente WHERE nome IN ('Serginho Groisman', 'Neymar Jr.')
    )
);

DELETE FROM item_venda 
WHERE id_venda IN (
    SELECT id_venda FROM Venda WHERE id_cliente IN (
        SELECT id_cliente FROM Cliente WHERE nome IN ('Serginho Groisman', 'Neymar Jr.')
    )
);
