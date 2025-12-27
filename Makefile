.PHONY: run test test-lint build

run:
	go run ./cmd/api/main.go

test:
	go test -v $$(go list ./tests/... | grep -v /testapp) -p 20 -count=1

build:
	go build -o bin/api ./cmd/api
