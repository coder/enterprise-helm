#!/usr/bin/env bash

set -euo pipefail
cd "$(dirname "$0")"

# Ensure $GOPATH/bin exists within the $PATH. This is where gobin and helm-docs
# will be installed.
if [[ "$PATH" != *"$(go env GOPATH)"* ]]; then
	echo "Please ensure $(go env GOPATH)/bin is within your \$PATH"
fi

# Ensure helm-docs and gobin are installed.
if ! which helm-docs &> /dev/null; then 
	echo "Installing helm-docs..."
	if ! which gobin &> /dev/null; then 
		echo "Installing gobin..."
		GO111MODULE=off go get -u github.com/myitcv/gobin
	fi

	# Use gobin to install helm-docs since it is a modules enabled repo.
	gobin github.com/coadler/helm-docs/cmd/helm-docs
fi

helm-docs
prettier --write README.md
