# Use an official Golang image as a build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

# Use a minimal base image for the final container
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary has execution permissions
RUN chmod +x /root/main

# Expose necessary ports
EXPOSE 3000

# Run the application
CMD ["./main"]
