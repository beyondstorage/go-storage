SHELL := /bin/bash

.PHONY: all check format vet build test generate tidy integration_test build-all

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check               to do static check"
	@echo "  build               to create bin directory and build"
	@echo "  generate            to generate code"
	@echo "  test                to run test"
	@echo "  build-all           to build all packages"

check: vet

format:
	gofmt -w -l .

vet:
	go vet ./...

generate:
	go generate -tags tools ./...
	gofmt -w -l .

build: tidy generate format check
	go build -tags tools ./...

build-all:
	for f in $$(find . -name go.mod);     \
		do make -C $$(dirname $$f) build;  \
	done

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...
	go tool cover -html="coverage.txt" -o "coverage.html"

test-all:
	for f in $$(find . -name go.mod);     \
		do make -C $$(dirname $$f) test;  \
	done

tidy:
	go mod tidy
	go mod verify

tidy-all:
	for f in $$(find . -name go.mod);     \
		do make -C $$(dirname $$f) tidy;  \
	done

clean:
	find . -type f -name 'generated.go' -delete
