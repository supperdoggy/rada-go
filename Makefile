.PHONY: run build test lint fmt

RUN_TARGET ?= ./cmd/...
GOLANGCI_LINT ?= golangci-lint

run:
	go run $(RUN_TARGET)

build:
	go build ./...

test:
	go test ./...

lint:
	$(GOLANGCI_LINT) run

fmt:
	go fmt ./...
