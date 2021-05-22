#!/usr/bin/env bash
#
# This script generates the README from the Go template.

set -euo pipefail
cd "$(dirname "$0")"
source "./lib.sh"

check_dependencies \
  git \
  helm-docs \
  make

PROJECT_ROOT="$(git rev-parse --show-toplevel)"

pushd "$PROJECT_ROOT" >/dev/null
  run_trace false make --always-make README.md

  FILES="$(git ls-files --other --modified --exclude-standard)"
  if [ -n "$FILES" ]; then
    echo "Unstaged files after running 'make fmt':"
    echo "$FILES"
    exit 1
  fi
popd >/dev/null
