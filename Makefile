.PHONY: proto
proto:
	@buf generate
.PHONY: lint
lint:
	@golangci-lint run
.PHONY: test
test:
	@gotestsum -- -coverprofile=cover.out github.com/binarymatt/optimus
