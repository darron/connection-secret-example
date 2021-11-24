CONTAINER_NAME ?= darron/connection-secret-example:hit2

all: docker

docker: ## Build Docker image
	docker build . -t $(CONTAINER_NAME)

.PHONY: docker