.PHONY: start

start:
	@docker compose up -d
	@go run cmd/main.go -env .env

down:
	@docker compose down

restart:
	@make down
	@make start

total-down:
	@docker compose down --volumes --remove-orphans

total-restart:
	@make total-down
	@make start

wire:
	@wire ./internal/wire