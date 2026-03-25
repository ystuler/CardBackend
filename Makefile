.PHONY: help run test swagger db-up db-down db-logs tidy fmt

help:
	@echo Available targets:
	@echo   make run      - start backend
	@echo   make test     - run all tests
	@echo   make swagger  - generate Swagger docs from annotations
	@echo   make db-up    - start local PostgreSQL in Docker
	@echo   make db-down  - stop local Docker services
	@echo   make db-logs  - show PostgreSQL logs
	@echo   make tidy     - run go mod tidy
	@echo   make fmt      - run go fmt on all packages

run:
	go run ./cmd/main.go

test:
	go test ./...

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.8.1 init -g cmd/main.go -o docs --parseInternal

db-up:
	docker compose up -d db

db-down:
	docker compose down

db-logs:
	docker compose logs db

tidy:
	go mod tidy

fmt:
	go fmt ./...
