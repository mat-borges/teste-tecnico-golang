# 🧩 Go GraphQL Aggregator

Este projeto é um desafio técnico que expõe uma API GraphQL responsável por consolidar dados de duas APIs REST públicas:

- Users: <https://jsonplaceholder.typicode.com/users>
- Posts: <https://jsonplaceholder.typicode.com/posts>

## 🚀 Funcionalidades

Para um `userId`, retorna

- Nome do usuário
- Email
- Quantidade de posts

Query GraphQL:

```graphql
query {
	userSummary(userId: 1) {
		name
		email
		postCount
	}
}
```

## 🧠 Stack

- Go 1.25.1+ (testado com 1.25.1)
- gqlgen
- net/http
- Docker

---

## 🧩 Como rodar

### Local

```bash
# Clonar repositório
git clone https://github.com/mat-borges/go-graphql-aggregator.git
cd go-graphql-aggregator

# Instalar dependências
go mod tidy

# Rodar servidor
go run cmd/api/main.go
```

Acesse: [http://localhost:8080/](http://localhost:8080/) (playground)

### Docker

Buildar e subir:

```bash
docker build -t go-graphql-aggregator .
docker run -p 8080:8080 -e PORT=8080 go-graphql-aggregator
```

Ou com docker-compose:

```bash
docker-compose up --build
```

## Testes e cobertura

Rode:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Notas de performance e arquitetura

- As chamadas externas são feitas de forma concorrente e limitadas por contexto/timeout.
- Para melhorar performance, os `posts` são buscados utilizando `?userId=...`, reduzindo payload.
- `http.Client` configurado com timeouts e transporte mais robusto.
