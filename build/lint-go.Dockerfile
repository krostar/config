FROM golang:1.11-alpine AS lint-go
VOLUME [ "/go/pkg/mod" ]
WORKDIR /app

# hadolint ignore=DL3018
RUN apk add --no-cache bash git

RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint

ENTRYPOINT [ "scripts/lint.sh", "go" ]
