include .env

GOOSE_MIGRATION_DIR ?= ./internal/database/migration


dev: # Run development
	@air -c .air.toml

swag: # Generate Swagger docs
	@swag fmt
	@swag init -g cmd/api/main.go -o ./docs

goose-up: # DB migration up
	@goose -dir=./internal/database/migration postgres $(DATABASE_URL) up

goose-down: # DB migration down
	@goose -dir=./internal/database/migration postgres $(DATABASE_URL) down

goose-down-to: # DB migration down to a specific version
	@if [ -z "$(version)" ]; then \
		echo "❌ Error: missing version. Usage: make goose-down-to version=20230910120000"; \
		exit 1; \
	fi
	@goose -dir=./internal/database/migration postgres $(DATABASE_URL) down-to $(version)

goose-create: # Create a migration sql file
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: missing name. Usage: make goose-create name=add_users_table"; \
		exit 1; \
	fi
	@goose -dir $(GOOSE_MIGRATION_DIR) create $(name) sql