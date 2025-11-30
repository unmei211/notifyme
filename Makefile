PROTO_DIR := internal/pkg/proto
.PHONY: clean start soft-clean start-follow

.PHONY: proto-gen
proto-gen:
	@echo "Generating proto for all services"
	@cd $(PROTO_DIR) && \
	PROTO_FILES=$$(find . -name "*.proto") && \
	if [ -z "$$PROTO_FILES" ]; then \
		echo "No proto files found"; exit 1; \
	fi && \
	protoc \
	  --go_out=. \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=. \
	  --go-grpc_opt=paths=source_relative \
	  $$PROTO_FILES

start:
	DOCKER_BUILDKIT=1 docker compose -f compose.yaml -f compose.override.yaml up -d --build
start-follow:
	DOCKER_BUILDKIT=1 docker compose -f compose.yaml -f compose.override.yaml up --build
clean:
	docker compose -f compose.yaml -f compose.override.yaml down -v
soft-clean:
	docker compose -f compose.yaml -f compose.override.yaml down

reload-services:
	docker compose up -d hub_submitter --build

infra:
	docker compose up -d \
	infra-kafka infra-postgres infra-kafka2 infra-kafka3 \
	infra-tofu infra-kafka-ui --build