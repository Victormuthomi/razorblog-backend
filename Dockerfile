# 1. Build stage
FROM golang:1.24.3-alpine AS builder


WORKDIR /app

# Install git (needed if you use go modules with private repos)
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN go build -o razorblog ./cmd/main.go

# 2. Final stage
FROM alpine:3.18

WORKDIR /app

# Optional: add CA certs if you make HTTPS requests
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/razorblog .

# Expose port (match your Gin server)
EXPOSE 8080

# Run the binary
CMD ["./razorblog"]

