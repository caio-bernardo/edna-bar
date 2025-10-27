CREATE TABLE Imprime (
    lisbn VARCHAR(20),
    grafica_id INT,
    nto_copias INT,
    data_entrega DATE,
    PRIMARY KEY (lisbn, grafica_id),
	FOREIGN KEY (lisbn) REFERENCES Livros(isbn),
	FOREIGN KEY (grafica_id) REFERENCES Grafica(id)
);
