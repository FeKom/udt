.PHONY: up down migrate-up migrate-down migrate-new run build test vet tidy

up:            ## Start Postgres and apply migrations
	docker compose up -d --wait postgres
	docker compose run --rm migrate up

down:          ## Stop containers (data volume is kept)
	docker compose down

migrate-up:
	docker compose run --rm migrate up

migrate-down:  ## Roll back the last migration
	docker compose run --rm migrate down 1

migrate-new:   ## make migrate-new name=create_families
	docker compose run --rm migrate create -ext sql -dir /migrations -seq $(name)

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy
