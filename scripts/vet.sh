#!/bin/bash

REPO_ROOT=$(git rev-parse --show-toplevel)

(
    cd $REPO_ROOT
    go vet ./...
)