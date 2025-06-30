-include .env
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

# Поднимает все сервисы для тестов: Kafka, Zookeeper, Postgres, Prometheus
compose_test_up:
	@docker-compose --env-file .env -f docker/docker-compose.test.yml up

# Останавливает и удаляет все тестовые контейнеры, созданные docker-compose
compose_test_down:
	@docker-compose --env-file .env -f docker/docker-compose.test.yml down

# Поднимает все сервисы для CI/CD: Kafka, Zookeeper, Postgres, Prometheus
compose_ci:
	@docker-compose --env-file .env -f docker/docker-compose.ci.yml up -d

# Останавливает и удаляет все контейнеры, созданные docker-compose для CI/CD
compose_ci_down:
	@docker-compose --env-file .env -f docker/docker-compose.ci.yml down

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
	@docker-compose --env-file .env -f docker/docker-compose.yml up -d prometheus

# Открывает интерфейс Prometheus в браузере
open_prometheus:
	@open http://localhost:9090

# ==== Тестирование ====
# Запуск всех модульных тестов
test:
	@go test ./... -v

# Запуск только модулей processor
test_processor:
	@go test -tags=unit-test-processor ./internal/processor -v

# Интеграционные тесты
test_integration:
	@DB_INTEGRATION_DSN=$(DB_DSN) go test -tags=integration ./internal/tests -v

	# go test -tags=integration ./internal/tests -v

# ==== CI Тестирование ====

# Запускает окружение CI и выполняет интеграционные тесты
test_ci: compose_ci
	@echo "Waiting for services to be ready..."
	@sleep 5

	@echo "Running database migrations..."
	@$(MAKE) migrate_up

	@sleep 3
	@echo "Running integration tests..."
	@DB_INTEGRATION_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE) \
		go test -tags=integration ./internal/tests -v

	@echo "🧹 Shutting down CI containers..."
	@$(MAKE) compose_ci_down
