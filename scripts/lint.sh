#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set +o posix

# shellcheck disable=SC1090
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

test::dockerfile() {
    for file in "$(project::path::build)"/*.Dockerfile; do
        hadolint "$file"
    done
}

test::go() {
    CGO_ENABLED=0 golangci-lint run \
        --config "$(dirname "${BASH_SOURCE[0]}")/lint-go-config.yaml"
}

test::markdown() {
    remark \
        --rc-path "$(dirname "${BASH_SOURCE[0]}")/lint-markdown.yaml" \
        --frail \
        .
}

test::sh() {
    shellcheck \
        --check-sourced \
        --external-sources \
        --severity=info \
        --shell=bash \
        ./**/*.sh
}

test::yaml() {
    yamllint \
        --config-file "$(dirname "${BASH_SOURCE[0]}")/lint-yaml-config.yaml" \
        --strict \
        .
}

lint() {
    local -r lint_type="$1"

    echo "running ${lint_type} linters ..."
    case "$lint_type" in
    "dockerfile")
        test::dockerfile
        ;;
    "go")
        test::go
        ;;
    "markdown")
        test::markdown
        ;;
    "sh")
        test::sh
        ;;
    "yaml")
        test::yaml
        ;;
    *)
        exit::error -1 "${lint_type} is not a lint type"
    ;;
    esac
    echo "${lint_type} linters ran without errors"
}

lint "$1"
