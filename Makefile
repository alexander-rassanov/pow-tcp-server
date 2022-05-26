REGISTRY ?=
TAG ?= latest

run-client:
	go run cmd/client/main.go getQuote

run-server:
	go run cmd/server/main.go run

lint:
	golangci-lint run

build-client-image:
	DOCKER_BUILDKIT=0 docker build -t ${REGISTRY}client:${TAG} . -f cmd/client/Dockerfile

build-server-image:
	docker build -t ${REGISTRY}server:${TAG} . -f cmd/server/Dockerfile

test:
	go test ./...
