.PHONY: build test lint ci

## build: compile the project
build:
	go build ./...

## test: run unit tests
test:
	go test ./...

## lint: run golangci-lint
lint:
	golangci-lint run

## ci: run lint and tests (used by GitHub Actions)
ci: lint test
