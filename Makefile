include .env

run-dev:
	go run ./cmd/kinopoisk-api/main.go -config .env

build:
	go build -o api ./cmd/kinopoisk-api/main.go

clean:
	rm -f api

run: build
	./api -config .env

migrate:
	goose -dir ./database/migrations postgres "postgresql://admin:admin@185.209.20.129:5432/postgres?sslmode=disable" up
