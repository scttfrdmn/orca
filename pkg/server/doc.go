// Package server provides HTTP endpoints for ORCA.
//
// The server exposes:
// - Health check endpoints (/healthz, /readyz) for Kubernetes probes
// - Prometheus metrics endpoint (/metrics) for monitoring
//
// The health check endpoints follow Kubernetes conventions:
// - /healthz: Liveness probe - returns 200 if server is running
// - /readyz: Readiness probe - returns 200 if ready to serve traffic, 503 otherwise
package server
