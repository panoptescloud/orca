#!/bin/bash

REPO_ROOT=$(git rev-parse --show-toplevel)

go build \
    -ldflags="-X 'main.version=${VERSION:-unknown}' -X 'main.commit=${COMMIT:-unknown}' -X 'main.date=${DATE:-unknown}'" \
    -o $REPO_ROOT/.local/bin/orca \
    $REPO_ROOT/cmd/orca/*.go