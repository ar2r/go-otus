BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

clean:
	rm -rf $(BIN)

version: build
	$(BIN) version

test:
	# Если тесты запускаются с TZ=UTC или TZ="", происходит какая-то странная магия. Подробности в пул реквесте.
	# Поэтому запускаем тесты с TZ=Europe/Moscow 🤬
	TZ=Europe/Moscow go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.0

lint: install-lint-deps
	golangci-lint run ./...

generate:
	rm -rf ./internal/server/grpc/protobuf
	mkdir -p ./internal/server/grpc/protobuf

	protoc -I ./api \
	   --go_out internal/server/grpc/protobuf \
	   --go_opt paths=source_relative \
	   --go-grpc_out internal/server/grpc/protobuf \
	   --go-grpc_opt paths=source_relative \
	   ./api/*.proto
	go generate ./...

migrate: build
	$(BIN) -config ./configs/config.toml migrate

up:
	docker compose up -d

down:
	docker compose down -v

.PHONY: build run build-img run-img version test lint
