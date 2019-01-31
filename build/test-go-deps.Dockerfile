FROM golang:1.11-alpine AS lint-go
VOLUME [ "/go/pkg/mod" ]
WORKDIR /app

# hadolint ignore=DL3018
RUN apk add --no-cache bash git

ENTRYPOINT [ "scripts/test.sh", "go-deps" ]
