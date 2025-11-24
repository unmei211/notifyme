FROM golang:1.25-alpine AS build

WORKDIR /app

COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/hub_submitter/go.mod services/hub_submitter/go.sum ./src/service/

WORKDIR /app/src/service

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/hub_submitter/ ./src/service

WORKDIR /app/src/service/cmd/

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -gcflags="all=-N -l" -o /opt/service/app .

RUN mv ../migrations /opt/service/migrations

WORKDIR /opt/service

EXPOSE 8056
CMD ["./app"]

