# Dev Go APIs

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
| `ACCESS_TOKEN_TTL`     | Access token's time to live  | `Development`, `Docker` | Default valid time units: `ns, us (or µs), ms, s, m, h` and alternative time units: `d` for `day`, `w` for `week`, `mth` for `month` and `y` for `year` | 15m (15 minutes)                                                                         |
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

```
services:
  postgres:
    image: bitnami/postgresql:17.6.0
    container_name: postgres_db
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-postgres}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-postgres}
      POSTGRES_DB: ${POSTGRES_DB:-main}
    ports:
      - ${POSTGRES_PORT:-8002}:5432
    volumes:
      - pg-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d main"]
      interval: 5s
      timeout: 3s
      retries: 3

  redis:
    image: redis:alpine3.20
    container_name: redis_db
    command: ["redis-server", "--requirepass", "$$REDIS_PASSWORD"]
    environment:
      REDIS_PASSWORD: ${REDIS_PASSWORD:-redis123}
    ports:
      - ${REDIS_PORT:-8003}:6379
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "$$REDIS_PASSWORD", "ping"]
      interval: 5s
      timeout: 3s
      retries: 3

  api:
    image: phamgiaphuc/dev-go-apis:1.0.0
    container_name: apis
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      DATABASE_URL: ${DATABASE_URL}
      REDIS_URL: ${REDIS_URL}
      PORT: ${PORT}
      SERVER_URL: ${SERVER_URL}
      API_KEY: ${API_KEY}
      ACCESS_TOKEN_SECRET: ${ACCESS_TOKEN_SECRET}
      REFRESH_TOKEN_SECRET: ${REFRESH_TOKEN_SECRET}
      ACCESS_TOKEN_TTL: ${ACCESS_TOKEN_TTL}
      REFRESH_TOKEN_TTL: ${REFRESH_TOKEN_TTL}
      MIGRATION_MODE: ${MIGRATION_MODE}
    ports:
      - ${PORT:-8001}:${PORT:-8001}
  dbgate:
    image: dbgate/dbgate:alpine
    container_name: dbgate
    ports:
      - "8004:3000"
    volumes:
      - dbgate-data:/root/.dbgate
    environment:
      LOGIN: ${DBGATE_USER}
      PASSWORD: ${DBGATE_PASSWORD}

      CONNECTIONS: con1,con2

      LABEL_con1: postgres_db
      SERVER_con1: postgres
      USER_con1: postgres
      PASSWORD_con1: postgres
      PORT_con1: 5432
      ENGINE_con1: postgres@dbgate-plugin-postgres

      LABEL_con2: redis_db
      URL_con2: redis://default:redis123@redis:6379/0
      ENGINE_con2: redis@dbgate-plugin-redis

volumes:
  pg-data:
  redis-data:
  dbgate-data:

networks:
  default:
    driver: bridge

```
