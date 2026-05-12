.PHONY: help dev run build test clean proto deps
# установить сначало brew install make && brew install golang-migrate
#Go: Restart Language Server все желтым

ifneq (,$(wildcard .env.dev))
    include .env.dev
    export
endif

MIGRATIONS_DIR := ./migrations
ifeq ($(strip $(STORE_SERVICE_DATABASE_URL)),)
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
else
DATABASE_URL := $(STORE_SERVICE_DATABASE_URL)
endif

help:
	@echo "Доступные команды:"
	@echo "  make dev       - Запуск в dev режиме с hot reload (рекомендуется)"
	@echo "  make run       - Запуск сервиса без hot reload"
	@echo "  make build     - Сборка бинарника"
	@echo "  make test      - Запуск тестов"
	@echo "  make proto     - Генерация proto (из shared-protos)"
	@echo "  make deps      - Установка зависимостей"
	@echo "  make clean     - Очистка"
	@echo "  make install   - Установка зависимостей"
	@echo "  make migrate-up   - Применение миграций"
	@echo "  make migrate-down - Откат последней миграции"

install:
	#go mod init gitlab.kalina-tech.ru/backend/store-service
	go install github.com/air-verse/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go get github.com/IBM/sarama@latest
	go get github.com/google/uuid@latest
	go get github.com/lib/pq@latest
	go get github.com/redis/go-redis/v9@latest
	go get github.com/spf13/viper@latest
	go get golang.org/x/crypto@latest
	go get google.golang.org/grpc@latest
	go get github.com/kalina-malina/IM-PROTOS@latest
	go get github.com/go-playground/validator/v10@latest


dev:
	@echo "Сервис запущен в dev режиме с hot reload"
	@PATH="$$PATH:$$HOME/go/bin" air

run:
	@echo "Сервис запущен в режиме без hot reload"
	@go run ./cmd/main.go

build:
	@echo "Сборка бинарника..."
	go build -o bin/auth-service cmd/main.go
	@echo "Сервис собран в бинарный файл"

test:
	@echo "Запуск тестов..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Тесты завершены"

deps:
	go mod download
	go mod tidy


update-protos:
	@echo "Обновление зависимости store_service"
	@go get gitlab.kalina-tech.ru/protos/store_service@latest
	@go mod download gitlab.kalina-tech.ru/protos/store_service
	@echo "Обновление зависимости set_service"
	@go get gitlab.kalina-tech.ru/protos/set_service@latest
	@go mod download gitlab.kalina-tech.ru/protos/set_service
	@echo "Зависимости proto обновлены"

clean:
	rm -rf bin/ tmp/
	rm -f coverage.out



migrate-create:
	@echo "Создание новой миграции..."
	@echo "MIGRATIONS_DIR: $(MIGRATIONS_DIR)"
	@PATH="$$PATH:$$HOME/go/bin" migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)


migrate-up:
	@echo "STORE_SERVICE_DATABASE_URL: $(STORE_SERVICE_DATABASE_URL)"
	@echo "Применяем миграции..."
	@echo "MIGRATIONS_DIR: $(MIGRATIONS_DIR)"
	@echo "DATABASE_URL: $(DATABASE_URL)"
	@PATH="$$PATH:$$HOME/go/bin" migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up || (echo "migrate не найден. Выполни: make install (или: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)" && exit 1)

migrate-down:
	@echo "Откатываем последнюю миграцию..."
	@echo "MIGRATIONS_DIR: $(MIGRATIONS_DIR)"
	@echo "DATABASE_URL: $(DATABASE_URL)"
	@PATH="$$PATH:$$HOME/go/bin" migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down || (echo "migrate не найден. Выполни: make install (или: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)" && exit 1)

migrate-force:
	@echo "Принудительно устанавливаем версию миграции..."
	@echo "DATABASE_URL: $(DATABASE_URL)"
	@PATH="$$PATH:$$HOME/go/bin" migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(VERSION)


docker:
	docker build --platform linux/amd64 -t kalinamalinadev/store-service:v9 .
	docker push kalinamalinadev/store-service:v9