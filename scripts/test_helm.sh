#!/usr/bin/env bash
#
# This script renders Helm templates from the local repository,
# then validates the result.

set -euo pipefail
cd "$(dirname "$0")"
source "./lib.sh"

check_dependencies \
  git \
  helm \
  kube-linter

PROJECT_ROOT=$(git rev-parse --show-toplevel)

mapfile -t EXAMPLES < <(
  find "$PROJECT_ROOT/examples/" -mindepth 1 -type d -printf "%f\n"
)

BUILD="$PROJECT_ROOT/build"
mkdir --parents "$BUILD"

for example in "${EXAMPLES[@]}"; do
  run_trace false helm template "$example" "$PROJECT_ROOT" \
    --create-namespace \
    --namespace=coder-test \
    --release-name \
    --values="$PROJECT_ROOT/examples/images.yaml" \
    --values="$PROJECT_ROOT/examples/$example/$example.values.yaml" \
    --output-dir="$BUILD"
done

run_trace false kube-linter lint --config="$PROJECT_ROOT/kube-linter.yaml" "$BUILD"
