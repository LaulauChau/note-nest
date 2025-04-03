.PHONY: help
help:
	@echo "build - build the project"
	@echo "clean - clean the project"
	@echo "run - run the project"
	@echo "dev - run the project with hot reload"
	@echo "docker-build - build the docker image"
	@echo "docker-down - stop the docker container"
	@echo "docker-up - start the docker container"
	@echo "test - run tests"

.PHONY: build
build:
	go build -o bin/note-nest cmd/note-nest/main.go

.PHONY: clean
clean:
	rm -rf bin/note-nest tmp/

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

.PHONY: test
test: test-unit test-integration

test-unit:
	go test -v ./...

test-integration:
	go test -v ./tests/integration
