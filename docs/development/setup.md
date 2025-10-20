# ORCA Development Guide

This guide covers local development setup, testing, and contribution workflows for ORCA.

## Prerequisites

### Required Tools

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)
- **Docker** - [Install Docker](https://docs.docker.com/get-docker/)
- **kubectl** - [Install kubectl](https://kubernetes.io/docs/tasks/tools/)
- **kind** (optional) - For local Kubernetes testing
- **golangci-lint** - For code linting

### Install Development Tools

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install kind (Kubernetes in Docker) - optional
go install sigs.k8s.io/kind@latest
```

## Getting Started

### 1. Clone and Setup

```bash
# Clone the repository
git clone https://github.com/scttfrdmn/orca.git
cd orca

# Download dependencies
make mod-download

# Verify setup
make test
```

### 2. Project Structure

```
orca/
├── cmd/orca/              # Main application entry point
├── pkg/                   # Public packages
│   ├── provider/          # Virtual Kubelet provider
│   ├── config/            # Configuration management
│   └── instances/         # Instance selection logic
├── internal/              # Private packages
│   ├── aws/               # AWS SDK integration
│   ├── container/         # Container runtime
│   └── metrics/           # Metrics and monitoring
├── docs/                  # Documentation
├── deploy/                # Deployment manifests
├── examples/              # User examples
└── scripts/               # Utility scripts
```

## Development Workflow

### Building

```bash
# Build the binary
make build

# Output: bin/orca
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run tests for specific package
go test -v ./pkg/provider/...
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Run all quality checks
make fmt && make lint && make vet && make test
```

### Running Locally

```bash
# Create a config file
cat > config.yaml <<EOF
aws:
  region: us-west-2
  credentials:
    accessKeyID: AKIA...
    secretAccessKey: your-secret-key

node:
  name: orca-dev-node
  operatingSystem: Linux
  cpu: 1000
  memory: 1Ti
  pods: 1000

logging:
  level: debug
EOF

# Run ORCA
make run

# Or run directly
./bin/orca --config config.yaml --kubeconfig ~/.kube/config
```

## Testing with Local Kubernetes

### Setup kind Cluster

```bash
# Create a kind cluster
kind create cluster --name orca-dev

# Verify
kubectl cluster-info --context kind-orca-dev
```

### Deploy ORCA to kind

```bash
# Build Docker image
make docker-build

# Load image into kind
kind load docker-image orca:latest --name orca-dev

# Deploy
kubectl apply -f deploy/kubernetes/

# Verify
kubectl get pods -n kube-system | grep orca
```

### Test Pod Creation

```bash
# Create a test pod
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: test-burst-pod
  annotations:
    orca.research/instance-type: "t3.small"
spec:
  nodeSelector:
    orca.research/provider: "aws"
  tolerations:
  - key: orca.research/burst-node
    operator: Equal
    value: "true"
    effect: NoSchedule
  containers:
  - name: test
    image: busybox
    command: ["sleep", "3600"]
EOF

# Check pod status
kubectl get pod test-burst-pod -o wide

# Check ORCA logs
kubectl logs -n kube-system -l app=orca
```

## AWS Configuration

### Using AWS CLI Credentials

```bash
# Configure AWS CLI
aws configure

# ORCA will use ~/.aws/credentials automatically
```

### Using IAM Role (Recommended for Production)

```yaml
# config.yaml
aws:
  region: us-west-2
  # No credentials needed - uses IAM role
```

### Required IAM Permissions

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:RunInstances",
        "ec2:TerminateInstances",
        "ec2:DescribeInstances",
        "ec2:DescribeInstanceTypes",
        "ec2:CreateTags",
        "ec2:DescribeTags"
      ],
      "Resource": "*"
    }
  ]
}
```

## Debugging

### Enable Debug Logging

```bash
# Via config file
cat > config.yaml <<EOF
logging:
  level: debug
EOF

# Via command line
./bin/orca --log-level=debug
```

### Common Issues

#### Issue: Cannot connect to Kubernetes cluster

```bash
# Verify kubeconfig
kubectl cluster-info

# Specify kubeconfig explicitly
./bin/orca --kubeconfig ~/.kube/config
```

#### Issue: AWS credentials not found

```bash
# Verify AWS credentials
aws sts get-caller-identity

# Set credentials explicitly in config.yaml
```

#### Issue: Pod stuck in Pending state

```bash
# Check pod events
kubectl describe pod <pod-name>

# Check ORCA logs
kubectl logs -n kube-system -l app=orca

# Check AWS instance creation
aws ec2 describe-instances --filters "Name=tag:orca.research/pod,Values=<pod-name>"
```

## Architecture Patterns

### Provider Interface

```go
// pkg/provider/provider.go

type Provider interface {
    // Pod lifecycle
    CreatePod(ctx context.Context, pod *corev1.Pod) error
    UpdatePod(ctx context.Context, pod *corev1.Pod) error
    DeletePod(ctx context.Context, pod *corev1.Pod) error
    GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error)
    GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error)
    GetPods(ctx context.Context) ([]*corev1.Pod, error)
}
```

### Instance Selection

```go
// pkg/instances/selector.go

// Selector chooses the appropriate EC2 instance type
type Selector interface {
    // Select returns instance type for pod
    Select(pod *corev1.Pod) (string, error)
}

// Three selection strategies
type ExplicitSelector struct{}  // Priority 1: User-specified
type TemplateSelector struct{}  // Priority 2: Named templates
type AutoSelector struct{}      // Priority 3: Auto-selection
```

### Configuration Management

```go
// pkg/config/config.go

type Config struct {
    AWS       AWSConfig       `yaml:"aws"`
    Node      NodeConfig      `yaml:"node"`
    Instances InstancesConfig `yaml:"instances"`
    Logging   LoggingConfig   `yaml:"logging"`
}
```

## Writing Tests

### Unit Test Example

```go
// pkg/instances/selector_test.go

func TestExplicitSelector(t *testing.T) {
    tests := []struct {
        name        string
        annotations map[string]string
        expected    string
        expectError bool
    }{
        {
            name: "explicit p5.48xlarge",
            annotations: map[string]string{
                "orca.research/instance-type": "p5.48xlarge",
            },
            expected:    "p5.48xlarge",
            expectError: false,
        },
        {
            name:        "no annotation",
            annotations: map[string]string{},
            expected:    "",
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            pod := &corev1.Pod{
                ObjectMeta: metav1.ObjectMeta{
                    Annotations: tt.annotations,
                },
            }

            selector := NewExplicitSelector()
            result, err := selector.Select(pod)

            if tt.expectError && err == nil {
                t.Error("expected error, got nil")
            }
            if !tt.expectError && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

### Integration Test Example

```go
// internal/aws/client_test.go

// +build integration

func TestCreateInstance(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Setup
    client := NewClient(/* config */)
    pod := createTestPod()

    // Execute
    err := client.CreateInstance(context.Background(), pod)
    if err != nil {
        t.Fatalf("failed to create instance: %v", err)
    }

    // Verify
    instances, err := client.DescribeInstances(context.Background(), pod)
    if err != nil || len(instances) != 1 {
        t.Errorf("expected 1 instance, got %d", len(instances))
    }

    // Cleanup
    defer client.TerminateInstance(context.Background(), instances[0].ID)
}
```

## Performance Profiling

### CPU Profiling

```bash
# Build with profiling
go build -o bin/orca ./cmd/orca

# Run with CPU profiling
./bin/orca --cpuprofile=cpu.prof

# Analyze profile
go tool pprof cpu.prof
```

### Memory Profiling

```bash
# Run with memory profiling
./bin/orca --memprofile=mem.prof

# Analyze profile
go tool pprof mem.prof
```

## Release Process

### Version Bump

```bash
# Update VERSION file
echo "0.2.0" > VERSION

# Update CHANGELOG.md
# Move [Unreleased] items to [0.2.0] section

# Commit
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to 0.2.0"

# Tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

## Getting Help

- **GitHub Issues** - Report bugs or request features
- **GitHub Discussions** - Ask questions
- **Research Partners** - Contact NRP, SDSU teams

## Additional Resources

- [Virtual Kubelet Documentation](https://virtual-kubelet.io/)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/docs/)
- [Kubernetes Client Go](https://github.com/kubernetes/client-go)
- [Effective Go](https://go.dev/doc/effective_go)
