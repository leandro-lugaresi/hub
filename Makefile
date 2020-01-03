SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)
TEST_PATTERN?=./...
TEST_OPTIONS?=-race

setup: ## Install all the build and lint dependencies
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.22.2
	GO111MODULE=off go get github.com/mfridman/tparse
	GO111MODULE=off go get -u golang.org/x/tools/cmd/cover
	go get -v -t ./...

test: ## Run all the tests
	go test $(TEST_OPTIONS) -covermode=atomic -coverprofile=coverage.txt -timeout=1m -cover -json $(SOURCE_FILES) | $$(go env GOPATH)/bin/tparse -all

throughput: ## Run the throughput tests
	go test -v -timeout 60s github.com/leandro-lugaresi/hub -run ^TestThroughput -args -throughput

bench: ## Run the benchmark tests
	go test -bench=. $(TEST_PATTERN)

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run all the linters
	$$(go env GOPATH)/bin/golangci-lint run

ci: lint test ## Run all the tests and code checks

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
