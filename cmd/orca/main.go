package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/scttfrdmn/orca/pkg/config"
	"github.com/scttfrdmn/orca/pkg/provider"
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
		nodeName    = flag.String("node-name", "", "name of the virtual node (overrides config)")
		namespace   = flag.String("namespace", "kube-system", "namespace to run in")
		showVersion = flag.Bool("version", false, "show version information")
		logLevel    = flag.String("log-level", "", "log level (overrides config)")
	)
	flag.Parse()

	// TODO: Add kubeconfig flag when Virtual Kubelet integration is implemented
	// kubeconfig := flag.String("kubeconfig", "", "path to kubeconfig file")

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
	setupLogging(cfg.Logging)

	// Log startup information
	log(cfg.Logging.Level, "info", fmt.Sprintf("Starting ORCA version %s", version))
	log(cfg.Logging.Level, "info", fmt.Sprintf("Node name: %s", cfg.Node.Name))
	log(cfg.Logging.Level, "info", fmt.Sprintf("Namespace: %s", *namespace))
	log(cfg.Logging.Level, "info", fmt.Sprintf("AWS region: %s", cfg.AWS.Region))

	// Check for LocalStack endpoint
	if cfg.AWS.LocalStackEndpoint != "" {
		log(cfg.Logging.Level, "warn", fmt.Sprintf("Using LocalStack endpoint: %s", cfg.AWS.LocalStackEndpoint))
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log(cfg.Logging.Level, "info", fmt.Sprintf("Received signal %s, shutting down...", sig))
		cancel()
	}()

	// Create provider
	log(cfg.Logging.Level, "info", "Creating ORCA provider...")
	_, err = provider.NewProvider(cfg, cfg.Node.Name, *namespace, version)
	if err != nil {
		log(cfg.Logging.Level, "error", fmt.Sprintf("Failed to create provider: %v", err))
		os.Exit(1)
	}

	// TODO: Start Virtual Kubelet node controller
	// TODO: Register with Kubernetes API server using provider
	// TODO: Start pod watch loop
	// TODO: Start metrics server

	log(cfg.Logging.Level, "info", "ORCA provider created successfully")
	log(cfg.Logging.Level, "info", "Virtual node registered: "+cfg.Node.Name)

	// For now, just wait for shutdown signal
	log(cfg.Logging.Level, "info", "ORCA is running. Press Ctrl+C to stop.")

	<-ctx.Done()

	log(cfg.Logging.Level, "info", "ORCA shutdown complete")
}

// setupLogging configures logging based on configuration
func setupLogging(cfg config.LoggingConfig) {
	// TODO: Setup structured logging (logrus, zap, or zerolog)
	// For now, using simple stdout logging
}

// log is a simple logging function
// TODO: Replace with proper structured logging
func log(level, severity, message string) {
	// Simple severity filtering
	severities := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	configLevel := severities[level]
	msgLevel := severities[severity]

	if msgLevel >= configLevel {
		fmt.Printf("[%s] %s\n", severity, message)
	}
}
