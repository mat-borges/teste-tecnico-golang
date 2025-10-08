# build
FROM golang:1.25.1-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# runtime
FROM alpine:3.18
WORKDIR /app

# 🔹 Variáveis de ambiente (com defaults)
ENV SERVER_PORT=8080
ENV USERS_BASE_URL=https://jsonplaceholder.typicode.com/users
ENV POSTS_BASE_URL=https://jsonplaceholder.typicode.com/posts
ENV HTTP_TIMEOUT=8s
ENV AGG_TIMEOUT=6s
ENV ENABLE_INTROSPECTION=1
ENV ENABLE_APQ=1

COPY --from=builder /app/server .
EXPOSE 8080

CMD ["./server"]
