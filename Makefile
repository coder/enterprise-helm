PROJECT_ROOT := $(shell git rev-parse --show-toplevel)

.DEFAULT_GOAL: help
help:
	@echo "Usage: $(MAKE) <target>"
	@echo
	@echo " * 'all' - Run everything"
	@echo " * 'gen' - Generate code and documentation"
	@echo " * 'lint' - Run linters"
	@echo " * 'fmt' - Run formatters"
.PHONY: help

all: lint fmt
.PHONY: all

clean:
	rm -vrf build
.PHONY: clean

docs:
.PHONY: docs

lint: lint/helm lint/shellcheck
.PHONY: lint

lint/helm:
	@echo "--- Running helm lint"
	helm lint --strict .
.PHONY: lint/helm

lint/shellcheck: $(shell scripts/depfind/sh.sh)
	@echo "--- Running shellcheck"
	shellcheck $^
.PHONY: lint/shellcheck

lint/kube:
	@echo "--- Linting rendered templates"
	kube-linter lint build
.PHONY: lint/kube

fmt: fmt/docs
.PHONY: fmt

fmt/docs: $(shell scripts/depfind/markdown.sh)
	@echo "--- Formatting documentation"
	prettier --write $^
.PHONY: fmt/docs

README.md: README.md.gotmpl
	@echo "--- Generating documentation"
	helm-docs --template-files=$<
