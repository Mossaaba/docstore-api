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

# Generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app .

# Final stage - minimal image
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/app /app

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check (removed problematic health-check flag)

# Run the binary
ENTRYPOINT ["/app"]