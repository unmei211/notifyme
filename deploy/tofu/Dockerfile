FROM ghcr.io/opentofu/opentofu:minimal as tofu

FROM eclipse-temurin:21-alpine

COPY --from=tofu /usr/local/bin/tofu /usr/local/bin/tofu

CMD ["tofu", "apply", "-auto-approve"]
