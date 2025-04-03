APP_NAME := note-nest
GO := go
GOBUILD := $(GO) build
GOFORMAT := go fmt
GOTEST := $(GO) test
GOVET := $(GO) vet

.PHONY: help
help:
	@echo "build - build the project"
	@echo "clean - clean the project"
	@echo "format - format the code"
	@echo "dev - run the project with hot reload"
	@echo "docker-build - build the docker image"
	@echo "docker-down - stop the docker container"
	@echo "docker-up - start the docker container"
	@echo "start - run the project"
	@echo "test - run tests"
	@echo "vet - run the vet tool"

.PHONY: build
build:
	$(GOBUILD) -o bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

.PHONY: clean
clean:
	rm -rf bin/$(APP_NAME) tmp/

.PHONY: format
format:
	$(GOFORMAT) ./...

.PHONY: dev
dev:
	air

.PHONY: docker-build
docker-build:
	docker compose build

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: start
start:
	./bin/$(APP_NAME)

.PHONY: test
test: test-unit test-integration

test-unit:
	$(GOTEST) -v ./internal/...

test-integration:
	$(GOTEST) -v ./tests/integration

.PHONY: vet
vet:
	$(GOVET) ./...