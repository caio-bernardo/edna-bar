CREATE TABLE IF NOT EXISTS contem_item_venda (
    valor_unitario decimal(6, 2),
    quantidade int,
    id_oferta int,
    id_produto int,
    id_venda int,
    PRIMARY KEY (id_produto, id_venda),
    FOREIGN KEY (id_produto) REFERENCES ProdutoComercial(id_produto),
    FOREIGN KEY (id_venda) REFERENCES Venda(id_venda) ON DELETE CASCADE,
    FOREIGN KEY (id_oferta) REFERENCES Oferta(id_oferta) ON DELETE SET NULL
);
