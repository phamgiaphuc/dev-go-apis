include .env

GOOSE_MIGRATION_DIR ?= ./internal/database/migration


dev: # Run development
	@air -c .air.toml
swag: # Swagger Docs
	@swag fmt
	@swag init -g cmd/api/main.go -o ./docs
goose-up: # Migration up
	@goose -dir=./internal/database/migration postgres $(DATABASE_URL) up
goose-down: # Migration down
	@goose -dir=./internal/database/migration postgres $(DATABASE_URL) down
goose-down-to: # Migration down to a specific version
	@if [ -z "$(version)" ]; then \
		echo "❌ Error: missing version. Usage: make goose-down-to version=20230910120000"; \
		exit 1; \
	fi
	@goose -dir=./internal/database/migration postgres $(DATABASE_URL) down-to $(version)
goose-create: # Create a migration sql file: make goose-create name=create_users_table	
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: missing name. Usage: make migration name=add_users_table"; \
		exit 1; \
	fi
	@goose -dir $(GOOSE_MIGRATION_DIR) create $(name) sql