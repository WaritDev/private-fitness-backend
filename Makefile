start:
	go run cmd/app/main.go

gen:
	sqlc generate
	wire ./internal/...

gen-sqlc:
	sqlc generate

gen-wire:
	wire ./internal/...

migrate-schema:
	go run cmd/migration/main.go --migrate:schema

migrate-up:
	go run cmd/migration/main.go --migrate:up

migrate-down:
	go run cmd/migration/main.go --migrate:down --step=$(step)

migrate-reset:
	go run cmd/migration/main.go --migrate:reset

migrate-make:
	go run cmd/migration/main.go --migrate:make --name=$(name)

docker-migrate-schema:
	docker compose exec api go run cmd/migration/main.go --migrate:schema

docker-migrate-up:
	docker compose exec api go run cmd/migration/main.go --migrate:up

docker-migrate-down:
	docker compose exec api go run cmd/migration/main.go --migrate:down --step=$(step)

docker-migrate-reset:
	docker compose exec api go run cmd/migration/main.go --migrate:reset

docker-migrate-make:
	docker compose exec api go run cmd/migration/main.go --migrate:make --name=$(name)