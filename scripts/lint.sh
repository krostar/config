#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set +o posix

# shellcheck disable=SC1090
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

test::go() {
    golangci-lint run --config "$(dirname "${BASH_SOURCE[0]}")/lint-go-config.yaml"
}

test::sh() {
    shellcheck --check-sourced --external-sources --severity=info --shell=bash ./**/*.sh
}

test::yaml() {
    yamllint --config-file "$(dirname "${BASH_SOURCE[0]}")/lint-yaml-config.yaml" --strict .
}

lint() {
    local -r lint_type="$1"

    echo "running ${lint_type} linters ..."

    case "$lint_type" in
    "go")
        test::go
        ;;
    "sh")
        test::sh
        ;;
    "yaml")
        test::yaml
        ;;
    "all")
        test::go
        test::sh
        test::yaml
        ;;
    *)
        exit::error -1 "${lint_type} is not a lint type"
    ;;
    esac

    echo "${lint_type} linters ran without errors"
}

lint "$1"
