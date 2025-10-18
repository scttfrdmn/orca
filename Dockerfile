# ORCA Dockerfile
# Multi-stage build for optimal image size

# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN make build

# Stage 2: Runtime
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S orca && adduser -S orca -G orca

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/bin/orca /app/orca

# Copy default configuration
COPY deploy/kubernetes/config.yaml /app/config.yaml

# Set ownership
RUN chown -R orca:orca /app

# Switch to non-root user
USER orca

# Expose metrics port
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/app/orca"]
CMD ["--config", "/app/config.yaml"]
