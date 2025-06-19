FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /app

# Install curl and other dependencies
RUN apk --no-cache add curl ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/
COPY --from=builder /app/.env /app/

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
