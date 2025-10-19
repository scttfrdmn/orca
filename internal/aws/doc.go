// Package aws provides AWS EC2 client functionality for ORCA.
//
// This package wraps the AWS SDK v2 EC2 client and provides ORCA-specific
// operations for managing EC2 instances, including:
// - Creating instances with proper tags and configuration
// - Terminating instances
// - Querying instance state
// - Support for both on-demand and spot instances
//
// The client automatically applies ORCA resource tags to all created instances
// for proper resource tracking and cost attribution.
package aws
