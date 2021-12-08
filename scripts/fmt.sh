#!/usr/bin/env bash
#
# This script renders Helm templates from the local repository,
# then validates the result.

set -euo pipefail
cd "$(dirname "$0")"
source "./lib.sh"

check_dependencies \
  git \
  helm-docs \
  make

PROJECT_ROOT="$(git rev-parse --show-toplevel)"

pushd "$PROJECT_ROOT" >/dev/null
  run_trace false make --always-make fmt

  FILES="$(git ls-files --other --modified --exclude-standard)"
  if [ -n "$FILES" ]; then
    mapfile -t files <<< "$FILES"

    echo "The following files contain unstaged changes:"
    echo
    for file in "${files[@]}"
    do
      echo "  - $file"
    done
    echo

    echo "These are the changes:"
    echo
    for file in "${files[@]}"
    do
      git --no-pager diff "$file"
    done

    exit 1
  fi
popd >/dev/null
