# Билд-этап (создаём бинарник)
FROM golang:1.23 AS builder

WORKDIR /app
COPY ./ ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/

# Финальный образ
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY .env .env

EXPOSE 8080
CMD ["/app/main"]
