include .env
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
GOOSE=goose -dir ./migrations postgres "$(DB_URL)"

.PHONY: run migrate_up migrate_down dbshell drop_db create_db ensure_db

run:
	go run ./cmd/worker

migrate_up: ensure_db
	@$(GOOSE) up

migrate_down: ensure_db
	@$(GOOSE) down

dbshell: ensure_db
	@ psql "$(DB_URL)"

drop_db:
	@psql "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres?sslmode=$(DB_SSLMODE)" \
		-c "DROP DATABASE IF EXISTS $(DB_NAME);"

create_db:
	@psql "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres?sslmode=$(DB_SSLMODE)" \
		-c "CREATE DATABASE $(DB_NAME);"

# Проверка наличия БД, если нет — создаёт
ensure_db:
	@psql "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres?sslmode=$(DB_SSLMODE)" \
		-tAc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1 \
		|| (echo "Database '$(DB_NAME)' not found. Creating..."; make create_db)
