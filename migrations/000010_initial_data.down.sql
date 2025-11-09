-- migrations/000009_seed_initial_data.down.sql

-- TRUNCATE limpa todas as tabelas listadas
-- RESTART IDENTITY reseta os contadores 'serial' (ex: id_produto voltará a ser 1)
-- CASCADE lida automaticamente com as dependências de chave estrangeira
TRUNCATE 
    Fornecedor,
    Produto, -- Limpa Produto e ProdutoComercial (devido à herança)
    Cliente,
    Funcionario,
    Lote,
    Oferta,
    Venda,
    contem_item_oferta
RESTART IDENTITY CASCADE;
