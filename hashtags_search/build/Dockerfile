# Use the official Golang image for building
FROM golang:1.23 AS builder

# Set environment variables for Go
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create an app directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project source
COPY . .

# Build the binary
RUN go build -o searchhashtags ./search_hashtags.go

# Use a minimal image for running the app
FROM alpine:latest

# Set environment variables
ENV GIN_MODE=release

# Install required dependencies
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/searchhashtags .

# Expose the application port
EXPOSE 8380

# Run the application
CMD ["./searchhashtags"]
