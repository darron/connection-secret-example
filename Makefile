BINARY_NAME ?= app
CONTAINER_NAME ?= darron/connection-secret-example

BUILD_COMMAND=-mod=vendor -o bin/$(BINARY_NAME) main.go
UNAME=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m)

all: build

deps: ## Install all dependencies.
	go mod vendor
	go mod tidy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Remove compiled binaries.
	rm -f bin/$(BINARY_NAME) || true
	rm -f bin/$(BINARY_NAME)*gz || true

docker: ## Build Docker image
	docker build . -t $(CONTAINER_NAME)

build: clean
	go build $(BUILD_COMMAND)

rebuild: clean ## Force rebuild of all packages.
	go build -a $(BUILD_COMMAND)

linux: clean ## Cross compile for linux.
	CGO_ENABLED=0 GOOS=linux go build $(BUILD_COMMAND)

gzip: ## Compress current compiled binary.
	gzip bin/$(BINARY_NAME)
	mv bin/$(BINARY_NAME).gz bin/$(BINARY_NAME)-$(UNAME)-$(ARCH).gz

release: build gzip ## Full release process.

unit: ## Run unit tests.
	go test -mod=vendor -cover -race -short ./... -v

lint: ## See https://github.com/golangci/golangci-lint#install for install instructions
	golangci-lint run ./...

redis: ## Run Redis server with local redis.conf
	redis-server redis.conf

## SSL Additions from: https://nishanths.svbtle.com/setting-up-redis-with-tls

FQDN ?= 127.0.0.1
OUTDIR ?= tls
# may be /etc/pki/tls in some machines.
# use `openssl version -a | grep OPENSSLDIR` to find out - this works on OSX.
OPENSSLDIR ?= /etc/ssl

generate: prepare redis.crt cleancerts ## Re-generate all SSL/TLS certificates

prepare:
	mkdir ${OUTDIR}

cleancerts:
	rm -f ${OUTDIR}/openssl.cnf

openssl.cnf:
	cat ${OPENSSLDIR}/openssl.cnf > ${OUTDIR}/openssl.cnf
	echo "" >> ${OUTDIR}/openssl.cnf
	echo "[ san_env ]" >> ${OUTDIR}/openssl.cnf
	echo "subjectAltName = IP:${FQDN}" >> ${OUTDIR}/openssl.cnf

ca.key:
	openssl genrsa 4096 > ${OUTDIR}/ca.key

ca.crt: ca.key
	openssl req \
		-new \
		-x509 \
		-nodes \
		-sha256 \
		-key ${OUTDIR}/ca.key \
		-days 3650 \
		-subj "/C=AU/CN=example" \
		-out ${OUTDIR}/ca.crt

redis.csr: openssl.cnf
	# is -extensions necessary?
	# https://security.stackexchange.com/a/86999
	SAN=IP:$(FQDN) openssl req \
		-reqexts san_env \
		-extensions san_env \
		-config ${OUTDIR}/openssl.cnf \
		-newkey rsa:4096 \
		-nodes -sha256 \
		-keyout ${OUTDIR}/redis.key \
		-subj "/C=AU/CN=$(FQDN)" \
		-out ${OUTDIR}/redis.csr

redis.crt: openssl.cnf ca.key ca.crt redis.csr
	SAN=IP:$(FQDN) openssl x509 \
		-req -sha256 \
		-extfile ${OUTDIR}/openssl.cnf \
		-extensions san_env \
		-days 3650 \
		-in ${OUTDIR}/redis.csr \
		-CA ${OUTDIR}/ca.crt \
		-CAkey ${OUTDIR}/ca.key \
		-CAcreateserial \
		-out ${OUTDIR}/redis.crt

.PHONY: help all deps clean build gzip release unit lint redis generate prepare cleancerts docker