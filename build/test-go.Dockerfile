FROM golang:1.11-stretch AS lint-go
VOLUME [ "/go/pkg/mod" ]
WORKDIR /app

ENTRYPOINT [ "scripts/test.sh", "go" ]
