<div style="text-align: center; margin-top: 50px;">
  <h1 style="font-size: 3rem; margin-bottom: 0.5rem;">
    Dev Go APIs
  </h1>
  <h3 style="font-weight: normal; color: #d4d4d4;">
    This is a useless API project built for learning purposes.
  </h3>
  <img src="./screeshot.png" alt="Screenshot" style="margin-top: 1rem; max-width: 100%; height: auto;">
</div>

## Table of Contents

- [Links and Diagrams](#links-and-diagrams-)
- [Commands](#commands-)
- [Environments](#environments-)
- [Features checklist](#features-checklist)

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

| Name                   | Description                  | Phase                   | Note                                                                                                                                                    | Default value                                                                            |
| ---------------------- | ---------------------------- | ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `PORT`                 | Server's port                | `Development`, `Docker` | Ex: `8000`                                                                                                                                              | 8000                                                                                     |
| `SERVER_URL`           | Server's url                 | `Development`, `Docker` | Ex: `http://localhost:8000`                                                                                                                             | http://localhost:8000                                                                    |
| `MIGRATION_MODE`       | Run database's migrations    | `Docker`                | `0` for `off`, `1` for `on`                                                                                                                             | 0                                                                                        |
| `CORS_ALLOWED_ORIGINS` | CORs allowed origins         | `Development`, `Docker` | -                                                                                                                                                       | localhost:3000, localhost:5173                                                           |
| `API_KEY`              | Key for accessing the apis   | `Development`, `Docker` | The API key middleware is not enabled, if `API_KEY` not set.                                                                                            | -                                                                                        |
| `HMAC_SECRET_KEY`      | Key for Hmac signature       | `Development`, `Docker` | -                                                                                                                                                       | @secret123                                                                               |
| `ACCESS_TOKEN_SECRET`  | Access token's secret        | `Development`, `Docker` | -                                                                                                                                                       | @secret123                                                                               |
| `REFRESH_TOKEN_SECRET` | Refresh token's secret       | `Development`, `Docker` | -                                                                                                                                                       | @secret123                                                                               |
| `ACCESS_TOKEN_TTL`     | Access token's time to live  | `Development`, `Docker` | Default valid time units: `ns, us (or µs), ms, s, m, h` and alternative time units: `d` for `day`, `w` for `week`, `mth` for `month` and `y` for `year` | 30m (30 minutes)                                                                         |
| `REFRESH_TOKEN_TTL`    | Refresh token's time to live | `Development`, `Docker` | Default valid time units: `ns, us (or µs), ms, s, m, h` and alternative time units: `d` for `day`, `w` for `week`, `mth` for `month` and `y` for `year` | 7d (7 days)                                                                              |
| `DATABASE_URL`         | Postgres DB connection url   | `Development`, `Docker` | -                                                                                                                                                       | postgres://{user}:{password}@localhost:5432/{db_name}?sslmode=disable&search_path=public |
| `REDIS_URL`            | Redis DB connection url      | `Development`, `Docker` | -                                                                                                                                                       | redis://{user}:{password}@localhost:6379/0                                               |
| `GOOGLE_CLIENT_ID`     | Google Client ID             | `Development`, `Docker` | -                                                                                                                                                       | -                                                                                        |
| `GOOGLE_CLIENT_SECRET` | Google Client Secret         | `Development`, `Docker` | -                                                                                                                                                       | -                                                                                        |
| `GOOGLE_REDIRECT_URL`  | Google Redirect URL          | `Development`, `Docker` | -                                                                                                                                                       | -                                                                                        |

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

## Docker deployment:

- [Dockerfile](./Dockerfile)
- [docker-compose.yml](./docker-compose.yml)
