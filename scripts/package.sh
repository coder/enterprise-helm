#!/usr/bin/env bash
#
# Package up the Helm chart and list the files.

set -euo pipefail
cd "$(dirname "$0")"
source "./lib.sh"

check_dependencies \
  git \
  helm

PROJECT_ROOT="$(git rev-parse --show-toplevel)"

BUILD="$PROJECT_ROOT/build"
mkdir -p "$BUILD"

run_trace false helm package "$PROJECT_ROOT" \
  --destination="$BUILD" \
  --version=0.1

run_trace false tar --verbose --list --gzip --file="$BUILD/coder-0.1.tgz"
