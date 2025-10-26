# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Generate swagger docs for production
RUN go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /build/src
RUN sed -e 's/@title Hornetsecurity-Document Store API/@title Hornet Security-Document Store API - Production/' \
    -e 's/@description A simple document storage API with CRUD operations/@description A simple document storage API with CRUD operations - Production Environment/' \
    main.go > main_temp.go && \
    swag init -g main_temp.go --output docs --instanceName prod && \
    rm main_temp.go

# Build the binary with optimizations
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app ./src

# Final stage - minimal image with basic filesystem
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create appuser for security
RUN adduser -D -g '' appuser

# Create necessary directories
RUN mkdir -p /home/appuser && chown appuser:appuser /home/appuser

# Copy the binary
COPY --from=builder /build/app /app

# Copy config directory for environment file loading
COPY --from=builder /build/config /config

# Use non-root user
USER appuser

# Set HOME directory
ENV HOME=/home/appuser

# Expose port
EXPOSE 8080

# Health check (removed problematic health-check flag)

# Run the binary
ENTRYPOINT ["/app"]
