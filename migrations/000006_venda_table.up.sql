DROP TYPE IF EXISTS tipo_de_pagamento;
-- Adicione 'fiado' na lista abaixo
CREATE TYPE tipo_de_pagamento AS ENUM ('credito', 'debito', 'pix', 'dinheiro', 'VA/VR', 'fiado');

CREATE TABLE IF NOT EXISTS Venda (
    id_venda serial PRIMARY KEY,
    data_hora_venda timestamp NOT NULL DEFAULT now(),
    data_hora_pagamento timestamp,
    tipo_pagamento tipo_de_pagamento,

    id_cliente int NOT NULL,
    id_funcionario int NOT NULL,

    FOREIGN KEY (id_cliente) REFERENCES Cliente(id_cliente) ON DELETE CASCADE,
    FOREIGN KEY (id_funcionario) REFERENCES Funcionario(id_funcionario) ON DELETE SET NULL
);
