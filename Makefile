SHELL := /bin/bash
PACKAGES = credential endpoint services/s3

.PHONY: all check format vet lint build test generate tidy integration_test $(PACKAGES)

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check               to do static check"
	@echo "  build               to create bin directory and build"
	@echo "  generate            to generate code"
	@echo "  test                to run test"
	@echo "  integration_test    to run integration test"

check: vet

format:
	@echo "gofmt"
	@gofmt -w -l .
	@echo "ok"

vet:
	@echo "go vet"
	@go vet ./...
	@echo "ok"

generate:
	@echo "generate code"
	@go generate -tags tools ./...
	@gofmt -w -l .
	@echo "ok"

build: tidy generate format check
	@echo "build storage"
	@go build -tags tools ./...
	@echo "ok"

build-all: build $(PACKAGES)

$(PACKAGES):
	pushd $@ && make build && popd

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
