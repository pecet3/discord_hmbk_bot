# Start from the official Go image
FROM golang:1.22.1-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy static files and .env
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env .

# Load environment variables
COPY .env* ./ 

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]