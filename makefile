# Variables
BINARY_NAME=bin/server
MAIN_PATH=cmd/server/v2/main.go

.PHONY: dev build run sqlc clean

# Start development with hot reloading (Air)
dev:
	air

# Build the binary
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the built binary
run: build
	./$(BINARY_NAME)

# Generate SQLC code
sqlc:
	sqlc generate

# Clean up binaries
clean:
	rm -rf bin/