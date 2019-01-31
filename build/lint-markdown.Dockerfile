FROM node:11.8-alpine AS lint-markdown
WORKDIR /app

# hadolint ignore=DL3018
RUN apk add --no-cache bash

RUN npm install -g \
    remark-cli@6.0.1 \
    remark-preset-lint-recommended@3.0.2 \
    remark-lint-sentence-newline@2.0.0 \
    remark-preset-lint-consistent@2.0.2 \
    remark-lint-no-dead-urls@0.4.1

ENTRYPOINT [ "scripts/lint.sh", "markdown" ]
