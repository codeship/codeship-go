GOTOOLS = \
	golang.org/x/tools/cmd/cover \
	github.com/golangci/golangci-lint/cmd/golangci-lint@de1d1ad \
	github.com/mattn/goveralls \

GOPACKAGES := $(go list ./... | grep -v /vendor/)

export GOBIN:=$(PWD)/bin
export PATH:=$(GOBIN):$(PATH)

.PHONY: setup
setup: ## Install all dependencies
	go get -v -t ./...

.PHONY: tools
tools: ## Install external tools
	go get -v $(GOTOOLS)

.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -timeout=30s $(GOPACKAGES)

.PHONY: integration
integration: ## Run integration tests
	go test -v -tags=integration ./integration/...

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

.PHONY: fmt
fmt: ## goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -s -w "$$file"; goimports -w "$$file"; done

.PHONY: lint
lint: ## Run all the linters
	golangci-lint run

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
