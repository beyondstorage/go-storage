SHELL := /bin/bash

-include Makefile.env

.PHONY: all check format lint build test generate tidy

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
	# Storj is much slower than we expected: https://github.com/beyondstorage/go-storage/issues/958
	# 10 minutes is not enough for storj to finish our test.
	# So we increase timeout to 30m.
	go test -race -count=1 -covermode=atomic -v -timeout=30m ./tests

tidy:
	go mod tidy
	go mod verify

clean:
	@echo "clean generated files"
	find . -type f -name 'generated.go' -delete
