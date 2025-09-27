# Dev Go APIs

## Links and Diagrams

- Swagger docs: Running on `localhost:{PORT}/docs` or `{domain}/docs`
- Environment sample file: [File](./.env.example)
- Databse design: [Diagram Link](https://dbdiagram.io/d/Dev-Go-APIs-68d77773d2b621e42226cab2)

## Commands

- Goose commands:
  `goose -dir ./internal/database/migration create add_extensions sql`
- Docker commands:
  `docker build -t dev-go-apis:<version> .`
  `docker compose --env-file .env.docker up -d`
