FROM hadolint/hadolint:v1.15.0-debian AS lint-dockerfile
WORKDIR /app

# hadolint ignore=DL3008
RUN apt-get update && \
    apt-get install --yes --no-install-recommends bash && \
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "scripts/lint.sh", "dockerfile" ]
