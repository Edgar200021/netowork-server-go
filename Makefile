.PHONY: run test test-lint build

run:
	go run ./cmd/api/main.go

test:
	go test -v $$(go list ./tests/... | grep -v /testapp) -p 100 -count=1

test-lint:
	testifylint ./...

test-lint-fix:
	testifylint --fix ./...

build:
	go build -o bin/api ./cmd/api
