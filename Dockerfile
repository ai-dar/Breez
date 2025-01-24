# Используем официальное изображение Golang
FROM golang:1.22

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем весь проект в контейнер
COPY . .
RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

# Устанавливаем зависимости и собираем приложение
RUN go mod tidy && go build -o main .

# Открываем порт для приложения
EXPOSE 8080

# Запускаем приложение

CMD ["./main"]
