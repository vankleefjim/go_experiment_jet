# TODO PHONY

ifndef CI
	include .env
endif

.EXPORT_ALL_VARIABLES:

migrate-up:
	go run cmd/migrate/main.go up

run: start gen-models
	go run cmd/main/main.go

docker-jet: migrate-up
	docker build -f Dockerfile.jet -t jet \
	--build-arg GO_VERSION=${GO_VERSION} \
	.

gen-models: docker-jet
	docker run --name jet --env-file .env --network host -v .:/work jet

gen-models2: migrate-up
	jet -dsn=postgresql://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}:${DB_PORT}/${DB_NAME}?sslmode=disable -schema=public -path=./internal/db/.gen

start:
	${COMPOSE_CMD} up -d --wait

stop:
	${COMPOSE_CMD} down --volumes

restart-clean: stop run
	

test:
	go test -shuffle on -race ./...

lint: # TODO add this to .nix
	golangci-lint run