include .env

run-dev:
	go run ./cmd/kbox-api/main.go -config .env

build:
	go build -o api ./cmd/kbox-api/main.go

clean:
	rm -f api

run: build
	./api -config .env

migrate:
	@source .env && goose -dir ./database/migrations postgres "$${DATABASE_URL}" up

create-migration:
	goose -dir ./database/migrations create example_migration sql

swagger:
	swag init -g cmd/kbox-api/main.go -o ./docs
