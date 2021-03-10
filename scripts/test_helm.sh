#!/usr/bin/env bash
#
# This script renders Helm templates from the local repository,
# then validates the result.

set -euo pipefail
cd "$(dirname "$0")"
source "./lib.sh"

check_dependencies \
  git \
  helm

PROJECT_ROOT="$(git rev-parse --show-toplevel)"

EXAMPLES=" \
  kind \
"

BUILD="$PROJECT_ROOT/build"
mkdir -p "$BUILD"

for example in $EXAMPLES; do
  run_trace false helm template "$PROJECT_ROOT" \
    --release-name \
    --name-template="$example" \
    --values="$PROJECT_ROOT/examples/release.yaml" \
    --values="$PROJECT_ROOT/examples/$example/$example.values.yaml" \
    --output-dir="$BUILD" \| indent
done

run_trace false kube-linter lint "$BUILD"
