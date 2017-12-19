GOTOOLS = \
	github.com/alecthomas/gometalinter \
	golang.org/x/tools/cmd/cover \
	github.com/golang/dep/cmd/dep \

GOPACKAGES := $(go list ./... | grep -v /vendor/)

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get -u $(GOTOOLS)
	gometalinter --install --update

.PHONY: dep
dep: ## Run dep ensure and prune
	dep ensure
	dep prune

.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && go test -covermode=atomic -coverprofile=coverage.txt -v -timeout=30s $(GOPACKAGES)

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

.PHONY: fmt
fmt: ## goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -s -w "$$file"; goimports -w "$$file"; done

.PHONY: lint
lint: ## Run all the linters
	gometalinter --exclude=vendor --exclude=/go/src --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gofmt \
		--enable=goimports \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=megacheck \
		--deadline=10m \
		$(GOPACKAGES)

.PHONY: ci
ci: lint test ## Run all the tests and code checks

.PHONY: build
build: ## Build a version
	CGO_ENABLED=0 go build $(GOPACKAGES)

.PHONY: clean
clean: ## Remove temporary files
	go clean

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
