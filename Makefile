PROTO_DIR := internal/pkg/proto

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
