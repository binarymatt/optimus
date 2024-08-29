.PHONY: proto
proto:
	@buf generate
.PHONY: lint
lint:
	@golangci-lint run
.PHONY: test
test:
	@gotestsum $(shell go list ./... | grep -v gen | grep -v cmd | grep -v mocks)

.PHONY: test-with-coverage
test-with-coverage:
	@gotestsum $(shell go list ./... | grep -v gen | grep -v cmd | grep -v mocks) -coverprofile=cover.out

.PHONY: show-coverage
coverage:
	go tool cover -html=cover.out
