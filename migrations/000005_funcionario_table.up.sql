DROP TYPE IF EXISTS tipo_de_funcionario;
CREATE TYPE tipo_de_funcionario AS ENUM ('garcom', 'seguranca',
'caixa', 'faxineiro', 'balconista');

DROP TYPE IF EXISTS tipo_de_expediente;
CREATE TYPE tipo_de_expediente AS ENUM ('manha', 'tarde', 'noite', 'madrugada');

CREATE TABLE IF NOT EXISTS Funcionario (
    id_funcionario serial PRIMARY KEY,
    nome varchar(50) NOT NULL,
    CPF char(11) NOT NULL,
    tipo tipo_de_funcionario NOT NULL,
    expediente tipo_de_expediente NOT NULL,
    data_contratacao date NOT NULL,
    salario decimal(8, 2) NOT NULL
);
