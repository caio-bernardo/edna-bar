CREATE TABLE Particular (
    grafica_id INT PRIMARY KEY,
    FOREIGN KEY (grafica_id) REFERENCES Grafica(id)
);
