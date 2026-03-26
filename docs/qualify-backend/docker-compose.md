# Docker Compose - Qualify Backend

Esta documentação explica o fluxo de `docker-compose` para o backend Qualify (PostgreSQL + migrações + app).

## Estrutura do `docker-compose.yml`

- `db`:
  - Imagem: `postgres:17-alpine`
  - Usuário: `gouser`
  - Senha: `gopassword`
  - DB: `godb`
  - Porta: `5432`
  - Volume: `pgdata` (persistência do banco)

- `migrate`:
  - Imagem: `migrate/migrate:v4.16.0`
  - Dependência: `db`
  - Monta migrations do host: `./internal/database/migrations:/migrations`
  - Comando default: `migrate -path /migrations -database "postgres://gouser:gopassword@db:5432/godb?sslmode=disable" up`

- `app`:
  - Build: `.` (diretório `qualify-backend`)
  - Porta: `8001`
  - Dependência: `migrate` (garante migrações aplicadas antes)
  - Ambiente: `DATABASE_URL=postgres://gouser:gopassword@db:5432/godb?sslmode=disable`

- `volumes`:
  - `pgdata`: volume nomeado para dados do postgres.

## O que acontece em `docker compose up --build`

1. Inicia `db`.
2. Inicia `migrate` depois de `db` e aplica `up` nas migrations (tabelas e esquema).
3. Inicia `app` depois que `migrate` termina com sucesso.

> O app lê `DATABASE_URL` em `cmd/main.go` e consulta a tabela `widgets` (id=42). 

## Comandos úteis

- Iniciar (build + run):
  - `docker compose up --build`
- Ver logs:
  - `docker compose logs -f app`
- Parar/remover:
  - `docker compose down`
- Resetar migrações (dev):
  - `docker compose run --rm migrate down` (rollback último batch)
  - `docker compose run --rm migrate force <version>` (forçar versão)

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

