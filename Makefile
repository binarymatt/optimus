.PHONY: proto
proto:
	@buf generate
.PHONY: lint
lint:
	@golangci-lint run
.PHONY: test
test:
	@gotestsum

.PHONY: test-with-coverage
test-with-coverage:
	@gotestsum -- -coverprofile=cover.out ./...
