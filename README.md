# 🧩 Go GraphQL Aggregator

Este projeto é um **desafio técnico em Go** que expõe uma **API GraphQL** responsável por consolidar dados de duas APIs REST públicas:

- **Users:** [https://jsonplaceholder.typicode.com/users](https://jsonplaceholder.typicode.com/users)
- **Posts:** [https://jsonplaceholder.typicode.com/posts](https://jsonplaceholder.typicode.com/posts)

---

## 🚀 Funcionalidade principal

Para um determinado `userId`, a API retorna:

- **name** — Nome do usuário
- **email** — Email do usuário
- **postCount** — Quantidade de posts que ele possui

### Query GraphQL

```graphql
query {
	userSummary(userId: 1) {
		name
		email
		postCount
	}
}
```

---

## 🧠 Stack técnica

- **Go 1.25.1+**
- **gqlgen** — geração do schema e resolvers GraphQL
- **errgroup.WithContext** — concorrência segura com cancelamento automático
- **net/http** — servidor e clientes HTTP configurados com timeouts
- **log/slog** — logging estruturado nativo
- **Docker + docker-compose**
- **Middleware customizado** — logs HTTP e panic recovery
- **Makefile** - automatização de testes e build

## ⚙️ Arquitetura e design

- As duas APIs REST são consumidas **concorrentemente** usando `errgroup.WithContext`.
- Cada requisição tem **timeout configurável (`AGG_TIMEOUT`)**.
- **HTTP client** otimizado com reuso de conexões.
- **Logs estruturados** com `log/slog`, suportando `text`, `json` ou `silent`.
- **Variáveis de ambiente centralizadas** via `internal/config`.
- **Graceful shutdown** e **tratamento de panic** integrados.

---

## ⚙️ Configuração de ambiente

O projeto lê as configurações de um arquivo `.env` (exemplo abaixo):

```dotenv
SERVER_PORT=8080
USERS_BASE_URL=https://example.com/users
POSTS_BASE_URL=https://example.com/posts
HTTP_TIMEOUT=8s
AGG_TIMEOUT=6s
ENABLE_INTROSPECTION=1
ENABLE_APQ=1
LOG_MODE=text
```

---

## 💻 Rodando o projeto

### Localmente

```bash
git clone https://github.com/mat-borges/go-graphql-aggregator.git
cd go-graphql-aggregator

go mod tidy
go run cmd/api/main.go
```

---

### Com Docker

```bash
docker build -t go-graphql-aggregator .
docker run -p 8080:8080 --env-file .env go-graphql-aggregator
```

Ou com Docker Compose:

```bash
docker-compose up --build
```

---

Acesse o playground:
👉 [http://localhost:8080/](http://localhost:8080/)

---

## 🧪 Testes

O projeto possui um **Makefile** que automatiza os testes e geração de relatórios de cobertura.
Ele foca apenas nos pacotes principais (`aggregator`, `graph` e `fetchers`) para resultados relevantes.

### Manualmente

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Makerfile

#### 🔹 Rodar testes silenciosos

```bash
make test
```

#### 🔹 Gerar relatório HTML de cobertura

```bash
make test-report
```

Isso cria o arquivo `coverage.html`, que pode ser aberto no navegador.

#### 🔹 Limpar arquivos temporários

```bash
make clean
```

### 🪄 Como usar o Makefile no Windows

Se aparecer o erro:

```
'make' não é reconhecido como um comando interno
ou externo...
```

➡️ Soluções:

1. **Usar Git Bash**: abra o terminal Git Bash e rode `make test`.
2. **Instalar make via Chocolatey:**

   ```bash
   choco install make
   ```

3. **Ou executar os comandos Go manualmente:**

   ```bash
   go test ./internal/aggregator/... ./internal/graph/... ./internal/fetchers/... -v -cover
   ```

---

## 📊 Logs

Configure o formato via variável de ambiente `LOG_MODE`:

| Valor           | Descrição                    |
| --------------- | ---------------------------- |
| `text` (padrão) | Logs legíveis no terminal    |
| `json`          | Estruturado para produção    |
| `silent`        | Silencia logs durante testes |

---

## 🧩 Estrutura resumida

```
cmd/
  api/main.go     → Inicialização do servidor
internal/
  aggregator/     → lógica de agregação e concorrência
  fecther/        → comunicação HTTP com APIs externas
  config/         → configurações via env
  graph/          → schema e resolvers GraphQL (gqlgen)
  middleware/     → logger HTTP e recovery
  logger/         → setup do slog global
  test/           → inicialização do servidor
Makefile          → automação de testes e build
```

---

## 💬 Notas finais

- Código segue princípios de **Clean Code** e **SOLID**.
- As requisições externas são **concorrentes** e canceláveis.
- `errgroup.WithContext` cancela a outra chamada automaticamente em caso de erro.
- Cobertura de testes alta nos módulos críticos.
- Docker e Makefile prontos para CI/CD e execução local.
- Logs estruturados, prontos para observabilidade.

---

**Autor:** Mateus Borges
📍 Campinas/SP — [github.com/mat-borges](https://github.com/mat-borges)
---
