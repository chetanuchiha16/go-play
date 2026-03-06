# Step 1: Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
# Using your path from the makefile
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/server/v2/main.go

# Step 2: Run Stage
FROM alpine:latest
WORKDIR /root/
# Copy the binary from the exact path used in the builder
COPY --from=builder /app/bin/server .
COPY .env . 
EXPOSE 8080
# Run the renamed binary
CMD ["./server"]