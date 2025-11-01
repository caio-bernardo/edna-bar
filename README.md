# E.D.N.A - Ecossistema D. Negociação Alcoólica

Sistema de banco de dados de Bar. Gerencie produtos, atendimentos e outras atividades do seu bar com nosso sistema.

> Esse projeto é parte do projeto final para a disciplina *Banco de Dados I*.

## Funcionalidades
- [ ]

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

Inicie a base de dados com `make docker-run` (use `sudo` se necessário). Vá para a migração mais recente (se já não estiver) com `./migrate.sh up`.

Por último rode projeto com `make run`, para rodar com _hot reloading_ (alterações serão atualizadas instantâneamente) use `make watch`.

## Migrações

> Migrações são scripts SQL que são rodados na base de dados e permitem criar um histórico de alterações e navegar por elas.

Não esqueça de acionar a base de dados antes de rodar migrações.

As migrações vivem em `migrations`. Utilize o script `migrate.sh` para gerenciar migrações. Crie novas migrações com `./migrate.sh create <migration_name>`. Vá para a migração mais recente com `./migrate.sh up`, volte uma migração com `./migrate.sh down 1`. Veja mais comandos em `./migrate.sh help`.


## Diagrama Conceitual

![Diagrama Entidade-Relacionamento do Sistema EDNA bar](assets/E.D.N.A_conceitual_1.png)

## MakeFile

Compile e rode os testes.
```bash
make all
```

Compile a aplicação
```bash
make build
```

Rode a aplicação
```bash
make run
```

Cria o banco de dados por um _container_ docker:
```bash
make docker-run
```

Desativa o container:
```bash
make docker-down
```

Testes de integração no Banco de Dados:
```bash
make itest
```

Live reload aplicação:
```bash
make watch
```

Rode os testes:
```bash
make test
```

Remova o executavél
```bash
make clean
```
