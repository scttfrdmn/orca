package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/scttfrdmn/orca/pkg/config"
	"github.com/scttfrdmn/orca/pkg/node"
	"github.com/scttfrdmn/orca/pkg/server"
)

var (
	// Version information (set by build flags)
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Parse command-line flags
	var (
		configFile  = flag.String("config", "config.yaml", "path to config file")
		kubeconfig  = flag.String("kubeconfig", "", "path to kubeconfig file (uses in-cluster config if empty)")
		nodeName    = flag.String("node-name", "", "name of the virtual node (overrides config)")
		namespace   = flag.String("namespace", "kube-system", "namespace to run in")
		showVersion = flag.Bool("version", false, "show version information")
		logLevel    = flag.String("log-level", "", "log level (overrides config)")
	)
	flag.Parse()

	// Show version and exit
	if *showVersion {
		fmt.Printf("ORCA version %s\n", version)
		fmt.Printf("  Git commit: %s\n", gitCommit)
		fmt.Printf("  Built:      %s\n", buildDate)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Override config with command-line flags
	if *nodeName != "" {
		cfg.Node.Name = *nodeName
	}
	if *logLevel != "" {
		cfg.Logging.Level = *logLevel
	}

	// Setup logging
	logger := setupLogging(cfg.Logging)

	// Log startup information
	logger.Info().
		Str("version", version).
		Str("git_commit", gitCommit).
		Str("build_date", buildDate).
		Msg("Starting ORCA")

	logger.Info().
		Str("node_name", cfg.Node.Name).
		Str("namespace", *namespace).
		Str("aws_region", cfg.AWS.Region).
		Str("vpc_id", cfg.AWS.VPCID).
		Str("subnet_id", cfg.AWS.SubnetID).
		Msg("Configuration loaded")

	// Check for LocalStack endpoint
	if cfg.AWS.LocalStackEndpoint != "" {
		logger.Warn().
			Str("endpoint", cfg.AWS.LocalStackEndpoint).
			Msg("Using LocalStack for development")
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
		cancel()
	}()

	// Create HTTP server for health checks and metrics
	logger.Info().Msg("Creating HTTP server")
	httpServer := server.NewServer(cfg, logger)

	// Start HTTP server in background
	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logger.Error().Err(err).Msg("HTTP server error")
		}
	}()

	// Create node controller
	logger.Info().Msg("Creating ORCA node controller")
	controller, err := node.NewController(cfg, *kubeconfig, *namespace, version, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create node controller")
	}

	// Start the controller in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := controller.Run(ctx); err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	// Mark server as ready once node is registered
	// (Virtual Kubelet node controller returns after registration)
	httpServer.SetReady(true)

	logger.Info().
		Int("http_port", cfg.Metrics.Port).
		Msg("ORCA is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal or error
	select {
	case <-ctx.Done():
		logger.Info().Msg("Shutting down...")
	case err := <-errChan:
		if err != nil {
			logger.Error().Err(err).Msg("Controller error")
		}
	}

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server first (stop accepting new requests)
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("HTTP server shutdown error")
	}

	// Then shutdown node controller
	if err := controller.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Node controller shutdown error")
	}

	logger.Info().Msg("ORCA shutdown complete")
}

// setupLogging configures zerolog based on configuration.
func setupLogging(cfg config.LoggingConfig) zerolog.Logger {
	// Set global log level
	switch cfg.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Configure output format
	var logger zerolog.Logger
	if cfg.Format == "json" {
		// JSON output (default)
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Human-readable console output
		logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	// Set as global logger
	log.Logger = logger

	return logger
}
