PROJECT_ROOT := $(shell git rev-parse --show-toplevel)

.DEFAULT_GOAL: help
help:
	@echo "Usage: $(MAKE) <target>"
	@echo
	@echo " * 'all' - Run everything"
	@echo " * 'fmt' - Run formatters"
	@echo " * 'lint' - Run linters"
	@echo " * 'clean' - Remove generated build files"
.PHONY: help

all: lint fmt
.PHONY: all

lint: lint/helm lint/kubernetes lint/shellcheck
.PHONY: lint

lint/helm:
	@echo "--- Running helm lint"
	helm lint --strict .
.PHONY: lint/helm

lint/kubernetes:
	@echo "--- Linting rendered templates"
	./scripts/test_helm.sh
.PHONY: lint/kubernetes

lint/shellcheck: $(shell scripts/depfind/sh.sh)
	@echo "--- Running shellcheck"
	shellcheck $^
.PHONY: lint/shellcheck

# TODO(jawnsy): this will be modified to format using Prettier
fmt: README.md
.PHONY: fmt

README.md: README.md.gotmpl
	@echo "--- Generating documentation"
	helm-docs --template-files=$<

clean:
	rm -vrf build/
.PHONY: clean
