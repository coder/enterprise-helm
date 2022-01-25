#!/usr/bin/env bash
#
# This script installs dependencies to /usr/local/bin.

set -euo pipefail
cd "$(dirname "$0")"
source "./lib.sh"

TMPDIR=$(mktemp -d)
BINDIR="/usr/local/bin"

curl_flags=(
  --silent
  --show-error
  --location
)

run_trace false sudo apt-get install --no-install-recommends --yes \
  shellcheck

# gotestsum makes test output more readable
GOTESTSUM_VERSION=1.7.0
run_trace false curl "${curl_flags[@]}" "https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz" \| \
  tar --extract --gzip --directory="$TMPDIR" --file=- gotestsum

# golangci-lint to lint Go code with multiple tools
GOLANGCI_LINT_VERSION=1.43.0
run_trace false curl "${curl_flags[@]}" "https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz" \| \
  tar --extract --gzip --directory="$TMPDIR" --file=- --strip-components=1 "golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64/golangci-lint"

# Helm for packaging and validating chart
HELM_VERSION=3.8.0
run_trace false curl "${curl_flags[@]}" "https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz" \| \
  tar --extract --gzip --directory="$TMPDIR" --file=- --strip-components=1 linux-amd64/helm

# helm-docs for generating README.md
HELM_DOCS_VERSION=1.5.0
run_trace false curl "${curl_flags[@]}" "https://github.com/norwoodj/helm-docs/releases/download/v${HELM_DOCS_VERSION}/helm-docs_${HELM_DOCS_VERSION}_Linux_x86_64.tar.gz" \| \
  tar --extract --gzip --directory="$TMPDIR" --file=- helm-docs

# kube-linter for checking generated YAML output
KUBE_LINTER_VERSION=0.2.5
run_trace false curl "${curl_flags[@]}" "https://github.com/stackrox/kube-linter/releases/download/${KUBE_LINTER_VERSION}/kube-linter-linux.tar.gz" \| \
  tar --extract --gzip --directory="$TMPDIR" --file=- kube-linter

run_trace false sudo install --mode=0755 --target-directory="$BINDIR" "$TMPDIR/*"

run_trace false command -v \
  gotestsum \
  golangci-lint \
  helm \
  helm-docs \
  kube-linter

run_trace false helm version --short

run_trace false rm --verbose --recursive --force "$TMPDIR"
