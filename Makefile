# TODO PHONY

## TODO not in CI ;)
include .env

.EXPORT_ALL_VARIABLES:

migrate-up:
	go run cmd/migrate/main.go up

run: migrate-up
	go run cmd/main/main.go

docker-jet:
	docker build -f Dockerfile.jet -t jet \
	--build-arg GO_VERSION=${GO_VERSION} \
	.

gen-models: docker-jet
	docker run --env-file .env --network host -v .:/work jet

gen-models2:
	jet -dsn=postgresql://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}:${DB_PORT}/${DB_NAME}?sslmode=disable -schema=public -path=./internal/db/.gen

start:
	docker-compose up

restart-clean:
	docker-compose down --volumes && docker-compose up