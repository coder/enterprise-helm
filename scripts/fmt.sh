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
  make \
  prettier

PROJECT_ROOT="$(git rev-parse --show-toplevel)"

make fmt

FILES=$(git ls-files --other --modified --exclude-standard)
if [ -n "$FILES" ]; then
  echo "Unstaged files after running 'make fmt':"
  echo "$FILES"
  exit 1
fi
