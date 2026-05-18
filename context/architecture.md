# Backend Architecture Context

## Точка входа

- `cmd/main.go` - старт приложения.
- `internal/app/app.go` - сборка зависимостей приложения.
- `internal/app/server.go` - HTTP-сервер.

## Слои

### Config

`internal/config/config.go`

Загружает YAML-конфиг и переменные окружения. Собирает DSN для PostgreSQL.

### Application

`internal/application`

Сервисный слой. Здесь находится бизнес-логика для:

- auth;
- pages;
- events;
- plans;
- honor;
- achievements;
- founder;
- images.

Главный объект: `Service` в `service.go`.

### Domain

`internal/domain`

Содержит:

- DTO для HTTP-запросов и ответов;
- DB entities;
- интерфейсы репозиториев;
- доменные ошибки.

### Repository

`internal/repository/db`

SQL-реализация репозитория. Обычно один файл на сущность.

### Infrastructure

`internal/infrastructure/repository/postgres.go`

Подключение к PostgreSQL и создание DB-репозитория.

### HTTP transport

`internal/transport/http`

Gin handlers, router, middleware, upload.

`router.go` регистрирует:

- публичные маршруты;
- `/api/auth/login`;
- защищённую группу `/api/admin`.

## Общий поток запроса

HTTP handler -> application service -> repository interface -> PostgreSQL -> service -> handler response.

Ошибки мапятся в HTTP-ответы через `error_mapper.go`.

## Загрузка файлов

Файлы загружаются через `POST /api/admin/upload`.

Папка загрузок берётся из `upload.dir`. Роутер автоматически создаёт директорию, если её нет, и раздаёт её по `upload.url_prefix`.
