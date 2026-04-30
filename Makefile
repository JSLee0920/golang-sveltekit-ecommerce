.PHONY: dev-backend dev-frontend build test migrate-up migrate-down migrate-status sqlc docker-up docker-down docker-start docker-stop

include .env
export

# dev
dev-backend:
	air -c .air.toml

dev-frontend:
	cd web && bun dev

# build
build:
	go build -o bin/server ./cmd/server

# test
test:
	go test ./... -v

# Database
migrate-up:
	goose -dir internal/db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir internal/db/migrations postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir internal/db/migrations postgres "$(DATABASE_URL)" status

# ── sqlc ──────────────────────────────────────
sqlc:
	sqlc generate

# ── Docker ────────────────────────────────────
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-start:
	docker compose start

docker-stop:
	docker compose stop
