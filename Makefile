PROJECT_ROOT := $(shell git rev-parse --show-toplevel)

.DEFAULT_GOAL: help
help:
	@echo "Usage: $(MAKE) <target>"
	@echo
	@echo " * 'all' - Run everything"
	@echo " * 'fmt' - Run formatters"
	@echo " * 'lint' - Run linters"
.PHONY: help

all: lint fmt
.PHONY: all

lint: lint/helm
.PHONY: lint

lint/helm:
	@echo "--- Running helm lint"
	# TODO(jawnsy): enable --strict once we fix the warnings
	helm lint .
.PHONY: lint/helm

# TODO(jawnsy): this will be modified to format using Prettier
fmt: README.md
.PHONY: fmt

README.md: README.md.gotmpl
	@echo "--- Generating documentation"
	helm-docs --template-files=$<
