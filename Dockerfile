# Step 1: Build the Go binary
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/server/v2/main.go

# Step 2: Create a tiny image to run the binary
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env . 
EXPOSE 8080
CMD ["./main"]