# Variables
BINARY_NAME=bin/server
MAIN_PATH=cmd/server/v2/main.go

.PHONY: dev build run sqlc clean up down

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