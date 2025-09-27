# Dev Go APIs

## Diagrams

- Databse Design: [Diagram Link](https://dbdiagram.io/d/Dev-Go-APIs-68d77773d2b621e42226cab2)

## Commands

- Goose commands:
  `goose -dir ./internal/database/migration create add_extensions sql`
- Docker commands:
  `docker build -t dev-go-apis:<version> .`
  `docker compose --env-file .env.docker up -d`
