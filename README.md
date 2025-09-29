# Dev Go APIs

## Table of Contents

- [Links and Diagrams](#links-and-diagrams-)
- [Commands](#commands-)
- [Environments](#environments-)

## Links and Diagrams 📈

- Swagger docs: Running on `localhost:{PORT}/docs` or `{domain}/docs`
- Environment sample file (.env.example): [File](./.env.example)
- Databse design: [Diagram link](https://dbdiagram.io/d/Dev-Go-APIs-68d77773d2b621e42226cab2)

## Commands 💻

- Make commands:
  - `make dev`: Run development
  - `make swag`: Generate Swagger docs
  - `make goose-up`: DB migration up
  - `make goose-down`: DB migration down
  - `make goose-down-to name=<version>`: DB migration down to a specific version
  - `make goose-create`: Create a migration sql file
- Goose commands:
  `goose -dir ./internal/database/migration create add_extensions sql`
- Docker commands:
  `docker build -t phamgiaphuc/dev-go-apis:<version> .`
  `docker compose --env-file .env.docker up -d`

## Environments 🔐

- **Server variables**:

| Name                   | Description                  | Phase                   | Note                                                                                                                                                     | Default value                                                                        |
| ---------------------- | ---------------------------- | ----------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| `PORT`                 | Server's port                | `Development`, `Docker` | Ex: `8000`                                                                                                                                               | 8000                                                                                 |
| `SERVER_URL`           | Server's url                 | `Development`, `Docker` | Ex: `http://localhost:8000`                                                                                                                              | http://localhost:8000                                                                |
| `MIGRATION_MODE`       | Run database's migrations    | `Docker`                | `0` for `off`, `1` for `on`                                                                                                                              | 0                                                                                    |
| `ACCESS_TOKEN_SECRET`  | Access token's secret        | `Development`, `Docker` | -                                                                                                                                                        | @scecret123                                                                          |
| `REFRESH_TOKEN_SECRET` | Refresh token's secret       | `Development`, `Docker` | -                                                                                                                                                        | @scecret123                                                                          |
| `ACCESS_TOKEN_TTL`     | Access token's time to live  | `Development`, `Docker` | - Default valid time units: `ns, us (or µs), ms, s, m, h` and alternative time units: `d` for `day`, `w` for `week` `mth` for `month` and `y` for `year` | 15m (15 minutes)                                                                     |
| `REFRESH_TOKEN_TTL`    | Refresh token's time to live | `Development`, `Docker` | - Default valid time units: `ns, us (or µs), ms, s, m, h` and alternative time units: `d` for `day`, `w` for `week` `mth` for `month` and `y` for `year` | 7d (7 days)                                                                          |
| `DATABASE_URL`         | Postgres DB connection url   | `Development`, `Docker` | -                                                                                                                                                        | postgres://{user}:{pass}@localhost:5432/{db_name}?sslmode=disable&search_path=public |
| `REDIS_URL`            | Redis DB connection url      | `Development`, `Docker` | -                                                                                                                                                        | redis://{user}:{pass}@localhost:6379/0                                               |

- **Postgres variables** (Docker deployment)

| Name                | Description       | Default value |
| ------------------- | ----------------- | ------------- |
| `POSTGRES_USER`     | Postgres username | postgres      |
| `POSTGRES_PASSWORD` | Postgres password | postgres      |
| `POSTGRES_DB`       | Database name     | main          |
| `POSTGRES_PORT`     | Postgres port     | 5432          |

- **Redis variables** (Docker deployment)

| Name             | Description    | Default value |
| ---------------- | -------------- | ------------- |
| `REDIS_PASSWORD` | Redis password | redis123      |
| `REDIS_PORT`     | Redis port     | 6379          |

- **DBGate variables** (DBGate deployment)

| Name              | Description         | Default value |
| ----------------- | ------------------- | ------------- |
| `DBGATE_USER`     | DBGate web user     | admin         |
| `DBGATE_PASSWORD` | DBGate web password | admin123      |

## Features checklist

- **Auth**:

  - [x] Log in
  - [x] Register

- **Role**:

  - [x] Get a role list
  - [x] Create a role
  - [x] Update a role with assigning permissions to that role
  - [x] Get a role
  - [x] Delete a role
