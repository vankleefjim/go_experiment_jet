# TODO PHONY

ifndef CI
	include .env
endif

.EXPORT_ALL_VARIABLES:

migrate-up:
	go run cmd/migrate/main.go up

run: start gen-models
	HTTP_PORT=8081 go run cmd/apiserver/main.go

run-ui: gen-templ
	HTTP_PORT=8082 go run cmd/ui/main.go

docker-jet: migrate-up
	docker build -f Dockerfile.jet -t jet \
	--build-arg GO_VERSION=${GO_VERSION} \
	.

gen-models: docker-jet
	docker run --env-file .env --network host -v .:/work jet

gen-models2: migrate-up
	jet -dsn=postgresql://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}:${DB_PORT}/${DB_NAME}?sslmode=disable -schema=public -path=./internal/db/.gen

gen-templ:
	templ generate

start:
	${COMPOSE_CMD} up -d --wait

stop:
	${COMPOSE_CMD} down --volumes

restart-clean: stop run
	

test:
	go test -shuffle on -race ./...

lint: # TODO add this to .nix
	golangci-lint run