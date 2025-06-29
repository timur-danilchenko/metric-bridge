include .env
export

# DSN для подключения к PostgreSQL
DB_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
GOOSE=goose -dir ./migrations postgres "$(DB_DSN)"

.PHONY: run migrate_up migrate_down dbshell drop_db create_db ensure_db compose_up compose_down prometheus open_prometheus

# ==== Приложение ====
# Запускает основной воркер (основной entrypoint приложения)
run:
	@go run ./cmd/worker

# ==== Docker ====
# Поднимает все сервисы: Kafka, Zookeeper, Postgres, Prometheus
compose_up:
	@docker-compose --env-file .env -f docker/docker-compose.yml up 

# Останавливает и удаляет все контейнеры, созданные docker-compose
compose_down:
	@docker-compose --env-file .env -f docker/docker-compose.yml down

# ==== Миграции ====
# Применяет все миграции к базе данных
migrate_up: ensure_db
	@$(GOOSE) up

# Откатывает последнюю миграцию
migrate_down: ensure_db
	@$(GOOSE) down

# ==== Работа с БД ====
# Открывает интерактивную оболочку PostgreSQL (psql)
dbshell: ensure_db
	@psql "$(DB_DSN)"

# Удаляет базу данных (если существует)
drop_db:
	@psql "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres?sslmode=$(DB_SSLMODE)" \
		-c "DROP DATABASE IF EXISTS $(DB_NAME);"

# Создаёт базу данных (если не существует)
create_db:
	@psql "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres?sslmode=$(DB_SSLMODE)" \
		-c "CREATE DATABASE $(DB_NAME);"

# Проверяет наличие базы и создаёт её, если она не существует
ensure_db:
	@psql "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres?sslmode=$(DB_SSLMODE)" \
		-tAc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1 \
		|| (echo "⏳ Database '$(DB_NAME)' not found. Creating..."; make create_db)

# ==== Работа с Prometheus ====
# Поднимает только сервис Prometheus
prometheus:
	docker-compose --env-file .env -f docker/docker-compose.yml up -d prometheus

# Открывает интерфейс Prometheus в браузере
open_prometheus:
	open http://localhost:9090

# ==== Тестирование ====
# Запуск всех модульных тестов
test:
	go test ./... -v

# Запуск только модулей processor
test_processor:
	go test ./internal/processor -v

# (в будущем) Интеграционные тесты
test_integration:
	go test ./internal/tests -tags=integration -v
