FROM golang:1.11-alpine AS publish-cover
WORKDIR /app

# hadolint ignore=DL3018
RUN apk add --no-cache git

RUN go get github.com/schrej/godacov

ENTRYPOINT [ "godacov" ]
