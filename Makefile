.PHONY: proto
proto:
	@buf generate
.PHONY: lint
lint:
	@golangci-lint run
.PHONY: test
test:
	@gotestsum $(shell go list ./... | grep -v gen | grep -v cmd)

.PHONY: test-with-coverage
test-with-coverage:
	@gotestsum $(shell go list ./... | grep -v gen | grep -v cmd) -coverprofile=cover.out

.PHONY: show-coverage
coverage:
	go tool cover -html=cover.out
