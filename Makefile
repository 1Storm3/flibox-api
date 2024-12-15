include .env

PROTO_DIR=api/proto
OUT_DIR=pkg/proto

proto:
	protoc --proto_path=$(PROTO_DIR) \
	       --go_out=$(OUT_DIR) \
	       --go-grpc_out=$(OUT_DIR) \
	       $(PROTO_DIR)/*.proto

run-dev:
	go run ./cmd/flibox-api/main.go -config .env

build:
	go build -o api ./cmd/flibox-api/main.go

clean:
	rm -f api

run: build
	./api -config .env

migrate:
	@source .env && goose -dir ./database/migrations postgres "$${DATABASE_URL}" up
down:
	@source .env && goose -dir ./database/migrations postgres "$${DATABASE_URL}" down

create-migration:
	goose -dir ./database/migrations create example_migration sql

swagger:
	swag init -g cmd/flibox-api/main.go -o ./docs

test:
	go test -v ./...

testauth:
	go test -cover -coverprofile=coverage.out ./internal/modules/auth

testTool:
	go tool cover -func=coverage.out

testHtml:
	go tool cover -html=coverage.out
