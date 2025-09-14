GO_FILTER_NAME ?= graphql-federation
GOPROXY := $(shell go env GOPROXY)
GOARCH ?= arm64
BUILD_DIR := build

.PHONY: all build clean

all: build

test-up:
	docker-compose -f graphql-plugin/scripts/docker-compose.yaml up -d
	docker-compose -f graphql-plugin/scripts/docker-compose.yaml ps
	docker logs -f scripts-envoy-1

test-down:
	docker-compose -f graphql-plugin/scripts/docker-compose.yaml down

clean:
	rm -rf $(BUILD_DIR)

build:
	DOCKER_BUILDKIT=1 docker build --build-arg GOPROXY=$(GOPROXY) \
								    --build-arg GO_FILTER_NAME=${GO_FILTER_NAME} \
									--build-arg GOARCH=${GOARCH} \
									-t ${GO_FILTER_NAME} \
									--output ./graphql-plugin/build/ \
									.