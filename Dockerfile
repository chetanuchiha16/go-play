# Step 1: Build the Go binary
FROM golang:1.25.5-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/server/v2/main.go

# Step 2: Create a tiny image to run the binary
FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env . 
EXPOSE 8080
CMD ["./main"]