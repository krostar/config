# hadolint ignore=DL3006
FROM python:3.7.2-alpine3.8 AS lint-yaml
WORKDIR /app

# hadolint ignore=DL3018
RUN apk add --no-cache bash
RUN pip install --no-cache-dir yamllint==1.14.0

ENTRYPOINT [ "scripts/lint.sh", "yaml" ]
