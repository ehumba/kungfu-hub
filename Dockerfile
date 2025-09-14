# syntax=docker/dockerfile:1

# ---------- Builder Stage ----------
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./

# Download dependencies (cached unless go.mod/go.sum change)
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o server .

# ---------- Final Stage ----------
FROM debian:bullseye-slim

WORKDIR /app

# Copy only the compiled binary from builder stage
COPY --from=builder /app/server .

# Expose app port
EXPOSE 8080

# Start the server
CMD ["./server"]
