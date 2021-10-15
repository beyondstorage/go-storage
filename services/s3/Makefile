SHELL := /bin/bash

-include Makefile.env

.PHONY: all check format vet lint build test generate tidy

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check               to do static check"
	@echo "  build               to create bin directory and build"
	@echo "  generate            to generate code"
	@echo "  test                to run test"

check: vet

format:
	go fmt ./...

vet:
	go vet ./...

generate:
	@echo "generate code"
	go generate ./...
	go fmt ./...

build: tidy generate check
	go build ./...

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -v .
	go tool cover -html="coverage.txt" -o "coverage.html"

integration_test:
	go test -count=1 -race -covermode=atomic -v ./tests

tidy:
	go mod tidy
	go mod verify

clean:
	@echo "clean generated files"
	find . -type f -name 'generated.go' -delete
