SHELL=bash

.PHONY: test
test:
	go test -count=1 -race -cover ./...

.PHONY: audit
audit:
	go list -json -m all | nancy sleuth --exclude-vulnerability-file ./.nancy-ignore

.PHONY: build
build:
	go build ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test-integration
test-integration:
	docker-compose down
	docker-compose up -d
	sleep 10
	go test -count=1 -race -cover -tags="integration" ./...
	docker-compose down
