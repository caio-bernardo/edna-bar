CREATE TABLE Contratada (
    grafica_id INT PRIMARY KEY,
    endereco VARCHAR(255) NOT NULL,
    FOREIGN KEY (grafica_id) REFERENCES Grafica(id)
);
