FROM golang:1.25-alpine AS builder

COPY ./internal/pkg/go.mod ./internal/pkg/go.sum /src/pkg/
COPY ./internal/services/hub_submitter/go.mod ./internal/services/hub_submitter/go.sum /src/services/hub_submitter/

WORKDIR /src/services/hub_submitter

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

WORKDIR /

COPY ./internal/pkg /src/pkg/
COPY ./internal/services/hub_submitter /src/services/hub_submitter/

WORKDIR /src/services/hub_submitter/cmd/

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /opt/hub_submitter/app .

CMD ["app"]

