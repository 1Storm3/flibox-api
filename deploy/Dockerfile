FROM golang:1.23.1-alpine AS build

WORKDIR /app

COPY ../go.mod go.sum internal ./

RUN go mod download

COPY .. .

RUN go build -o main ./cmd/flibox-api/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]
