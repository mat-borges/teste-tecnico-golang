# 🧩 Go GraphQL Aggregator

Este projeto é um desafio técnico que expõe uma API GraphQL responsável por consolidar dados de duas APIs REST públicas:

- [Users](https://jsonplaceholder.typicode.com/users)
- [Posts](https://jsonplaceholder.typicode.com/posts)

## 🚀 Funcionalidades

Para um `userId`, retorna:

- Nome do usuário
- Email
- Quantidade de posts

## 🧠 Stack

- Go 1.25.1
- gqlgen
- net/http
- Docker

---

## 🧩 Como rodar o projeto (modo local)

```bash
# Clonar repositório
git clone https://github.com/seuusuario/go-graphql-aggregator.git
cd go-graphql-aggregator

# Instalar dependências
go mod tidy

# Rodar servidor
go run cmd/main.go
```
