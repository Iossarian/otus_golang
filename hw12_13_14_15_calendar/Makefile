CALENDAR_BIN := "./bin/calendar"
SCHEDULER_BIN := "./bin/sheduler"
SENDER_BIN := "./bin/sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

up: build
	docker-compose --env-file deployments/.env -f deployments/docker-compose.yaml up -d --build --remove-orphans

down:
	docker-compose --env-file deployments/.env -f deployments/docker-compose.yaml down

generate:
	go generate ./...

migrate-up:
	goose --dir=migrations postgres "postgres://calendar:calendar@calendar_db:5432/calendar?sslmode=disable" up

migrate-down:
	goose --dir=migrations postgres "postgres://calendar:calendar@calendar_db:5432/calendar?sslmode=disable" down


build-calendar:
	go build -v -o $(CALENDAR_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(SCHEDULER_BIN) -ldflags "$(LDFLAGS)" ./cmd/scheduler

build-sender:
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/sender

build: build-calendar build-scheduler build-sender

run-calendar: build-calendar
	$(CALENDAR_BIN) -config ./configs/calendar_config.yml

run-scheduler: build-scheduler
	$(SCHEDULER_BIN) -config ./configs/scheduler_config.yml

run-sender: build-sender
	$(SENDER_BIN) -config ./configs/sender_config.yml

integration-tests: up
	chmod +x ./integration_test.sh && ./integration_test.sh

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
