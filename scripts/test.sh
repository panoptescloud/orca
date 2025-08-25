#!/bin/bash

REPO_ROOT=$(git rev-parse --show-toplevel)

(
    cd $REPO_ROOT
    go test -v -coverprofile $REPO_ROOT/assets/testing/coverage.out "$@"
    go tool cover -html $REPO_ROOT/assets/testing/coverage.out -o $REPO_ROOT/assets/testing/coverage.html
)