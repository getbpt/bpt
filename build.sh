#!/usr/bin/env bash

function run() {
  source "$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)/scripts/common.sh"
  prepare || fatal "failed to prepare: $?"
  build || fatal "failed to build: $?"
  package || fatal "failed to package: $?"
}

run "$@"
