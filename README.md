# Private Fitness Backend

Backend service using Clean Architecture, Fiber, sqlc, and MariaDB.

## Quick Start

1) Start services (MariaDB + API dev container):

```bash
docker compose up -d
```

2) Run migrations:

```bash
make docker-migrate-up
```

3) Generate code (sqlc + wire):

```bash
make gen
```

4) Start app locally (without Docker):

```bash
make start
```

API will be available at http://localhost:8000

## Notes
- SQL placeholders for MySQL/MariaDB must be `?` (not `$1`, `$2`, …).
- Makefile runs Wire via `go run`, so you don’t need a `wire` binary installed.
- Environment variables live in `.env`. Important keys:
  - `DB_DRIVER=mysql`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PARAMS`
  - `USE_DOCKER=true` when running in Docker
  - `JWT_SECRET` for signing access tokens

## Useful Commands
- `make gen-sqlc` — generate sqlc models/queries
- `make gen-wire` — generate DI code (Google Wire)
- `make migrate-make name=create_xxx` — create a new migration
- `make docker-migrate-up` — apply migrations in container
- `make docker-migrate-down step=1` — rollback last migration

## Project Docs
- See `docs/overview.md` for architecture overview
- See `docs/develop_101.md` for step-by-step feature implementation
