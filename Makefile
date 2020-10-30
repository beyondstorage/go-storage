SHELL := /bin/bash

.PHONY: all check format vet lint build test generate tidy integration_test

EXE_SUFFIX := ""
ifeq ($(GOOS), "linux")
	EXE_SUFFIX += ".exe"
endif

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

build_definitions:
	@echo "build storage generator"
	@pushd cmd/definitions \
		&& go generate ./... \
		&& mkdir -p ../../bin/
		&& CGO_ENABLED=0 go build -o ../../bin/ . \
		&& popd
	@echo "build iterator generator"
	@pushd internal/cmd && mkdir -p ../bin/ && go build -o ../bin/ ./iterator && popd
	@echo "Done"

generate: build_definitions
	@echo "generate code"
	@go generate ./...
	@go fmt ./...
	@echo "ok"

build: generate tidy check
	@echo "build storage"
	@go build ./...
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
	@pushd internal/cmd && go build -o ../bin/gomod${EXE_SUFFIX} ./gomod && popd
	@./internal/bin/gomod

clean:
	@echo "clean generated files"
	@find . -type f -name 'generated.go' -delete
	@echo "Done"