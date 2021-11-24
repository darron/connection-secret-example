CONTAINER_NAME ?= darron/connection-secret-example:hit2

all: docker

docker: ## Build Docker image
	docker build . -t $(CONTAINER_NAME)

docker_two: ## Build second Docker image
	docker build -f Dockerfile.two . -t $(CONTAINER_NAME)2

.PHONY: docker