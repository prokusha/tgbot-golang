include .env
export

export PROJECT_ROOT=$(shell pwd)

up-env:
	docker compose up -d db

down-env:
	docker compose down db

migrate-up:
	make migrate-command command=up

migrate-down:
	make migrate-command command=down

app-run-tgbot:
	go mod tidy && \
	go run ./cmd/tgbot-golang
