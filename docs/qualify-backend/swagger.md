# Swagger - Qualify Backend

Este arquivo descreve como gerar e visualizar a documentação Swagger (OpenAPI) do backend Qualify.

## Requisitos

- Go instalado e configurado (GOPATH).
- Ter o diretório `$GOPATH/bin` no `PATH` para acessar o binário `swag`.

## Instalar o CLI `swag`

Execute (uma vez):

```bash
cd qualify-backend
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

Observação: se preferir não alterar o `PATH` global, use o caminho completo `$(go env GOPATH)/bin/swag` ao invocar o comando.

### Instalação local (por projeto)

Se preferir manter o binário dentro do repositório/projeto (sem instalar globalmente), há duas abordagens comuns:

- Usando `GOBIN` para instalar no diretório `bin` do projeto:

```bash
cd qualify-backend
mkdir -p bin
GOBIN=$(pwd)/bin go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(pwd)/bin
# Agora `swag` estará disponível como ./bin/swag no projeto
```

- Baixando a release oficial e colocando em `bin` (útil quando não há Go disponível):

```bash
cd qualify-backend
mkdir -p bin && cd bin
# escolha o asset compatível com sua arquitetura em https://github.com/swaggo/swag/releases
curl -L -O https://github.com/swaggo/swag/releases/latest/download/swag_linux_amd64.tar.gz
tar -xzf swag_linux_amd64.tar.gz
rm swag_linux_amd64.tar.gz
cd -
export PATH=$PATH:$(pwd)/bin
```

Observação: substitua `swag_linux_amd64.tar.gz` pelo asset correto para sua plataforma (por exemplo `swag_darwin_amd64.tar.gz` em macOS). A URL exata depende da release e arquitetura.

## Multiplataforma (Linux / macOS / Windows)

As instruções acima funcionam na maioria das plataformas; abaixo há passos específicos para cada sistema.

- Linux / macOS (com Go instalado):

```bash
cd qualify-backend
go install github.com/swaggo/swag/cmd/swag@latest
# Em macOS/Linux, adicione: export PATH=$PATH:$(go env GOPATH)/bin
```

- macOS (sem Go): baixe o asset `swag_darwin_amd64.tar.gz`/`swag_darwin_arm64.tar.gz`, extraia em `bin/` e `chmod +x bin/swag`.

- Windows (PowerShell) com Go instalado:

```powershell
cd qualify-backend
go install github.com/swaggo/swag/cmd/swag@latest
# para uso na sessão atual:
$env:PATH += ";$(go env GOPATH)\bin"
# ou para instalar local no projeto:
mkdir bin; $env:GOBIN = (Resolve-Path .\bin).Path; go install github.com/swaggo/swag/cmd/swag@latest
```

- Windows (sem Go): baixe o asset `swag_windows_amd64.zip` correspondente, extraia para `bin\` e adicione `bin` ao `PATH` (temporário com `$env:PATH += ";$(Resolve-Path .\bin)"`, permanente via `setx`).

Observação sobre assets: escolha o arquivo correto para sua arquitetura, por exemplo `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64`, `windows_amd64`.

Exemplo de comandos para download e extração (Linux/macOS):

```bash
cd qualify-backend
mkdir -p bin && cd bin
curl -L -O https://github.com/swaggo/swag/releases/latest/download/swag_linux_amd64.tar.gz
tar -xzf swag_linux_amd64.tar.gz
chmod +x swag
cd -
export PATH=$PATH:$(pwd)/bin
```

Exemplo para Windows (PowerShell):

```powershell
cd qualify-backend
mkdir bin; cd bin
(New-Object System.Net.WebClient).DownloadFile('https://github.com/swaggo/swag/releases/latest/download/swag_windows_amd64.zip','swag_windows_amd64.zip')
Expand-Archive swag_windows_amd64.zip -DestinationPath .
cd ..
$env:PATH += ";$(Resolve-Path .\bin)"
```

## Gerar a documentação

No diretório raiz do serviço backend execute:

```bash
cd qualify-backend
swag init --parseDependency --parseInternal --dir ./cmd,./internal/database/handlers,./pkg \
	--generalInfo main.go --output ./docs
```

Parâmetros principais:
- `--dir`: diretórios a serem parseados para comentários de anotação.
- `--generalInfo`: arquivo que contém as informações gerais da API (ex.: `main.go`).
- `--output`: pasta onde os artefatos gerados serão gravados (`docs/`).
- `--parseDependency` / `--parseInternal`: ajudam a incluir pacotes internos e dependências na análise.

## Arquivos gerados

Após a execução, a pasta `qualify-backend/docs` deverá conter pelo menos:

- `swagger.yaml` e/ou `swagger.json` — especificação OpenAPI gerada.
- `docs.go` — arquivo Go com informações da spec (usado por bibliotecas como `http-swagger`).

Verifique com:

```bash
ls qualify-backend/docs
```

## Anotações e comentários

Para que o `swag` gere corretamente os endpoints, comente os handlers e modelos seguindo o formato de comentários do swag, por exemplo:

```go
// @Summary Cria um novo usuário
// @Description Cria um usuário com os dados fornecidos
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User payload"
// @Success 201 {object} models.User
// @Security     BearerAuth
// @Router /users [post]
func CreateUser(c *gin.Context) { ... }
```

## Dicas de troubleshooting

- `swag: command not found`: verifique se `go install` executou com sucesso e se `$(go env GOPATH)/bin` está no `PATH`.
- Se endpoints não aparecerem: confirme que os comentários estão no padrão `@Summary`, `@Param`, `@Success`, `@Router` etc. e que os diretórios corretos foram passados em `--dir`.

## Referências rápidas

- Comando de geração recomendado:

```bash
swag init --parseDependency --parseInternal --dir ./cmd,./internal/database/handlers,./pkg --generalInfo main.go --output ./docs
```

- Exemplo de execução do Swagger UI via Docker: veja a seção "Como visualizar a documentação" acima.

---
Se quiser, eu executo a instalação do `swag` e gero os arquivos aqui no ambiente (se permitir), ou adapto o texto para outro idioma/estilo.

PATH="$PATH:$(go env GOBIN):$(go env GOPATH)/bin" swag init --parseDependency --parseInternal --dir ./cmd,./internal/database/handlers,./pkg --generalInfo main.go --output ./docs