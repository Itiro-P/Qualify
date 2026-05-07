# Docker Compose - Qualify Backend

Esta documentação explica o fluxo de `docker-compose` para o backend Qualify (PostgreSQL + migrações + app).

## O que acontece em `docker compose up --build`

1. Inicia `db`.
2. Inicia `migrate` depois de `db` e aplica `up` nas migrations (tabelas e esquema).
3. Inicia `app` depois que `migrate` termina com sucesso.

> O app lê `DATABASE_URL` em `cmd/main.go` e consulta a tabela `user` (id=42). 

## Comandos úteis

- Iniciar (build + run):
  - `docker compose up --build`
- Ver logs:
  - `docker compose run logs -f app`
- Parar/remover:
  - `docker compose down`
- Resetar migrações (dev):
  - `docker compose --profile tools run --rm migrate down <número de migrations>` (rollback último batch)
  - `docker compose --profile tools run --rm migrate force <version>` (forçar versão)

## Migrations

- Rode o comando:
  - `docker compose run --rm --user "$(id -u):$(id -g)" migrate create -ext=sql -dir=/migrations -seq init`
- Exemplo `up`:
  ```sql
  CREATE TABLE user (
    id bigint PRIMARY KEY,
  );
  ```
- Exemplo `down`:
  ```sql
  DROP TABLE user;
  ```

## Ajustes comuns

- `cmd/main.go` usa `os.Getenv("DB_URL")`; manter variável igual no compose.
- Se migrator deve executar somente manualmente, remova `depends_on` de `app` ou use `docker compose run --rm migrate` fora do fluxo de startup.

