#!make

include .env
export $(shell sed 's/=.*//' .env)

.PHONY:run
run:
	@go run ./main.go

.PHONY:lint
lint-local:
	@golangci-lint run

.PHONY:lint
lint:
	@docker run --rm \
	  -v ${PWD}:/app \
	  -w /app \
	  golangci/golangci-lint:v1.44-alpine \
	  golangci-lint run -v --timeout 5m

.PHONY:test
test:
	@go test -v ./...

.PHONY:swagger
swagger:
	@swag init
	@swag fmt

.PHONY:migrate-up
migrate-up:
	@docker run --rm \
		-v $(shell pwd)/migrations:/migrations \
		--network host \
		migrate/migrate:4 \
		-path=/migrations/ \
		-database ${DB_POSTGRESQL_URI}?sslmode=disable \
		up

.PHONY:migrate-down
migrate-down:
	@docker run --rm \
		-v $(shell pwd)/migrations:/migrations \
		--network host \
		migrate/migrate:4 \
		-path=/migrations/ \
		-database ${DB_POSTGRESQL_URI}?sslmode=disable \
		down 1

gosec:
	@gosec -quiet -tests ./...
