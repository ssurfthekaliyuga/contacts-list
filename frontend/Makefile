.PHONY: up down rebuild bash start

CONTAINER_NAME=contacts-list-frontend-dev

up:
	@docker compose up -d

down:
	@docker compose down

rebuild:
	@docker compose build --no-cache

bash:
	@docker exec -it $(CONTAINER_NAME) bash

start:
	@make up && make bash