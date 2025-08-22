# 🔗 URL Shortener (Golang + Postgres + Redis + Gin)
Простой сервис сокращения ссылок, написанный на Go.  
Хранит данные в **Postgres**, кэширует обращения через **Redis**,  
использует **Gin** для HTTP-обработки и поддерживает Swagger-документацию.
---
## 📑 Оглавление
- [✨ Возможности](#-возможности)
- [⚙️ Стек](#️-стек)
- [📂 Структура проекта](#-структура-проекта)
- [📦 Запуск](#-запуск)
    - [1. Клонировать проект](#1-клонировать-проект)
    - [2. Создать файл env](#2-создать-файл-env)
    - [3. Поднять контейнеры](#3-поднять-контейнеры)
- [🚀 API](#-api)
    - [POST /shorten](#post-shorten)
    - [GET /short](#get-short)
- [📖 Swagger](#-swagger)
---

## ✨ Возможности
- Сокращение длинных URL до коротких кодов
- Редирект по коротким ссылкам
- Валидация и нормализация URL
- Кэширование в Redis с TTL
- Идемпотентность: один и тот же URL всегда даёт одинаковый short-код
- Rate limiting (ограничение запросов, можно включить через middleware)
- TTL для ссылок (опционально)
- Swagger-документация (`/swagger/index.html`)

---

## ⚙️ Стек
- **Go** 1.24.4
- **Postgres** — основное хранилище
- **Redis** — кэш
- **Gin** — HTTP-сервер
- **golang-migrate** — миграции базы
- **swaggo/swag** — Swagger-документация

---

## 📂 Структура проекта
- cmd/app/main.go # точка входа
- internal/config # конфигурация и env
- internal/database # подключение к Postgres/Redis
- internal/repository # работа с БД и кэшем
- internal/usecase # бизнес-логика
- internal/handler # HTTP-хэндлеры (Gin)
- migrations/ # SQL-миграции
- internal/docs/ # Swagger (генерируется swag init)

## 📦 Запуск

### 1. Клонировать проект
```bash
git clone https://github.com/iviv660/url-shortener.git
cd url-shortener
```

### 2. Создать файл .env
```bash
# Postgres
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=

# Postgres (строка подключения)
DATABASE_URL=postgresql://_:_@postgres:5432/_?sslmode=disable

# Redis (без пароля)
REDIS_URL=redis://redis:6379/0

# Базовый URL приложения
BASE_URL=http://localhost:3000

# Секрет для генерации коротких кодов
SHORT_CODE_SECRET=_

# TTL кэша (секунды)
CACHE_TTL_SECONDS=_

# Окружение
APP_ENV=_
```

### 3. Поднять контейнеры
```bash
docker-compose up
```

## 🚀 API
### POST /shorten
#### Сократить длинный URL.

### Запрос
```bash
curl -X POST http://localhost:3000/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/iviv660"}'
```
### Ответ
```bash
{
  "short_code": "iHwL01z-",
  "short_url": "http://localhost:3000/сокращеный url"
}
```

### GET /:short
#### Редирект на оригинальный URL.

```bash
curl -i http://localhost:3000/сокращеный url
```
### Ответ
```bash
{
  "short_code": "сокращеный url",
  "short_url": "http://localhost:3000/сокращеный url"
}
```
## 📖 Swagger
#### После запуска доступно по адресу:
```bash
    http://localhost:3000/swagger/index.html
```