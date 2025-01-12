# Используем официальный образ Go
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию
WORKDIR /Rates
# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Переходим в директорию cmd и собираем приложение
WORKDIR /Rates/cmd
RUN go build -o main .

# Запуск приложения
CMD ["./main"]
