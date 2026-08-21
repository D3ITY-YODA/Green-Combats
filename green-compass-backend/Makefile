GO ?= go
APP_NAME ?= green-compass-api
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: run test lint fmt vet tidy build dev-up dev-down docker-build clean

run:
	$(GO) run ./cmd/api

test:
	$(GO) test -race -count=1 ./...

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
