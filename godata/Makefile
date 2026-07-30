.PHONY: build run test lint clean docker-build docker-up docker-down migrate

APP_NAME := phoenix-api
GO := go
GOFMT := gofmt

build:
	$(GO) build -o bin/$(APP_NAME) ./cmd/api

run:
	$(GO) run ./cmd/api

test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

test-short:
	$(GO) test -short ./...

lint:
	golangci-lint run ./...

fmt:
	$(GOFMT) -s -w .

vet:
	$(GO) vet ./...

clean:
	rm -rf bin/ coverage.out

docker-build:
	docker build -t phoenix-api .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

migrate-up:
	migrate -path migrations -database "postgres://phoenix:phoenix@localhost:5432/phoenix?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://phoenix:phoenix@localhost:5432/phoenix?sslmode=disable" down

mod:
	$(GO) mod tidy

all: fmt vet test build
