package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configYAML  string
		expectError bool
		validate    func(*testing.T, *Config)
	}{
		{
			name: "valid minimal config",
			configYAML: `
aws:
  region: us-east-1
  vpcID: vpc-12345
  subnetID: subnet-12345
  securityGroupIDs:
    - sg-12345

node:
  name: test-node
  cpu: "100"
  memory: "1Ti"
  pods: "1000"
`,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.AWS.Region != "us-east-1" {
					t.Errorf("expected region us-east-1, got %s", cfg.AWS.Region)
				}
				if cfg.Node.Name != "test-node" {
					t.Errorf("expected node name test-node, got %s", cfg.Node.Name)
				}
			},
		},
		{
			name: "valid complete config",
			configYAML: `
aws:
  region: us-west-2
  vpcID: vpc-abc123
  subnetID: subnet-def456
  securityGroupIDs:
    - sg-ghi789
  amiID: ami-jkl012
  localStackEndpoint: http://localhost:4567
  credentials:
    accessKeyID: test
    secretAccessKey: test

node:
  name: orca-test
  labels:
    environment: test
    team: research
  taints:
    - key: test
      value: "true"
      effect: NoSchedule
  operatingSystem: Linux
  cpu: "100"
  memory: "1Ti"
  pods: "1000"
  gpu: "50"

instances:
  selectionMode: auto
  templates:
    test-template:
      instanceType: t3.small
      launchType: spot
  defaultLaunchType: on-demand
  allowedInstanceTypes:
    - t3.small
    - t3.medium
  maxSpotPrices:
    t3.small: "0.05"

limits:
  maxConcurrentInstances: 10
  maxInstancesPerNamespace: 5
  dailyBudget: 100.00
  monthlyBudget: 1000.00
  maxInstanceLifetime: 2h

logging:
  level: debug
  format: json
  logAWSCalls: true

metrics:
  enabled: true
  port: 9090
  path: /metrics
  prometheus:
    enabled: true

development:
  mockAWS: false
  dryRun: false
  enableFastCleanup: true
  cleanupInterval: 1m
`,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.AWS.Region != "us-west-2" {
					t.Errorf("expected region us-west-2, got %s", cfg.AWS.Region)
				}
				if cfg.AWS.LocalStackEndpoint != "http://localhost:4567" {
					t.Errorf("expected LocalStack endpoint, got %s", cfg.AWS.LocalStackEndpoint)
				}
				if cfg.Node.CPU != "100" {
					t.Errorf("expected CPU 100, got %s", cfg.Node.CPU)
				}
				if cfg.Instances.SelectionMode != "auto" {
					t.Errorf("expected selectionMode auto, got %s", cfg.Instances.SelectionMode)
				}
				if cfg.Limits.MaxConcurrentInstances != 10 {
					t.Errorf("expected max concurrent 10, got %d", cfg.Limits.MaxConcurrentInstances)
				}
				if cfg.Logging.Level != "debug" {
					t.Errorf("expected log level debug, got %s", cfg.Logging.Level)
				}
				if cfg.Metrics.Port != 9090 {
					t.Errorf("expected metrics port 9090, got %d", cfg.Metrics.Port)
				}
			},
		},
		{
			name: "missing required AWS region",
			configYAML: `
aws:
  vpcID: vpc-12345
  subnetID: subnet-12345
  securityGroupIDs:
    - sg-12345
node:
  name: test-node
`,
			expectError: true,
		},
		{
			name: "missing required VPC ID",
			configYAML: `
aws:
  region: us-east-1
  subnetID: subnet-12345
  securityGroupIDs:
    - sg-12345
node:
  name: test-node
`,
			expectError: true,
		},
		{
			name: "missing required subnet ID",
			configYAML: `
aws:
  region: us-east-1
  vpcID: vpc-12345
  securityGroupIDs:
    - sg-12345
node:
  name: test-node
`,
			expectError: true,
		},
		{
			name: "missing security group IDs",
			configYAML: `
aws:
  region: us-east-1
  vpcID: vpc-12345
  subnetID: subnet-12345
node:
  name: test-node
`,
			expectError: true,
		},
		{
			name: "missing node name",
			configYAML: `
aws:
  region: us-east-1
  vpcID: vpc-12345
  subnetID: subnet-12345
  securityGroupIDs:
    - sg-12345
node:
  labels:
    test: value
`,
			expectError: true,
		},
		{
			name: "invalid selection mode",
			configYAML: `
aws:
  region: us-east-1
  vpcID: vpc-12345
  subnetID: subnet-12345
  securityGroupIDs:
    - sg-12345
node:
  name: test-node
instances:
  selectionMode: invalid
`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			if err := os.WriteFile(configPath, []byte(tt.configYAML), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			// Load config
			cfg, err := LoadConfig(configPath)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := &Config{
			AWS: AWSConfig{
				Region:           "us-east-1",
				VPCID:            "vpc-12345",
				SubnetID:         "subnet-12345",
				SecurityGroupIDs: []string{"sg-12345"},
			},
			Node: NodeConfig{
				Name:   "test-node",
				CPU:    "100",
				Memory: "1Ti",
				Pods:   "1000",
			},
			Instances: InstancesConfig{
				SelectionMode: "explicit",
			},
		}

		if err := cfg.Validate(); err != nil {
			t.Errorf("unexpected validation error: %v", err)
		}
	})

	t.Run("invalid AWS config - missing region", func(t *testing.T) {
		cfg := &Config{
			AWS: AWSConfig{
				VPCID:            "vpc-12345",
				SubnetID:         "subnet-12345",
				SecurityGroupIDs: []string{"sg-12345"},
			},
			Node: NodeConfig{
				Name: "test-node",
			},
		}

		if err := cfg.Validate(); err == nil {
			t.Error("expected validation error for missing region")
		}
	})

	t.Run("invalid node config - missing name", func(t *testing.T) {
		cfg := &Config{
			AWS: AWSConfig{
				Region:           "us-east-1",
				VPCID:            "vpc-12345",
				SubnetID:         "subnet-12345",
				SecurityGroupIDs: []string{"sg-12345"},
			},
			Node: NodeConfig{},
		}

		if err := cfg.Validate(); err == nil {
			t.Error("expected validation error for missing node name")
		}
	})

	t.Run("invalid selection mode", func(t *testing.T) {
		cfg := &Config{
			AWS: AWSConfig{
				Region:           "us-east-1",
				VPCID:            "vpc-12345",
				SubnetID:         "subnet-12345",
				SecurityGroupIDs: []string{"sg-12345"},
			},
			Node: NodeConfig{
				Name: "test-node",
			},
			Instances: InstancesConfig{
				SelectionMode: "invalid-mode",
			},
		}

		if err := cfg.Validate(); err == nil {
			t.Error("expected validation error for invalid selection mode")
		}
	})
}

func TestNodeCapacity(t *testing.T) {
	cfg := &NodeConfig{
		CPU:    "100",
		Memory: "1Ti",
		Pods:   "500",
		GPU:    "50",
	}

	capacity := cfg.Capacity()

	// Check CPU
	cpuQuantity := capacity["cpu"]
	if cpuQuantity.String() != "100" {
		t.Errorf("expected CPU 100, got %s", cpuQuantity.String())
	}

	// Check Memory
	memQuantity := capacity["memory"]
	if memQuantity.String() != "1Ti" {
		t.Errorf("expected Memory 1Ti, got %s", memQuantity.String())
	}

	// Check Pods
	podsQuantity := capacity["pods"]
	if podsQuantity.String() != "500" {
		t.Errorf("expected Pods 500, got %s", podsQuantity.String())
	}

	// Check GPU
	gpuQuantity := capacity["nvidia.com/gpu"]
	if gpuQuantity.String() != "50" {
		t.Errorf("expected GPU 50, got %s", gpuQuantity.String())
	}
}

func TestNodeAllocatable(t *testing.T) {
	cfg := &NodeConfig{
		CPU:    "100",
		Memory: "1Ti",
		Pods:   "500",
		GPU:    "50",
	}

	allocatable := cfg.Allocatable()

	// Check CPU
	cpuQuantity := allocatable["cpu"]
	if cpuQuantity.String() != "100" {
		t.Errorf("expected CPU 100, got %s", cpuQuantity.String())
	}

	// Check Memory
	memQuantity := allocatable["memory"]
	if memQuantity.String() != "1Ti" {
		t.Errorf("expected Memory 1Ti, got %s", memQuantity.String())
	}

	// Check Pods
	podsQuantity := allocatable["pods"]
	if podsQuantity.String() != "500" {
		t.Errorf("expected Pods 500, got %s", podsQuantity.String())
	}

	// Check GPU
	gpuQuantity := allocatable["nvidia.com/gpu"]
	if gpuQuantity.String() != "50" {
		t.Errorf("expected GPU 50, got %s", gpuQuantity.String())
	}
}

func TestSetDefaults(t *testing.T) {
	cfg := &Config{
		AWS: AWSConfig{
			Region:           "us-east-1",
			VPCID:            "vpc-12345",
			SubnetID:         "subnet-12345",
			SecurityGroupIDs: []string{"sg-12345"},
		},
		Node: NodeConfig{
			Name:   "test-node",
			CPU:    "100",
			Memory: "1Ti",
			Pods:   "1000",
		},
	}

	// Call Validate which calls setDefaults
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check defaults were set
	if cfg.Node.OperatingSystem != "Linux" {
		t.Errorf("expected default OS Linux, got %s", cfg.Node.OperatingSystem)
	}

	if cfg.Instances.SelectionMode != "explicit" {
		t.Errorf("expected default selectionMode explicit, got %s", cfg.Instances.SelectionMode)
	}

	if cfg.Instances.DefaultLaunchType != "on-demand" {
		t.Errorf("expected default launchType on-demand, got %s", cfg.Instances.DefaultLaunchType)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("expected default log level info, got %s", cfg.Logging.Level)
	}

	if cfg.Metrics.Port != 8080 {
		t.Errorf("expected default metrics port 8080, got %d", cfg.Metrics.Port)
	}
}
