# LocalStack Testing Guide

This guide explains how to use LocalStack for local development and testing of ORCA without incurring AWS costs.

## Overview

LocalStack is a fully functional local AWS cloud stack that emulates AWS services on your local machine. This allows you to:

- Develop and test AWS integrations without AWS costs
- Test faster (no internet latency)
- Test in isolation (no conflicts with production resources)
- Run integration tests in CI/CD pipelines

## Prerequisites

- Docker and Docker Compose installed
- AWS CLI installed (for awslocal commands)
- Python 3 and pip3 (for awslocal wrapper)

### Installing awslocal

The `awslocal` CLI is a wrapper around the AWS CLI that automatically configures it to use LocalStack:

```bash
pip3 install awscli-local
```

Alternatively, you can use the regular AWS CLI with LocalStack by setting environment variables:

```bash
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-west-2
export AWS_ENDPOINT_URL=http://localhost:4566
```

## Quick Start

### 1. Start LocalStack

```bash
make localstack-start
```

This will:
- Start LocalStack in a Docker container
- Expose the LocalStack gateway on port 4566
- Automatically run initialization scripts to create test resources
- Store data in `/tmp/localstack` for persistence

### 2. Check Initialization Status

```bash
# View LocalStack logs
make localstack-logs

# Check resource IDs
make localstack-status
```

The initialization script creates:
- VPC with DNS support
- Internet Gateway
- Public subnet in us-west-2a
- Route table with internet access
- Security group with SSH access
- Test AMI for launching instances

### 3. Verify LocalStack is Ready

```bash
./scripts/wait-for-localstack.sh
```

### 4. Run Tests

```bash
# Run integration tests
make test-integration

# Or run all tests
make test
```

### 5. Run ORCA Locally

```bash
make run-local
```

This starts ORCA with the LocalStack configuration (`config.localstack.yaml`).

### 6. Stop LocalStack

```bash
make localstack-stop
```

## Configuration

### LocalStack Configuration

The LocalStack environment is configured in `docker-compose.localstack.yml`:

- **Services**: EC2, IAM, STS, CloudWatch, CloudWatch Logs
- **Endpoint**: http://localhost:4566
- **Persistence**: Enabled (data survives container restarts)
- **VM Manager**: Docker (for EC2 instances)

### ORCA Configuration

ORCA uses `config.localstack.yaml` for LocalStack testing:

```yaml
aws:
  region: us-west-2
  localStackEndpoint: http://localhost:4566
  credentials:
    accessKeyID: test
    secretAccessKey: test
  # ... other settings
```

## Working with LocalStack

### Querying Resources

```bash
# List EC2 instances
awslocal ec2 describe-instances --region us-west-2

# List VPCs
awslocal ec2 describe-vpcs --region us-west-2

# List security groups
awslocal ec2 describe-security-groups --region us-west-2

# List subnets
awslocal ec2 describe-subnets --region us-west-2
```

### Launching Test Instances

```bash
# Get resource IDs from initialization
source /tmp/localstack-orca-resources.env

# Launch an instance
awslocal ec2 run-instances \
  --image-id $LOCALSTACK_AMI_ID \
  --instance-type t3.medium \
  --subnet-id $LOCALSTACK_SUBNET_ID \
  --security-group-ids $LOCALSTACK_SG_ID \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=test-instance}]'
```

### Viewing Logs

```bash
# Follow all LocalStack logs
make localstack-logs

# View specific service logs
docker exec orca-localstack cat /var/log/localstack/ec2.log
```

### Opening a Shell in LocalStack

```bash
make localstack-shell
```

## Testing Workflow

### Unit Tests

Unit tests don't require LocalStack and use mocks:

```bash
make test
```

### Integration Tests

Integration tests connect to LocalStack and test real AWS SDK interactions:

```bash
# Start LocalStack
make localstack-start

# Wait for it to be ready
./scripts/wait-for-localstack.sh

# Run integration tests
make test-integration

# Stop LocalStack when done
make localstack-stop
```

### Testing ORCA End-to-End

```bash
# 1. Start LocalStack
make localstack-start

# 2. Build ORCA
make build

# 3. Run ORCA with LocalStack config
make run-local

# 4. In another terminal, deploy a test pod
kubectl apply -f examples/test-pod.yaml

# 5. Check ORCA logs to see instance creation
# 6. Verify instance in LocalStack
awslocal ec2 describe-instances

# 7. Stop ORCA (Ctrl+C)
# 8. Stop LocalStack
make localstack-stop
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make localstack-start` | Start LocalStack container |
| `make localstack-stop` | Stop LocalStack container |
| `make localstack-restart` | Restart LocalStack |
| `make localstack-logs` | Follow LocalStack logs |
| `make localstack-status` | Show created resource IDs |
| `make localstack-shell` | Open bash in LocalStack container |
| `make test-integration` | Run integration tests |
| `make run-local` | Run ORCA with LocalStack config |

## Troubleshooting

### LocalStack Won't Start

```bash
# Check if port 4566 is already in use
lsof -i :4566

# Check Docker status
docker ps -a | grep localstack

# View container logs
docker logs orca-localstack

# Remove and restart
make localstack-stop
docker rm -f orca-localstack
make localstack-start
```

### Resources Not Initialized

```bash
# Check initialization logs
make localstack-logs | grep -A 20 "Initializing LocalStack"

# Resource IDs are saved here
cat /tmp/localstack-orca-resources.env

# Re-run initialization manually
docker exec orca-localstack /docker-entrypoint-initaws.d/01-init-ec2.sh
```

### Tests Fail with Connection Errors

```bash
# Verify LocalStack is running
docker ps | grep localstack

# Check LocalStack health
curl http://localhost:4566/_localstack/health

# Verify EC2 service is available
awslocal ec2 describe-regions --region us-west-2

# Wait for LocalStack to be fully ready
./scripts/wait-for-localstack.sh
```

### EC2 Instances Don't Start

LocalStack uses Docker-in-Docker for EC2 instances. Ensure:

```bash
# Docker socket is mounted (in docker-compose.localstack.yml)
volumes:
  - "/var/run/docker.sock:/var/run/docker.sock"

# Check LocalStack EC2 logs
docker exec orca-localstack cat /var/log/localstack/ec2.log
```

**Note**: LocalStack's EC2 emulation has limitations compared to real AWS:
- Instance types are simulated (no real resource limits)
- Networking is simplified
- Some advanced EC2 features may not work

### Clear All LocalStack Data

```bash
# Stop LocalStack
make localstack-stop

# Clear persistent data
rm -rf /tmp/localstack/*

# Clear resource IDs
rm -f /tmp/localstack-orca-resources.env

# Restart fresh
make localstack-start
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Start LocalStack
        run: make localstack-start

      - name: Wait for LocalStack
        run: ./scripts/wait-for-localstack.sh

      - name: Run Integration Tests
        run: make test-integration

      - name: Stop LocalStack
        if: always()
        run: make localstack-stop
```

## Differences from Real AWS

Be aware of these differences when using LocalStack:

### Supported Features
- ✅ EC2 instance launch/terminate
- ✅ VPC, subnets, security groups
- ✅ Tags and resource filtering
- ✅ IAM roles and policies (basic)
- ✅ CloudWatch Logs
- ✅ Basic networking

### Limited/Unsupported Features
- ❌ Spot instance pricing (spot launches work, but no real pricing)
- ❌ Capacity Reservations
- ❌ Real instance resource limits
- ❌ Some advanced networking features
- ❌ AWS-specific instance metadata
- ❌ Real GPU support

### Testing Strategy

Use LocalStack for:
- **Unit-level integration tests**: Testing AWS SDK calls
- **Logic validation**: Testing tag application, instance selection
- **Error handling**: Testing timeout, failure scenarios
- **Development**: Fast iteration without AWS costs

Use real AWS for:
- **End-to-end testing**: Full production-like validation
- **Performance testing**: Real instance performance
- **Advanced features**: Capacity Reservations, Spot, GPUs
- **Pre-production validation**: Final testing before release

## Best Practices

1. **Always wait for initialization**: Use `./scripts/wait-for-localstack.sh`
2. **Use resource IDs from env file**: Source `/tmp/localstack-orca-resources.env`
3. **Check logs frequently**: LocalStack logs reveal issues quickly
4. **Clear data between test runs**: Ensure clean state
5. **Don't rely on persistence for CI**: Treat as ephemeral
6. **Test failures locally first**: LocalStack makes debugging easier
7. **Validate against real AWS**: LocalStack is a simulation

## Additional Resources

- [LocalStack Documentation](https://docs.localstack.cloud/)
- [LocalStack AWS Service Coverage](https://docs.localstack.cloud/references/coverage/)
- [awslocal CLI Reference](https://github.com/localstack/awscli-local)
- [LocalStack GitHub](https://github.com/localstack/localstack)

## Getting Help

If you encounter issues with LocalStack:

1. Check LocalStack logs: `make localstack-logs`
2. Verify health: `curl http://localhost:4566/_localstack/health`
3. Check GitHub issues: https://github.com/localstack/localstack/issues
4. LocalStack Community: https://discuss.localstack.cloud/
