cd.PHONY: build test clean run server dev-server

# Binary name
SERVER_BINARY=pmanager
TEST_AGENT_BINARY=pmtestagent

# Build settings
BUILD_DIR=bin
MAIN_PATH=./cmd/main.go
SERVER_PATH=./cmd/server/main.go
TEST_AGENT_PATH=./agent/testagent/main.go

MOCKERY_CMD=$(GOPATH)/bin/mockery

# Environment variables
export GO111MODULE=on
export CGO_ENABLED=0

# Install development tools
install-tools: install-mockery

install-mockery:
	go get github.com/vektra/mockery/v2@latest

# Generate mocks
generate-mocks: install-mockery
	$(MOCKERY_CMD) --case underscore --dir ./natsorchestration --all --output ./natsorchestration/mocks --outpkg mocks --with-expecter


# Build the application
build: build-server build-test-agent

# Build the server
build-server:
	go build -o $(BUILD_DIR)/$(SERVER_BINARY) $(SERVER_PATH)

build-test-agent:
	go build -o $(BUILD_DIR)/$(TEST_AGENT_BINARY) $(TEST_AGENT_PATH)

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Run the server in development mode
dev-server: build-server
	./$(BUILD_DIR)/$(SERVER_BINARY)

dev-test-agent: build-server
	./$(BUILD_DIR)/$(TEST_AGENT_BINARY)

# Run the server in production mode
server: build-server
	./$(BUILD_DIR)/$(SERVER_BINARY)

# Build for multiple platforms
build-all:
	# Server binaries
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(SERVER_BINARY)-linux-amd64 $(SERVER_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(SERVER_BINARY)-darwin-amd64 $(SERVER_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(SERVER_BINARY)-darwin-arm64 $(SERVER_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(SERVER_BINARY)-windows-amd64.exe $(SERVER_PATH)
	
