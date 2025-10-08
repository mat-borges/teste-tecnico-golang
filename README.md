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

---

## ⚙️ Arquitetura e design

- As duas APIs REST são consumidas **concorrentemente** usando `errgroup.WithContext`.
- Cada requisição é limitada por **timeout configurável (`AGG_TIMEOUT`)**.
- **HTTP client** otimizado com pool de conexões e timeouts de rede.
- **Logs estruturados** (texto ou JSON via `LOG_MODE`), incluindo tempos de execução.
- **Variáveis de ambiente centralizadas** via `internal/config/config.go`.
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

Acesse o playground:
👉 [http://localhost:8080/](http://localhost:8080/)

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

## 🧪 Testes

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

---

## 📊 Logs

Modo texto (padrão):

```
time=2025-10-08T15:25:42Z level=INFO msg="fetch user done" userId=1
```

Modo JSON:

```bash

```

---

## 🧩 Estrutura resumida

```
cmd/
	api/main.go   → Inicialização do servidor
internal/
  aggregator/     → Lógica de agregação e concorrência
  config/         → Configurações via env
  graph/          → Schema e resolvers gqlgen
  middleware/     → Logs HTTP, recovery
  logger/         → Configuração global de slog
```

---

## 💬 Notas finais

- As requisições externas são **concorrentes** e respeitam contexto e timeout.
- `errgroup.WithContext` cancela a outra chamada automaticamente em caso de erro.
- Código segue princípios de **Clean Code e SOLID**.
- Preparado para **deploy em container** e observabilidade básica.
