FROM golang:1.25-alpine AS build

WORKDIR /app

ARG SERVICE_NAME=hub_submitter
ARG SERVICE_HOST_PATH=./services/${SERVICE_NAME}/
ARG SERVICE_CONTAINER_PATH=src/${SERVICE_NAME}/
ARG PKG_HOST_PATH=./pkg/
ARG PKG_CONTAINER_PATH=pkg/

COPY ${PKG_HOST_PATH}go.mod ${PKG_HOST_PATH}go.sum ${PKG_CONTAINER_PATH}
COPY ${SERVICE_HOST_PATH}go.mod ${SERVICE_HOST_PATH}go.sum ${SERVICE_CONTAINER_PATH}

WORKDIR /app/${SERVICE_CONTAINER_PATH}

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

WORKDIR /app

COPY ${PKG_HOST_PATH} ${PKG_CONTAINER_PATH}
COPY ${SERVICE_HOST_PATH} ${SERVICE_CONTAINER_PATH}

WORKDIR ${SERVICE_CONTAINER_PATH}cmd/

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /opt/${SERVICE_NAME}/app .

WORKDIR /opt/${SERVICE_NAME}
CMD ["./app"]

