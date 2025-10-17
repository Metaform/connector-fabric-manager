.PHONY: test test-coverage clean lint generate-mocks install-mockery

MOCKERY_CMD=$(shell which mockery || echo $(GOPATH)/bin/mockery)

install-mockery:
	go install github.com/vektra/mockery/v2@latest

generate-mocks: install-mockery
	$(MOCKERY_CMD)

# Build settings
COVER_PROFILE=coverage.out

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=$(COVER_PROFILE) -covermode=atomic ./...
	go tool cover -html=$(COVER_PROFILE)

# Run linting
lint:
	golangci-lint run ./...

# Clean test artifacts
clean:
	rm -f $(COVER_PROFILE)
	go clean -testcache

# Run tests and linting
check: lint test