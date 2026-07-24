include .env
export

export PROJECT_ROOT=$(shell pwd)

up-env:
	docker compose up -d db

down-env:
	docker compose down db

create-migrate:
	@if [ -z "$(seq)" ]; then \
		echo "Error, miss seq. Example: make create-migrate seq=init" && \
		exit 1; \
	fi;

	docker compose run --rm migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-command:
	docker compose run --rm migrate \
		-path /migrations \
		-database postgres://${P_USER}:${P_PASSWORD}@db:5432/${P_DB}?sslmode=disable \
		$(command)

migrate-up:
	make migrate-command command=up

migrate-down:
	make migrate-command command=down

app-run-tgbot:
	go mod tidy -e && \
	go run ./cmd/tgbot-golang
