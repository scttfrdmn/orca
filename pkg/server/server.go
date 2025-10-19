package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/scttfrdmn/orca/pkg/config"
)

// Server provides HTTP endpoints for health checks and metrics.
type Server struct {
	config     *config.Config
	httpServer *http.Server
	logger     zerolog.Logger
	ready      bool
}

// NewServer creates a new HTTP server.
func NewServer(cfg *config.Config, logger zerolog.Logger) *Server {
	mux := http.NewServeMux()

	s := &Server{
		config: cfg,
		logger: logger,
		ready:  false,
	}

	// Register handlers
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/readyz", s.handleReadyz)
	
	// Register metrics endpoint if enabled
	if cfg.Metrics.Enabled {
		mux.Handle(cfg.Metrics.Path, promhttp.Handler())
	}

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Metrics.Port)
	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info().
		Int("port", s.config.Metrics.Port).
		Bool("metrics_enabled", s.config.Metrics.Enabled).
		Str("metrics_path", s.config.Metrics.Path).
		Msg("Starting HTTP server")

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context cancellation or error
	select {
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	case err := <-errChan:
		return fmt.Errorf("HTTP server error: %w", err)
	}
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("Shutting down HTTP server")
	
	// Set readiness to false
	s.ready = false

	// Shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	return nil
}

// SetReady marks the server as ready to serve traffic.
func (s *Server) SetReady(ready bool) {
	s.ready = ready
	s.logger.Info().Bool("ready", ready).Msg("Readiness status changed")
}

// handleHealthz handles liveness probe requests.
// Returns 200 if the server is running (always succeeds unless server is down).
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","service":"orca"}`)
	
	s.logger.Debug().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("remote_addr", r.RemoteAddr).
		Int("status", http.StatusOK).
		Msg("Health check")
}

// handleReadyz handles readiness probe requests.
// Returns 200 if the server is ready to serve traffic, 503 otherwise.
func (s *Server) handleReadyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if s.ready {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready","service":"orca"}`)
		
		s.logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status", http.StatusOK).
			Msg("Readiness check - ready")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `{"status":"not_ready","service":"orca"}`)
		
		s.logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status", http.StatusServiceUnavailable).
			Msg("Readiness check - not ready")
	}
}
