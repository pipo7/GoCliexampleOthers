PROJECT_NAME := "perftest"  ## Or if its current directory then No need for this.
## Use this makefile in gitlab for exceution of jobs and stages.

.PHONY: all
all: compile

FORCE:

.PHONY: compile
compile: ## Compile go project
	echo "grant executable permission to script run_in_cosmobox.sh"
	cd $(PROJECT_NAME) && go build -o . ./...

.PHONY: clean
clean:

.PHONY: release
release:

.PHONY: test
test:  ## Execute tests
	cd $(PROJECT_NAME) && go test ./...

.PHONY: coverage
coverage:  ## Generate global code coverage report
	cd $(PROJECT_NAME) && echo "run covereage shell script if any" ;

.PHONY: lint
lint: ## Lint the files
	cd $(PROJECT_NAME) && golangci-lint run --timeout 5m

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
