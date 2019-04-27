# use this rule as the default make rule
.DEFAULT_GOAL = help
.PHONY: help
help:
	@echo "Available targets descriptions:"
	@grep -E '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint-all lint-go lint-markdown lint-yaml
lint-all: docker-lint-go docker-lint-markdown docker-lint-yaml	## Run all possible linters
lint-go: docker-lint-go											## Lint go files
lint-markdown: docker-lint-markdown								## Lint markdown files
lint-yaml: docker-lint-yaml										## Lint yaml files

.PHONY: test-all test-go test-go-deps test-go-fast
test-all: docker-test-go docker-test-go-deps 					## Run all possible tests
test-go: docker-test-go											## Test go code
test-go-deps: docker-test-go-deps								## Test go dependencies
test-go-fast: docker-test-go-fast								## Fast test co code

.PHONY: docker-%
# specify a special reusable volume for go-related docker builds
docker-%-go: DOCKER_RUN_OPTS += --mount type=volume,source='gomodcache',target='/go/pkg/mod/'
docker-%: DIR_ABS := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
docker-%:
	@docker pull "krostar/ci:$*" 2>&1 > /dev/null
	@docker --log-level warn run \
		--rm \
		--tty \
		--mount type=bind,source="$(DIR_ABS)",target=/app \
		$(DOCKER_RUN_OPTS) \
		"krostar/ci:$*" \
		$(DOCKER_RUN_ARGS)
