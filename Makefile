PROJECT_ROOT := $(shell git rev-parse --show-toplevel)

.DEFAULT_GOAL: help
help:
	@echo "Usage: $(MAKE) <target>"
	@echo
	@echo " * 'all' - Run everything"
	@echo " * 'docs' - Generate documentation"
	@echo " * 'lint' - Run linters"
.PHONY: help

all: lint docs
.PHONY: all

lint: lint/helm
.PHONY: lint

lint/helm:
	@echo "--- Running helm lint"
	# TODO(jawnsy): enable --strict once we fix the warnings
	helm lint .
.PHONY: lint/helm

docs: README.md
.PHONY: docs

README.md: README.md.gotmpl
	@echo "--- Generating documentation"
	helm-docs --template-files=$<
