# ========================
# Этап сборки
# ========================
FROM golang:1.25-alpine AS builder

# Рабочая директория
WORKDIR /app

# Копируем файлы модулей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект (все .go файлы должны быть в одном пакете main)
COPY . .

# Собираем бинарник
RUN go build -o subscription-service .

# ========================
# Финальный образ
# ========================
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

# Копируем бинарник из билд-стейджа
COPY --from=builder /app/subscription-service /subscription-service

# Копируем .env
COPY .env ./

# Открываем порт 8080
EXPOSE 8080

# Запуск сервиса
ENTRYPOINT ["/subscription-service"]
