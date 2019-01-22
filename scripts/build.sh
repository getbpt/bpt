#!/usr/bin/env bash

function run() {
  source "$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)/common.sh"
  build || fatal "failed to build: $?"
}

run "$@"
