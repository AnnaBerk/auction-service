FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файл конфигурации
COPY config.toml /app/config.toml

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -o auction ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auction .
# Копируем конфиг на финальный этап
COPY --from=builder /app/config.toml /app/config.toml

EXPOSE 8080

CMD ["./auction"]








