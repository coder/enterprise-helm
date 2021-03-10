PROJECT_ROOT := $(shell git rev-parse --show-toplevel)

.DEFAULT_GOAL: help
help:
	@echo "Usage: $(MAKE) <target>"
	@echo
	@echo " * 'all' - Run everything"
	@echo " * 'lint' - Run linters"
.PHONY: help

all: lint
.PHONY: all

lint: lint/helm
.PHONY: lint

lint/helm:
	@echo "--- Running helm lint"
	helm lint --strict .
.PHONY: lint/helm
