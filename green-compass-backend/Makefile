GO ?= go
APP_NAME ?= green-compass-api
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

ifneq (,$(wildcard .env))
include .env
export
endif

GC_DEV_PG_PORT ?= 5432
GC_TEST_DATABASE_URL ?= postgres://greencompass:greencompass@127.0.0.1:$(GC_DEV_PG_PORT)/greencompass?sslmode=disable

.PHONY: run test lint fmt vet tidy build test-integration ci-test migrate-up migrate-down migrate-status dev-up dev-down docker-build clean

run:
	$(GO) run ./cmd/api

test:
	$(GO) test -race -count=1 ./...

test-integration:
	GC_TEST_DATABASE_URL=$(GC_TEST_DATABASE_URL) $(GO) test -race -count=1 -run Integration ./...

ci-test: dev-up
	GC_TEST_DATABASE_URL=$(GC_TEST_DATABASE_URL) $(GO) test -race -count=1 ./...

MIGRATE_DSN ?= postgres://greencompass:greencompass@127.0.0.1:$(GC_DEV_PG_PORT)/greencompass?sslmode=disable
MIGRATE_STEPS ?= 0

migrate-up:
	$(GO) run ./cmd/migrate -dsn $(MIGRATE_DSN) -steps $(MIGRATE_STEPS) up

migrate-down:
	$(GO) run ./cmd/migrate -dsn $(MIGRATE_DSN) -steps $(MIGRATE_STEPS) down

migrate-status:
	$(GO) run ./cmd/migrate -dsn $(MIGRATE_DSN) status

lint:
	$(GO) tool golangci-lint run

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

build:
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o bin/api ./cmd/api

dev-up:
	bash scripts/dev-up.sh

dev-down:
	bash scripts/dev-down.sh

docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(APP_NAME):$(VERSION) .

clean:
	rm -rf bin
