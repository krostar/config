# use this rule as the default make rule
.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "Available targets descriptions:"
	@grep -E '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint-all test-all
lint-all: lint-go lint-markdown lint-yaml	## Run all possible linters
test-all: test-go test-go-deps				## Run all possible tests

.PHONY: lint-go lint-markdown lint-yaml lint-sh
lint-go: ci-lint-go				## Lint go files
lint-sh: ci-lint-sh				## Lint sh files
lint-yaml: ci-lint-yaml			## Lint yaml files
lint-markdown: ci-lint-markdown	## Lint markdown files

.PHONY: test-go test-go-short test-go-deps
test-go: ci-test-go				## Test go code
test-go-deps: ci-test-go-deps	## Test go dependencies
test-go-short: override DOCKER_RUN_OPTS += --env TEST_SHORT=1
test-go-short: ci-test-go		## Test go (short tests only)

###
# docker-related
###
# specify a special reusable volume for go-related docker builds
ci-%-go: override DOCKER_RUN_OPTS += --mount type=volume,source='gomodcache',target='/go/pkg/mod/'
ci-%-go: override DOCKER_RUN_OPTS += --mount type=volume,source='gocache',target='/root/.cache/go-build'

# use krostar/ci to test, lint, and/or build stuff
.PHONY: ci-%
ci-%: DIR_ABS := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
ci-%:
	@docker --log-level warn run							\
		--rm												\
		--mount type=bind,source="$(DIR_ABS)",target=/app	\
		$(DOCKER_RUN_OPTS)									\
		"krostar/ci:$(*)"									\
		$(DOCKER_RUN_ARGS)
