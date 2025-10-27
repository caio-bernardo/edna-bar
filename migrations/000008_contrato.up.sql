CREATE TABLE Contrato (
    id SERIAL PRIMARY KEY,
    valor DECIMAL(10, 2) NOT NULL,
    nome_responsavel VARCHAR(255) NOT NULL,
    grafica_cont_id INT,
    FOREIGN KEY (grafica_cont_id) REFERENCES Contratada(grafica_id)
);
