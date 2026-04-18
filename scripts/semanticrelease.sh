#!/bin/bash
# Wrapper around semantic-release

npx --yes --package semantic-release@25 \
    --package conventional-changelog-conventionalcommits@9 \
    semantic-release "$@"