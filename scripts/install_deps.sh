#!/usr/bin/env bash
#
# This script installs dependencies to /usr/local/bin.

set -euo pipefail
cd "$(dirname "$0")"
source "./lib.sh"

TMPDIR="$(mktemp -d)"
BINDIR="/usr/local/bin"

curl_flags=(
  --silent
  --show-error
  --location
)

run_trace false sudo apt-get install --no-install-recommends --yes \
  shellcheck

run_trace false curl "${curl_flags[@]}" "https://get.helm.sh/helm-v3.5.2-linux-amd64.tar.gz" \| \
  tar -C "$TMPDIR" --strip-components=1 -zxf - linux-amd64/helm

run_trace false curl "${curl_flags[@]}" "https://github.com/norwoodj/helm-docs/releases/download/v1.5.0/helm-docs_1.5.0_Linux_x86_64.tar.gz" \| \
  tar -C "$TMPDIR" -zxf - helm-docs

run_trace false sudo install --mode=0755 "$TMPDIR/*" "$BINDIR"

run_trace false which \
  helm \
  helm-docs

run_trace false helm version --short

run_trace false rm -vrf "$TMPDIR"
