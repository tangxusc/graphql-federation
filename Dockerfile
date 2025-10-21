#FROM golang:1.22-bullseye AS golang-base
FROM docker.io/library/golang:1.23-bullseye AS golang-base

ARG GOPROXY=https://proxy.golang.org,direct
ARG GO_FILTER_NAME=graphql-federation
ARG GOARCH=arm64

ENV GOFLAGS=-buildvcs=false
ENV GOPROXY=${GOPROXY}
ENV GOARCH=${GOARCH}
ENV CGO_ENABLED=1

# 根据目标架构安装对应的编译工具
RUN if [ "$GOARCH" = "arm64" ]; then \
        echo "Installing ARM64 toolchain" && \
        apt-get update && \
        apt-get install -y gcc-aarch64-linux-gnu binutils-aarch64-linux-gnu; \
    else \
        echo "Installing AMD64 toolchain" && \
        apt-get update && \
        apt-get install -y gcc-x86-64-linux-gnu binutils-x86-64-linux-gnu; \
    fi

WORKDIR /workspace
COPY . .

RUN cd ./graphql-plugin && go mod tidy
RUN if [ "$GOARCH" = "arm64" ]; then \
       pwd && ls -la CC=aarch64-linux-gnu-gcc AS=aarch64-linux-gnu-as go build -o /tmp/plugin.so -buildmode=c-shared /workspace/graphql-plugin/cmd/graphql; \
    else \
        pwd && ls -la && CC=x86_64-linux-gnu-gcc AS=x86_64-linux-gnu-as go build -o /tmp/plugin.so -buildmode=c-shared /workspace/graphql-plugin/cmd/graphql; \
    fi

FROM higress-registry.cn-hangzhou.cr.aliyuncs.com/higress/gateway:latest
ARG GO_FILTER_NAME=graphql-federation
ARG GOARCH=arm64
COPY --from=golang-base /tmp/plugin.so /var/lib/istio/envoy/${GO_FILTER_NAME}_${GOARCH}.so
