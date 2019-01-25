# ex: /Users/alice/go/src/github.com/alice/project/
DIR_ABS    := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
DIR_SCRIPT := $(DIR_ABS)/scripts
DIR_BUILD  := $(DIR_ABS)/build
DIR_BIN    := $(DIR_BUILD)/bin
DIR_COVER  := $(DIR_BUILD)/cover

# use this rule as the default make rule
.DEFAULT_GOAL := help
.PHONY: help
help:
	@echo "Available targets descriptions:"
	@# absolutely awesome -> http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
	@grep -E '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test-% lint-%
test-%: SCRIPT := test.sh ## Test project, possible values are [dep, go, go-short, all]
lint-%: SCRIPT := lint.sh ## Lint files, possible values are [go, sh, yaml, all]
test-% lint-%: | $(DIR_COVER)
	@$(DIR_SCRIPT)/$(SCRIPT) "$*"

.PHONY: cover
cover: $(DIR_COVER)/coverage.out ## Display the coverage report
	@go tool cover -html=$(DIR_COVER)/coverage.out
$(DIR_COVER)/coverage.out: test-go

# create directory in case they doesn't exists
$(DIR_COVER):
	@mkdir -p $@
