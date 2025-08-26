#!/bin/bash

REPO_ROOT=$(git rev-parse --show-toplevel)

OUTDIR=${OUTDIR:-$REPO_ROOT/.local/bin/orca}

go build \
    -ldflags="-X 'main.version=${VERSION:-unknown}' -X 'main.commit=${COMMIT:-unknown}' -X 'main.date=${DATE:-unknown}'" \
    -o $OUTDIR \
    $REPO_ROOT/cmd/orca/*.go