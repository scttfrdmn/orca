# LocalStack Testing Guide

ORCA supports LocalStack for local AWS testing without charges. This enables rapid development and testing without touching real AWS resources.

## What is LocalStack?

LocalStack is a fully functional local AWS cloud stack that provides drop-in replacements for AWS services including EC2, S3, Lambda, and more.

## Quick Start

### 1. Start LocalStack

```bash
./scripts/setup-localstack.sh
```

This script will:
- Start LocalStack Docker container
- Create VPC, subnet, and security group
- Update `config.localstack.yaml` with resource IDs
- Display configuration details

### 2. Run ORCA with LocalStack

```bash
# Build ORCA
make build

# Run with LocalStack config
./bin/orca --config config.localstack.yaml
```

### 3. Test Pod Creation

```bash
# In another terminal, create a test pod
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: localstack-test-pod
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
    command: ["sleep", "300"]
EOF

# Watch pod status
kubectl get pod localstack-test-pod -w
```

## LocalStack vs Real AWS

| Feature | LocalStack | Real AWS |
|---------|-----------|----------|
| **Cost** | Free | Pay per use |
| **Speed** | Fast (local) | Network latency |
| **Data** | Ephemeral | Persistent |
| **Instances** | Simulated | Real VMs |
| **Networking** | Mocked | Full AWS networking |
| **AMIs** | Any ID works | Must be valid |

## Configuration

### LocalStack Endpoint

```yaml
# config.localstack.yaml
aws:
  region: us-east-1
  localStackEndpoint: "http://localhost:4566"
  credentials:
    accessKeyID: "test"
    secretAccessKey: "test"
```

### Environment Variables

LocalStack requires these environment variables:

```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_ENDPOINT_URL=http://localhost:4566
```

## Testing Workflow

### Development Cycle

1. **Write code** - Implement feature
2. **Unit tests** - Test in isolation with mocks
3. **LocalStack tests** - Test AWS integration locally
4. **Real AWS tests** - Final validation (minimal)

### Example: Testing Instance Creation

```bash
# 1. Start LocalStack
./scripts/setup-localstack.sh

# 2. Run ORCA
./bin/orca --config config.localstack.yaml &
ORCA_PID=$!

# 3. Create test pod
kubectl apply -f examples/test-pod.yaml

# 4. Verify instance created in LocalStack
aws ec2 describe-instances \
  --endpoint-url http://localhost:4566 \
  --filters "Name=tag:orca.research/provider,Values=orca"

# 5. Cleanup
kubectl delete pod localstack-test-pod
kill $ORCA_PID
```

## Debugging LocalStack

### Check LocalStack Health

```bash
curl http://localhost:4566/_localstack/health
```

### View LocalStack Logs

```bash
docker logs orca-localstack
```

### List Resources in LocalStack

```bash
# List instances
aws ec2 describe-instances \
  --endpoint-url http://localhost:4566

# List VPCs
aws ec2 describe-vpcs \
  --endpoint-url http://localhost:4566

# List security groups
aws ec2 describe-security-groups \
  --endpoint-url http://localhost:4566
```

## Common Issues

### LocalStack not starting

```bash
# Check Docker
docker ps | grep localstack

# Restart LocalStack
docker restart orca-localstack

# Check logs
docker logs orca-localstack
```

### Connection refused

```bash
# Verify endpoint
curl http://localhost:4566/_localstack/health

# Check port forwarding
docker port orca-localstack
```

### Invalid VPC/Subnet IDs

```bash
# Recreate resources
./scripts/setup-localstack.sh
```

## Cleanup

### Stop LocalStack

```bash
docker stop orca-localstack
docker rm orca-localstack
```

### Full Cleanup

```bash
# Stop and remove LocalStack
docker stop orca-localstack
docker rm orca-localstack

# Remove images (optional)
docker rmi localstack/localstack:latest
```

## Integration Tests with LocalStack

Integration tests can use LocalStack automatically:

```go
// +build integration

func TestCreateInstance_LocalStack(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Load LocalStack config
    cfg, err := config.LoadConfig("../../config.localstack.yaml")
    if err != nil {
        t.Fatal(err)
    }

    // Create client (will use LocalStack endpoint)
    client, err := aws.NewClient(context.Background(), cfg)
    if err != nil {
        t.Fatal(err)
    }

    // Create test pod
    pod := createTestPod()

    // Create instance
    instanceID, err := client.CreateInstance(context.Background(), pod, "t3.small")
    if err != nil {
        t.Fatal(err)
    }

    // Verify instance exists
    instance, err := client.DescribeInstance(context.Background(), instanceID)
    if err != nil {
        t.Fatal(err)
    }

    if instance.State != "running" && instance.State != "pending" {
        t.Errorf("expected running/pending, got %s", instance.State)
    }

    // Cleanup
    defer client.TerminateInstance(context.Background(), instanceID)
}
```

## CI/CD with LocalStack

GitHub Actions can use LocalStack in CI:

```yaml
# .github/workflows/integration-test.yml
jobs:
  integration-test:
    runs-on: ubuntu-latest
    services:
      localstack:
        image: localstack/localstack:latest
        ports:
          - 4566:4566
        env:
          SERVICES: ec2
          DEBUG: 1

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Setup LocalStack
        run: ./scripts/setup-localstack.sh

      - name: Run integration tests
        run: make integration-test
        env:
          AWS_ENDPOINT_URL: http://localhost:4566
```

## Best Practices

1. **Always use LocalStack first** - Test locally before real AWS
2. **Reset between tests** - Restart LocalStack for clean state
3. **Check health** - Verify LocalStack is ready before tests
4. **Use realistic data** - Test with production-like configurations
5. **Document assumptions** - Note LocalStack behavior differences

## Limitations

LocalStack is not a perfect AWS replica:

- ❌ No real networking between instances
- ❌ No real VM creation
- ❌ No actual GPU hardware
- ❌ Some API features may differ
- ❌ Performance characteristics differ

For production validation, always test on real AWS after LocalStack testing.

## Further Reading

- [LocalStack Documentation](https://docs.localstack.cloud/)
- [LocalStack AWS Coverage](https://docs.localstack.cloud/references/coverage/)
- [LocalStack GitHub](https://github.com/localstack/localstack)
