# Use an official Golang image as a base for building
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for dependency caching
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application binary
RUN go build -o scrambled-strings ./cmd

# Use Ubuntu 22.04 as the final stage
FROM ubuntu:22.04

# Install required libraries
RUN apt-get update && apt-get install -y libc6 && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/scrambled-strings /app/scrambled-strings

# Set the entrypoint to the CLI tool
ENTRYPOINT ["/app/scrambled-strings"]
