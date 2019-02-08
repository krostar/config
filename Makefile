# ex: /Users/alice/go/src/github.com/alice/project/
DIR_ABS := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

# use this rule as the default make rule
.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "Available targets descriptions:"
	@# absolutely awesome -> http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint-all test-all
lint-all: docker-lint-go docker-lint-markdown docker-lint-yaml ## Run all possible linters
test-all: docker-test-go docker-test-go-deps	## Run all possible tests

.PHONY: lint-go lint-markdown lint-yaml
lint-go: docker-lint-go					## Lint go files
lint-markdown: docker-lint-markdown		## Lint markdown files
lint-yaml: docker-lint-yaml				## Lint yaml files

.PHONY: test-go test-go-deps test-go-fast
test-go: docker-test-go							## Test go code
test-go-deps: docker-test-go-deps				## Test go dependencies
test-go-fast: docker-test-go-fast				## Fast test co code

.PHONY: docker-lint-% docker-test-%
# specify a special reusable volume for go-related docker builds
docker-%-go: DOCKER_RUN_OPTS += --mount type=volume,source='gomodcache',target='/go/pkg/mod/'
docker-%:
	@docker --log-level warn run \
		--rm \
		--tty \
		--mount type=bind,source="$(DIR_ABS)",target=/app \
		$(DOCKER_RUN_OPTS) \
		"krostar/ci:$*" \
		$(DOCKER_RUN_ARGS)
