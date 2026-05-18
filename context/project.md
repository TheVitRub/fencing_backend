# Backend Project Context

## Назначение

Бэкенд обслуживает публичный сайт и админ-панель школы Ferrum et Gloria:

- отдаёт страницы и сущности сайта;
- принимает CRUD-запросы из админки;
- выпускает JWT для администратора;
- принимает загрузку изображений;
- раздаёт загруженные файлы через `/uploads`.

## Стек

- Go 1.22.
- Gin.
- PostgreSQL.
- sqlx + lib/pq.
- Viper для конфигурации.
- JWT `github.com/golang-jwt/jwt/v5`.
- bcrypt для проверки пароля администратора.
- zap через локальный пакет `pkg/logger`.

## Команды

```bash
go run ./cmd/main.go
go test ./...
```

В `Makefile` также есть цели:

```bash
make run
make test
make migrate-up
make migrate-down
```

Часть `Makefile` выглядит унаследованной от другого сервиса, поэтому перед использованием целей сборки и docker-публикации лучше проверять команды.

## Конфигурация

Основной dev-конфиг: `configs/dev.yaml`.

Ключевые параметры:

- HTTP: `server.host`, `server.port`.
- PostgreSQL: `database.*`.
- JWT: `auth.jwt_secret`, `auth.jwt_ttl`.
- Upload: `upload.dir`, `upload.url_prefix`, `upload.max_size_mb`.

Переменные окружения могут переопределять конфиг. JWT secret в dev обычно приходит из `.env.dev`.

## Администратор по умолчанию

Seed-миграция создаёт администратора:

- login: `admin`
- password: `admin123`

После первого входа пароль лучше заменить.
