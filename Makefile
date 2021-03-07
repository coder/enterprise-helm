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

docs:
.PHONY: docs

lint: lint/markdown lint/helm lint/shellcheck
.PHONY: lint

lint/markdown:
	@echo "--- Running markdownlint"
.PHONY: lint/markdown

lint/helm:
.PHONY: lint/helm

lint/shellcheck: $(shell scripts/depfind/sh.sh)
	shellcheck $^
.PHONY: lint/shellcheck

fmt:
	@echo "abc"
.PHONY: fmt

fmt/docs: $(shell scripts/depfind/markdown.sh)
	@echo "--- Formatting documentation"
	prettier --write $^
.PHONY: fmt/docs

README.md: README.md.gotmpl
	@echo "--- Generating documentation"
	helm-docs
