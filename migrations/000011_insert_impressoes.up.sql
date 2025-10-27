-- Inserindo 3 gráficas
INSERT INTO Grafica (nome) VALUES
('Grafica Alfa'), ('Grafica Beta'), ('Grafica Gamma');

-- Inserindo algumas gráficas como particulares e contratadas
INSERT INTO Particular (grafica_id) VALUES (1);

INSERT INTO Contratada (grafica_id, endereco) VALUES
(2, 'Avenida Central, 111'),
(3, 'Rua dos Cravos, 222');

-- Inserindo contratos relacionados às gráficas contratadas
INSERT INTO Contrato (valor, nome_responsavel, grafica_cont_id) VALUES
(10000.00, 'Carlos Almeida', 2),
(15000.50, 'Maria Ferreira', 3);


INSERT INTO imprime (lisbn, grafica_id, nto_copias, data_entrega) VALUES
('978-3-16-148410-0', 1, 100, '2024-11-01'),
('978-3-16-148411-7', 2, 150, '2024-11-03'),
('978-3-16-148412-4', 3, 200, '2024-11-05'),
('978-3-16-148413-1', 1, 120, '2024-11-07'),
('978-3-16-148414-8', 2, 130, '2024-11-10'),
('978-3-16-148415-5', 3, 160, '2024-11-12'),
('978-3-16-148416-2', 1, 180, '2024-11-15'),
('978-3-16-148417-9', 2, 140, '2024-11-18'),
('978-3-16-148418-6', 3, 170, '2024-11-20'),
('978-3-16-148419-3', 1, 190, '2024-11-22');
