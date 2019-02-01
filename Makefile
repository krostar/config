# ex: /Users/alice/go/src/github.com/alice/project/
DIR_ABS   := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
DIR_BUILD := $(DIR_ABS)/build

# use this rule as the default make rule
.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "Available targets descriptions:"
	@# absolutely awesome -> http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint-all lint-dockerfile lint-go lint-markdown lint-sh lint-yaml
lint-all: docker-lint-dockerfile docker-lint-go docker-lint-markdown docker-lint-sh docker-lint-yaml ## Run all possible linters
lint-dockerfile: docker-lint-dockerfile	## Lint dockerfiles
lint-go: docker-lint-go					## Lint go files
lint-markdown: docker-lint-markdown		## Lint markdown files
lint-sh: docker-lint-sh					## Lint shell files
lint-yaml: docker-lint-yaml				## Lint yaml files

.PHONY: test-all test-go test-go-deps test-go-fast
test-all: docker-test-go docker-test-go-deps	## Run all possible tests
test-go: docker-test-go							## Test go code
test-go-deps: docker-test-go-deps				## Test go dependencies
test-go-fast: docker-test-go-fast				## Fast test co code

docker-build-%:
	@docker build \
		--tag "configue:$*" \
		--file "$(DIR_BUILD)/$(*).Dockerfile" \
		.

.PHONY: docker-lint-% docker-test-%
# specify a special reusable volume for go-related docker builds
docker-%-go: DOCKER_RUN_OPTS += --mount type=volume,source='gomodcache',target='/go/pkg/mod/'
docker-%: docker-build-%
	@docker run \
		--rm \
		--tty \
		--mount type=bind,source="$(DIR_ABS)",target=/app,readonly \
		--mount type=bind,source="$(DIR_ABS)/build",target=/app/build \
		$(DOCKER_RUN_OPTS) \
		"configue:$*" \
		$(DOCKER_RUN_ARGS)
