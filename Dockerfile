# ORCA Dockerfile
# Multi-stage build for minimal, secure container image

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates

WORKDIR /build

# Copy go mod files first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
ARG VERSION=dev
ARG GIT_COMMIT=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -trimpath \
    -o /orca \
    ./cmd/orca

# Final stage - distroless for minimal attack surface
FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.title="ORCA" \
      org.opencontainers.image.description="Orchestration for Research Cloud Access - Kubernetes Virtual Kubelet for AWS" \
      org.opencontainers.image.url="https://github.com/scttfrdmn/orca" \
      org.opencontainers.image.source="https://github.com/scttfrdmn/orca" \
      org.opencontainers.image.vendor="Scott Friedman" \
      org.opencontainers.image.licenses="Apache-2.0"

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /orca /orca

# Use nonroot user (UID 65532)
USER nonroot:nonroot

# Expose metrics port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/orca"]
