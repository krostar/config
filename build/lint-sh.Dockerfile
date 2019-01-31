FROM koalaman/shellcheck-alpine:v0.6.0 AS lint-sh
WORKDIR /app

# hadolint ignore=DL3018
RUN apk add --no-cache bash

ENTRYPOINT [ "scripts/lint.sh", "sh" ]
