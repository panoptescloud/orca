#!/bin/bash

REPO_ROOT=$(git rev-parse --show-toplevel)

go build -o $REPO_ROOT/bin/orca $REPO_ROOT/cmd/orca/*.go