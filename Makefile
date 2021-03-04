SHELL := /bin/bash

.PHONY: all check format vet lint build test generate tidy integration_test

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check               to do static check"
	@echo "  build               to create bin directory and build"
	@echo "  generate            to generate code"
	@echo "  test                to run test"
	@echo "  integration_test    to run integration test"

check: vet

format:
	@echo "go fmt"
	@go fmt ./...
	@echo "ok"

vet:
	@echo "go vet"
	@go vet ./...
	@echo "ok"

generate:
	@echo "generate code"
	@go generate -tags tools ./...
	@go fmt ./...
	@echo "ok"

build: tidy generate check
	@echo "build storage"
	@go build -tags tools ./...
	@echo "ok"

test:
	@echo "run test"
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...
	@go tool cover -html="coverage.txt" -o "coverage.html"
	@echo "ok"

integration_test:
	@echo "run integration test"
	@pushd tests \
		&& go test -race -v ./... \
		&& popd
	@echo "ok"

tidy:
	@go mod tidy
	@go mod verify

clean:
	@echo "clean generated files"
	@find . -type f -name 'generated.go' -delete
	@echo "Done"