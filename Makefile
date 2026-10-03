APP_NAME ?= simple-api
BIN      ?= bin/server
PKG      ?= ./...

.PHONY: all build test vet run docker clean

all: vet test build

## build: compile the server binary into ./bin
build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BIN) ./cmd/server

## test: run the test suite
test:
	go test -race $(PKG)

## vet: run go vet
vet:
	go vet $(PKG)

## run: start the server locally (PORT, LOG_LEVEL from the environment)
run:
	go run ./cmd/server

## docker: build the production image
docker:
	docker build -t $(APP_NAME):latest .

## clean: remove build artifacts
clean:
	rm -rf bin
