SHELL := /bin/bash

.PHONY: all check format vet lint build test generate tidy

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  check      to format, vet and lint "
	@echo "  build      to create bin directory and build"
	@echo "  generate   to generate code"
	@echo "  test       to run test"

# mockgen: go get github.com/golang/mock/mockgen
# golint: go get -u golang.org/x/lint/golint
# go-bindata: go get -u github.com/kevinburke/go-bindata/...
tools := mockgen golint go-bindata

$(tools):
	@command -v $@ >/dev/null 2>&1 || echo "$@ is not found, plese install it."

check: vet lint

format:
	@echo "go fmt"
	@go fmt ./...
	@echo "ok"

vet:
	@echo "go vet"
	@go vet ./...
	@echo "ok"

lint: golint
	@echo "golint"
	@golint ./...
	@echo "ok"

build_generator: go-bindata
	@echo "build storage generator"
	@pushd internal/cmd \
		&& go generate ./... \
		&& go build -o ../bin/service ./service \
		&& go build -o ../bin/pairs ./pairs \
		&& go build -o ../bin/metadata ./metadata \
		&& popd
	@echo "Done"

generate: build_generator mockgen
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

tidy:
	@echo "Tidy and check the go mod files"
	@go mod tidy && go mod verify
	@pushd internal/cmd && go mod tidy && go mod verify && popd
	@echo "Done"

clean:
	@echo "Clean generated files"
	@rm services/*/generated.go
	@echo "Done"