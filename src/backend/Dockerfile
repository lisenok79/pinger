# FROM golang:1.23.4 AS builder
# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o backend .

# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /app/backend .

# CMD [ "./backend" ]

# Базовый образ с Go
FROM golang:1.23.4

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Компилируем приложение
RUN go build -o backend .

# Указываем порт для прослушивания
EXPOSE 8080

# Указываем команду для запуска приложения
CMD ["./backend"]
