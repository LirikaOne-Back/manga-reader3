FROM golang:1.22-alpine AS builder

WORKDIR /app

# Устанавливаем необходимые зависимости
RUN apk --no-cache add ca-certificates git tzdata gcc musl-dev

# Настройка прямой загрузки модулей
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on
ENV GOSUMDB=off

# Копируем Go модули
COPY go.mod go.sum ./

# Принудительно пропускаем верификацию
RUN go mod tidy -v

# Копируем исходный код
COPY . .

# Собираем приложение, пропуская ошибки зависимостей
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -ldflags="-s -w" -o manga-reader ./cmd/api/main.go

# Создаем минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Копируем сертификаты и локальные данные
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем собранный бинарник из первого этапа
COPY --from=builder /app/manga-reader .

# Создаем каталоги для данных (чтобы избежать проблем с правами)
RUN mkdir -p /app/data/images

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./manga-reader"]