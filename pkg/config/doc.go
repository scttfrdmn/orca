// Package config provides configuration management for ORCA.
//
// It supports loading configuration from YAML files and provides
// validation and default values for all configuration options.
//
// Example usage:
//
//	cfg, err := config.LoadConfig("config.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Configuration includes AWS settings, node capacity, instance selection
// modes, resource limits, logging, and metrics.
package config
