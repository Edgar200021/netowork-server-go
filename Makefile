.PHONY: help
help:
	@echo "Usage:"

.PHONY: test
test:
	go run ./cmd/api/main.go
