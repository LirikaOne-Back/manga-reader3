# Manga Reader

Современное приложение для чтения манги с красивым интерфейсом, разработанное с использованием Go (backend) и React (frontend).

## Особенности

- 📱 Адаптивный дизайн для мобильных устройств и десктопов
- 🌙 Темная и светлая тема
- 📚 Удобная навигация по главам
- 🔍 Поиск и фильтрация манги
- 📑 Закладки и история чтения
- 🖼️ Два режима чтения: вертикальный (прокрутка) и постраничный
- 🔐 Авторизация и регистрация пользователей

## Технологический стек

### Backend
- Go (Golang)
- Gin - HTTP web framework
- log/slog - современный логгер для Go
- PostgreSQL - база данных
- JWT - аутентификация
- Docker - контейнеризация

### Frontend
- React
- TypeScript
- Redux Toolkit
- React Router v6
- Tailwind CSS
- Headless UI
- Axios
- Docker

## Требования

- Docker и Docker Compose
- Git
- Свободный порт 80 (для фронтенда) и 8080 (для API)

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/manga-reader.git
cd manga-reader
```

2. Создайте `.env` файл с необходимыми переменными окружения:
```bash
# Настройки PostgreSQL
POSTGRES_USER=manga_user
POSTGRES_PASSWORD=manga_password
POSTGRES_DB=manga_reader

# Настройки JWT
JWT_SECRET=your_super_secret_key_change_in_production
```

3. Запустите приложение с помощью Docker Compose:
```bash
docker-compose up -d
```

4. Приложение будет доступно по адресу:
    - Frontend: http://localhost
    - Backend API: http://localhost:8080/api

## Структура проекта

```
manga-reader/
├── backend/               # Go бэкенд
│   ├── cmd/               # Точка входа
│   ├── internal/          # Внутренние пакеты
│   ├── pkg/               # Общие пакеты
│   └── Dockerfile         # Dockerfile для бэкенда
│
├── frontend/              # React фронтенд
│   ├── public/            # Статические файлы
│   ├── src/               # Исходный код
│   ├── Dockerfile         # Dockerfile для фронтенда
│   └── nginx.conf         # Конфигурация Nginx
│
├── docker-compose.yml     # Конфигурация Docker Compose
└── README.md              # Документация
```

## Разработка

### Backend

Для локальной разработки без Docker:

```bash
cd backend
go run cmd/api/main.go
```

### Frontend

Для локальной разработки без Docker:

```bash
cd frontend
npm install
npm start
```

## Лицензия

MIT