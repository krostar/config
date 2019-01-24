#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set +o posix

exit::error() {
    if [ $# -eq 2 ]; then
        echo "$2"
    fi
    exit "$1"
}

project::path::root() {
    cd -P "$(dirname "$(dirname "${BASH_SOURCE[0]}")")" || exit 255 && pwd
}

project::path::build() {
    echo "$(project::path::root)/build"
}

project::path::build::cover() {
    echo "$(project::path::build)/cover"
}

project::repo() {
    local -r gomodfile="$(project::path::root)/go.mod"

    if [[ -f "$gomodfile" ]]; then
        head -n1 "${gomodfile}" | cut -d' ' -f2
    else
        echo "${$(project::path::roo)#"${GOPATH}/src/"}"
    fi
}

project::name() {
    basename "$(project::repo)"
}
