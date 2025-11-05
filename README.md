# E.D.N.A - Ecossistema D. Negociação Alcoólica

Sistema de banco de dados de Bar. Gerencie produtos, atendimentos e outras atividades do seu bar com nosso sistema.

> Esse projeto é parte do projeto final para a disciplina *Banco de Dados I*.

## Funcionalidades

- Controle (leitura, criação, modificação, etc...):
  - [x] Fornecedores
  - [x] Produtos
  - [ ] Lotes
  - [ ] Clientes
  - [ ] Funcionários
  - [ ] Vendas
  - [ ] Ofertas
- [x] Documentação interativa da API com Swagger UI
- [ ] Frontend

## Tecnologias
- Golang 1.24
- Postgres
- Docker & Docker Compose
- go-migrate

## Como rodar

Configure as variáveis de ambiente em um arquivo `.env` seguindo os exemplos em `.env.example`.


Primeiro, instale a versão `1.24` da linguagem Go, a versão mais recente do Docker e Docker-compose. Em adendo para realizar as migrações é preciso instalar o `go-migrate`, você pode fazer isso com:
```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3
```

Inicie a base de dados com `make docker-run` (use `sudo` se necessário). Vá para a migração mais recente (se já não estiver) com `./migrate.sh up`, isso irá criar as tabelas no seu banco de dados (e populá-lo com dados), se quiser saber mais sobre o conceito de migrações veja a próxima seção.

Por último rode projeto com `make run`. Para rodar com _hot reloading_ (alterações serão refletidas instantâneamente) use `make watch`.

## Migrações

> Migrações são scripts SQL que são rodados na base de dados e permitem criar um histórico de alterações e navegar por elas.

Não esqueça de acionar a base de dados antes de rodar migrações.

As migrações vivem em `migrations`. Utilize o script `migrate.sh` para gerenciar migrações. Crie novas migrações com `./migrate.sh create <migration_name>`. Vá para a migração mais recente com `./migrate.sh up`, volte **uma** migração com `./migrate.sh down 1`. Veja mais comandos em `./migrate.sh help`.


## Diagrama Conceitual

![Diagrama Entidade-Relacionamento do Sistema EDNA bar](assets/E.D.N.A_conceitual_1.png)
