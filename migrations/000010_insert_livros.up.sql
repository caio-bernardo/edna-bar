-- Inserindo 3 editoras
INSERT INTO Editora (nome, endereco) VALUES
('Editora Alfa', 'Rua das Letras, 123'),
('Editora Beta', 'Avenida do Conhecimento, 456'),
('Editora Gamma', 'Praça da Cultura, 789');

-- Inserindo 20 livros, todos associados às 3 editoras
INSERT INTO Livros (ISBN, titulo, data_de_publicacao, editora_id) VALUES
('978-3-16-148410-0', 'O Livro Alfa', '2020-01-01', 1),
('978-3-16-148411-7', 'O Livro Beta', '2021-02-02', 1),
('978-3-16-148412-4', 'O Livro Gamma', '2022-03-03', 1),
('978-3-16-148413-1', 'O Livro Delta', '2023-04-04', 2),
('978-3-16-148414-8', 'O Livro Épsilon', '2024-05-05', 2),
('978-3-16-148415-5', 'O Livro Zeta', '2020-06-06', 2),
('978-3-16-148416-2', 'O Livro Eta', '2021-07-07', 3),
('978-3-16-148417-9', 'O Livro Theta', '2022-08-08', 3),
('978-3-16-148418-6', 'O Livro Iota', '2023-09-09', 3),
('978-3-16-148419-3', 'O Livro Kappa', '2024-10-10', 1),
('978-3-16-148420-9', 'O Livro Lambda', '2021-11-11', 1),
('978-3-16-148421-6', 'O Livro Mu', '2022-12-12', 1),
('978-3-16-148422-3', 'O Livro Nu', '2023-01-13', 2),
('978-3-16-148423-0', 'O Livro Xi', '2024-02-14', 2),
('978-3-16-148424-7', 'O Livro Omicron', '2020-03-15', 2),
('978-3-16-148425-4', 'O Livro Pi', '2021-04-16', 3),
('978-3-16-148426-1', 'O Livro Rho', '2022-05-17', 3),
('978-3-16-148427-8', 'O Livro Sigma', '2023-06-18', 3),
('978-3-16-148428-5', 'O Livro Tau', '2024-07-19', 1),
('978-3-16-148429-2', 'O Livro Upsilon', '2020-08-20', 1);

-- Inserindo 20 autores
INSERT INTO Autor (RG, nome, endereco) VALUES
('1234567890', 'Carlos Silva', 'Rua das Flores, 10'),
('2345678901', 'Maria Oliveira', 'Avenida dos Ipês, 20'),
('3456789012', 'João Souza', 'Praça das Árvores, 30'),
('4567890123', 'Ana Paula', 'Rua das Rosas, 40'),
('5678901234', 'Pedro Santos', 'Avenida das Hortências, 50'),
('6789012345', 'Paula Ribeiro', 'Praça das Orquídeas, 60'),
('7890123456', 'Ricardo Nascimento', 'Rua das Palmeiras, 70'),
('8901234567', 'Fernanda Costa', 'Avenida das Mangueiras, 80'),
('9012345678', 'Luiz Alves', 'Praça das Acácias, 90'),
('0123456789', 'Cláudia Vieira', 'Rua das Bromélias, 100'),
('1123456789', 'Fábio Castro', 'Avenida dos Jacarandás, 110'),
('2123456789', 'Juliana Souza', 'Praça dos Pinhais, 120'),
('3123456789', 'Renato Pinto', 'Rua das Aroeiras, 130'),
('4123456789', 'Tatiana Rocha', 'Avenida dos Cedros, 140'),
('5123456789', 'Paulo Almeida', 'Praça dos Ipês, 150'),
('6123456789', 'Roberta Lima', 'Rua das Macieiras, 160'),
('7123456789', 'Sérgio Lopes', 'Avenida das Araucárias, 170'),
('8123456789', 'Alessandra Costa', 'Praça dos Flamboyants, 180'),
('9123456789', 'Carlos Gomes', 'Rua das Paineiras, 190'),
('1012345678', 'Mariana Pereira', 'Avenida dos Coqueiros, 200');

-- Associando autores e livros na tabela Escreve
INSERT INTO Escreve (ISBN, RG) VALUES
('978-3-16-148410-0', '1234567890'),
('978-3-16-148411-7', '2345678901'),
('978-3-16-148412-4', '3456789012'),
('978-3-16-148413-1', '4567890123'),
('978-3-16-148414-8', '5678901234'),
('978-3-16-148415-5', '6789012345'),
('978-3-16-148416-2', '7890123456'),
('978-3-16-148417-9', '8901234567'),
('978-3-16-148418-6', '9012345678'),
('978-3-16-148419-3', '0123456789'),
('978-3-16-148420-9', '1123456789'),
('978-3-16-148421-6', '2123456789'),
('978-3-16-148422-3', '3123456789'),
('978-3-16-148423-0', '4123456789'),
('978-3-16-148424-7', '5123456789'),
('978-3-16-148425-4', '6123456789'),
('978-3-16-148426-1', '7123456789'),
('978-3-16-148427-8', '8123456789'),
('978-3-16-148428-5', '9123456789'),
('978-3-16-148429-2', '1012345678');
