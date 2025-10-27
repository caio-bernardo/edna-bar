CREATE TABLE Livro (
    ISBN VARCHAR(20) PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    data_de_publicacao DATE NOT NULL,
    editora_id INT,
    FOREIGN KEY (editora_id) REFERENCES Editora(id)
);
