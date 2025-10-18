package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Config is the main configuration structure for ORCA.
type Config struct {
	AWS         AWSConfig         `yaml:"aws"`
	Node        NodeConfig        `yaml:"node"`
	Instances   InstancesConfig   `yaml:"instances"`
	Limits      LimitsConfig      `yaml:"limits"`
	Logging     LoggingConfig     `yaml:"logging"`
	Metrics     MetricsConfig     `yaml:"metrics"`
	Development DevelopmentConfig `yaml:"development"`
}

// AWSConfig contains AWS-specific configuration.
type AWSConfig struct {
	Region           string          `yaml:"region"`
	Credentials      *AWSCredentials `yaml:"credentials,omitempty"`
	VPCID            string          `yaml:"vpcID"`
	SubnetID         string          `yaml:"subnetID"`
	SecurityGroupIDs []string        `yaml:"securityGroupIDs"`
	AMIID            string          `yaml:"amiID,omitempty"`
	DevelopmentMode  bool            `yaml:"developmentMode"`
}

// AWSCredentials contains AWS access credentials.
type AWSCredentials struct {
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
}

// NodeConfig contains virtual node configuration.
type NodeConfig struct {
	Name            string            `yaml:"name"`
	Labels          map[string]string `yaml:"labels"`
	Taints          []corev1.Taint    `yaml:"taints"`
	OperatingSystem string            `yaml:"operatingSystem"`
	CPU             string            `yaml:"cpu"`
	Memory          string            `yaml:"memory"`
	Pods            string            `yaml:"pods"`
	GPU             string            `yaml:"gpu,omitempty"`
}

// Capacity returns the node capacity as a ResourceList.
func (n *NodeConfig) Capacity() corev1.ResourceList {
	capacity := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(n.CPU),
		corev1.ResourceMemory: resource.MustParse(n.Memory),
		corev1.ResourcePods:   resource.MustParse(n.Pods),
	}

	if n.GPU != "" {
		capacity["nvidia.com/gpu"] = resource.MustParse(n.GPU)
	}

	return capacity
}

// Allocatable returns the allocatable resources (same as capacity for virtual node).
func (n *NodeConfig) Allocatable() corev1.ResourceList {
	return n.Capacity()
}

// InstancesConfig contains instance selection configuration.
type InstancesConfig struct {
	SelectionMode        string                      `yaml:"selectionMode"`
	Templates            map[string]WorkloadTemplate `yaml:"templates"`
	DefaultLaunchType    string                      `yaml:"defaultLaunchType"`
	AllowedInstanceTypes []string                    `yaml:"allowedInstanceTypes"`
	MaxSpotPrices        map[string]string           `yaml:"maxSpotPrices"`
}

// WorkloadTemplate defines a template for common workloads.
type WorkloadTemplate struct {
	InstanceType string `yaml:"instanceType"`
	LaunchType   string `yaml:"launchType"`
	MaxSpotPrice string `yaml:"maxSpotPrice,omitempty"`
}

// LimitsConfig contains resource limits and budget controls.
type LimitsConfig struct {
	MaxConcurrentInstances   int                       `yaml:"maxConcurrentInstances"`
	MaxInstancesPerNamespace int                       `yaml:"maxInstancesPerNamespace"`
	DailyBudget              float64                   `yaml:"dailyBudget"`
	MonthlyBudget            float64                   `yaml:"monthlyBudget"`
	MaxInstanceLifetime      time.Duration             `yaml:"maxInstanceLifetime"`
	NamespaceQuotas          map[string]NamespaceQuota `yaml:"namespaceQuotas"`
}

// NamespaceQuota defines limits for a specific namespace.
type NamespaceQuota struct {
	MaxInstances int     `yaml:"maxInstances"`
	MaxGPUs      int     `yaml:"maxGPUs"`
	DailyBudget  float64 `yaml:"dailyBudget"`
}

// LoggingConfig contains logging configuration.
type LoggingConfig struct {
	Level       string `yaml:"level"`
	Format      string `yaml:"format"`
	LogAWSCalls bool   `yaml:"logAWSCalls"`
}

// MetricsConfig contains metrics configuration.
type MetricsConfig struct {
	Enabled    bool             `yaml:"enabled"`
	Port       int              `yaml:"port"`
	Path       string           `yaml:"path"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
}

// PrometheusConfig contains Prometheus-specific configuration.
type PrometheusConfig struct {
	Enabled bool `yaml:"enabled"`
}

// DevelopmentConfig contains development-specific settings.
type DevelopmentConfig struct {
	MockAWS           bool          `yaml:"mockAWS"`
	DryRun            bool          `yaml:"dryRun"`
	EnableFastCleanup bool          `yaml:"enableFastCleanup"`
	CleanupInterval   time.Duration `yaml:"cleanupInterval"`
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	// Validate AWS config
	if c.AWS.Region == "" {
		return fmt.Errorf("aws.region is required")
	}

	// Validate node config
	if c.Node.Name == "" {
		return fmt.Errorf("node.name is required")
	}
	if c.Node.OperatingSystem == "" {
		c.Node.OperatingSystem = "Linux"
	}
	if c.Node.CPU == "" {
		return fmt.Errorf("node.cpu is required")
	}
	if c.Node.Memory == "" {
		return fmt.Errorf("node.memory is required")
	}
	if c.Node.Pods == "" {
		return fmt.Errorf("node.pods is required")
	}

	// Validate instances config
	if c.Instances.SelectionMode == "" {
		c.Instances.SelectionMode = "explicit"
	}
	if c.Instances.SelectionMode != "explicit" && c.Instances.SelectionMode != "template" && c.Instances.SelectionMode != "auto" {
		return fmt.Errorf("instances.selectionMode must be explicit, template, or auto")
	}
	if c.Instances.DefaultLaunchType == "" {
		c.Instances.DefaultLaunchType = "on-demand"
	}
	if c.Instances.DefaultLaunchType != "on-demand" && c.Instances.DefaultLaunchType != "spot" {
		return fmt.Errorf("instances.defaultLaunchType must be on-demand or spot")
	}

	// Validate logging config
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}

	// Validate metrics config
	if c.Metrics.Port == 0 {
		c.Metrics.Port = 8080
	}
	if c.Metrics.Path == "" {
		c.Metrics.Path = "/metrics"
	}

	return nil
}
