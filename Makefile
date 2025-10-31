GO_FILTER_NAME ?= graphql-federation
GOPROXY := $(shell go env GOPROXY)
GOARCH ?= arm64
BUILD_DIR := build
IMAGE_TAG ?= latest

.PHONY: all build-local build-image clean test-up test-down

all: build-local

# 本地测试环境：复制so文件到宿主机
build-local:
	DOCKER_BUILDKIT=1 docker build --build-arg GOPROXY=$(GOPROXY) \
								    --build-arg GO_FILTER_NAME=${GO_FILTER_NAME} \
									--build-arg GOARCH=${GOARCH} \
									-t ${GO_FILTER_NAME}:local \
									--output ./build/ \
									-f Dockerfile_local \
									.

# 其他环境：直接打包为镜像
build-image:
	DOCKER_BUILDKIT=1 docker build --build-arg GOPROXY=$(GOPROXY) \
								    --build-arg GO_FILTER_NAME=${GO_FILTER_NAME} \
									--build-arg GOARCH=${GOARCH} \
									-t ${GO_FILTER_NAME}:${IMAGE_TAG} \
									-f Dockerfile \
									.

test-up:
	docker-compose -f scripts/docker-compose.yaml up -d
	docker-compose -f scripts/docker-compose.yaml ps
	docker logs -f scripts_envoy_1

test-down:
	docker-compose -f scripts/docker-compose.yaml down

clean:
	rm -rf $(BUILD_DIR)