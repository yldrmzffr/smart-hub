# Tool versions
BUF_VERSION := v1.28.1
PROTOC_GO_VERSION := v1.31.0
PROTOC_GO_GRPC_VERSION := v1.3.0

# Build-time variables
DOCKER_COMPOSE_TEST_FILE := docker-compose.test.yml

.PHONY: run
run:
	@echo "Running the server..."
	go run cmd/api/main.go


.PHONY: prepare-and-run
prepare-and-run: setup proto run


.PHONY: all
all: setup proto test-integration

# Install required development tools
.PHONY: setup
setup:
	@echo "Installing required development tools..."
	@go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GO_VERSION)
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GO_GRPC_VERSION)
	@echo "All tools installed successfully"

.PHONY: proto
proto: setup
	@echo "Generating protobuf code..."
	@buf generate
	@echo "Protobuf code generation completed"

.PHONY: test-integration-up
test-integration-up:
	@echo "Starting test containers..."
	docker-compose -f $(DOCKER_COMPOSE_TEST_FILE) up -d
	@echo "Waiting for containers to be ready..."
	@sleep 5

.PHONY: test-integration-down
test-integration-down:
	@echo "Stopping test containers..."
	docker-compose -f $(DOCKER_COMPOSE_TEST_FILE) down
	@echo "Test containers stopped"

.PHONY: test-integration
test-integration: test-integration-up
	@echo "Running integration tests..."
	go test ./tests/integration/... -v -count=1
	@make test-integration-down

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	go test ./... -v -count=1

.PHONY: all
all: setup proto test-unit test-integration

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	rm -rf gen/
	rm -f coverage.out
	@echo "Stopping any running test containers..."
	@make test-integration-down
	@echo "Cleanup completed"

.PHONY: coverage
coverage: test-integration-up
	@echo "Running test coverage..."
	go test $(go list ./... | grep -v '/gen/proto') -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  setup              - Install required development tools"
	@echo "  proto              - Generate protobuf code"
	@echo "  run                - Run the server"
	@echo "  prepare-and-run    - Setup, proto, run"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-unit          - Run unit tests"
	@echo "  all                - Setup, proto, test-integration"
	@echo "  coverage           - Run test coverage"
	@echo "  clean              - Clean generated files and stop test containers"
	@echo "  help               - Show this help message"