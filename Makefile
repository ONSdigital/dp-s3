SHELL=bash

test:
	go test -count=1 -race -cover ./...
.PHONY: test

audit:
	go list -json -m all | nancy sleuth --exclude-vulnerability-file ./.nancy-ignore
.PHONY: audit

build:
	go build ./...
.PHONY: build

lint:
	exit
.PHONY: lint

test-integration:
	docker-compose down
	docker-compose up -d
	sleep 10
	go test -count=1 -race -cover -tags="integration" ./...
	docker-compose down
