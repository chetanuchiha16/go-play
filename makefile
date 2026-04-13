# Variables
BINARY_NAME=bin/server
MAIN_PATH=cmd/server/main.go

.PHONY: dev build run sqlc clean up down generate

# Start local development with hot reloading (Air)
dev:
	air

# Build the binary
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the built binary locally
run: build
	./$(BINARY_NAME)

# Generate type-safe Go code from your SQL
sqlc:
	sqlc generate

# Generate oapi-codegen server stubs from OpenAPI spec
generate:
	oapi-codegen -generate types,std-http,spec -package api -o internal/api/generated.go api/openapi.yaml

# --- DOCKER COMMANDS ---

# Start everything (DB + App) in containers
up:
	docker compose up --build

# Stop everything and remove containers
down:
	docker compose down

# Clean up binaries
clean:
	rm -rf bin/

mock:
	mockery

# Run all tests in the project (excludes integration tests)
test:
	go test -v ./...

# Run integration tests (requires server on :9999 and PostgreSQL)
test-integration:
	go test -tags=integration -v ./tests/integration/

# Ensure all dependencies are in sync
tidy:
	go mod tidy

# Generate test coverage report
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out