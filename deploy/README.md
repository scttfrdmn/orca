# ORCA Deployment Guide

This directory contains Kubernetes manifests for deploying ORCA to your cluster.

## Prerequisites

- Kubernetes cluster (1.28+)
- `kubectl` configured to access your cluster
- AWS account with appropriate permissions
- AWS VPC, subnet, and security group created

## Quick Start

### 1. Configure AWS Resources

Edit `configmap.yaml` and update the AWS configuration:

```yaml
aws:
  region: us-west-2  # Your AWS region
  vpcID: "vpc-XXXXXXXX"  # Your VPC ID
  subnetID: "subnet-XXXXXXXX"  # Your subnet ID
  securityGroupIDs:
    - "sg-XXXXXXXX"  # Your security group ID
```

### 2. Configure AWS Credentials

**Option A: IAM Role for Service Accounts (IRSA) - Recommended**

1. Create an IAM role with EC2 permissions
2. Associate the role with the service account:

```bash
kubectl annotate serviceaccount orca \
  -n orca-system \
  eks.amazonaws.com/role-arn=arn:aws:iam::ACCOUNT_ID:role/orca-role
```

3. Leave `secret.yaml` empty or omit it

**Option B: Explicit Credentials**

Edit `secret.yaml` and add your AWS credentials:

```yaml
stringData:
  accessKeyID: "AKIAIOSFODNN7EXAMPLE"
  secretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
```

**Option C: EC2 Instance Profile**

If your cluster nodes run on EC2 with an instance profile that has EC2 permissions, credentials will be auto-discovered. Leave `secret.yaml` empty.

### 3. Deploy ORCA

```bash
# Apply all manifests
kubectl apply -f deploy/

# Or use Kustomize
kubectl apply -k deploy/
```

### 4. Verify Deployment

```bash
# Check ORCA pod is running
kubectl get pods -n orca-system

# Check virtual node registered
kubectl get nodes
# Should see: orca-aws-node

# Check node details
kubectl describe node orca-aws-node

# Check logs
kubectl logs -n orca-system -l app.kubernetes.io/name=orca -f
```

## Deploy a Test Pod

Create a test pod that will burst to AWS:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: burst-test
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
```

Apply and watch:

```bash
kubectl apply -f test-pod.yaml
kubectl get pod burst-test -w
```

## Configuration Options

### Instance Selection Modes

Edit `configmap.yaml` to change instance selection behavior:

```yaml
instances:
  # explicit: Users must specify instance type via annotation
  # template: Users specify workload template name
  # auto: ORCA auto-selects based on resource requests
  selectionMode: explicit
```

### Budget Controls

Configure spending limits in `configmap.yaml`:

```yaml
limits:
  maxConcurrentInstances: 100
  maxInstancesPerNamespace: 50
  dailyBudget: 5000.00
  monthlyBudget: 100000.00
  maxInstanceLifetime: 24h
```

### Workload Templates

Define named workload templates for common use cases:

```yaml
instances:
  templates:
    llm-training:
      instanceType: p5.48xlarge
      launchType: spot

    inference:
      instanceType: g5.xlarge
      launchType: on-demand
```

Users can then reference templates:

```yaml
annotations:
  orca.research/workload-template: "llm-training"
```

## Monitoring

ORCA exposes Prometheus metrics on port 8080:

```bash
# Port-forward to access metrics
kubectl port-forward -n orca-system svc/orca-metrics 8080:8080

# View metrics
curl http://localhost:8080/metrics
```

### ServiceMonitor for Prometheus Operator

If using Prometheus Operator, create a ServiceMonitor:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: orca
  namespace: orca-system
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: orca
  endpoints:
    - port: metrics
      interval: 30s
```

## Upgrading

### Update Configuration

```bash
# Edit ConfigMap
kubectl edit configmap orca-config -n orca-system

# Restart ORCA to pick up changes
kubectl rollout restart deployment/orca -n orca-system
```

### Update Image

```bash
# Update to specific version
kubectl set image deployment/orca -n orca-system \
  orca=ghcr.io/scttfrdmn/orca:v0.2.0

# Or edit deployment
kubectl edit deployment orca -n orca-system
```

## Troubleshooting

### Pod Not Scheduling

```bash
# Check node status
kubectl describe node orca-aws-node

# Check pod tolerations and node selector
kubectl describe pod <pod-name>

# Check ORCA logs
kubectl logs -n orca-system -l app.kubernetes.io/name=orca
```

### AWS Errors

```bash
# Check ORCA has AWS credentials
kubectl exec -n orca-system -it deploy/orca -- env | grep AWS

# Check AWS permissions
# ORCA needs: ec2:RunInstances, ec2:TerminateInstances, ec2:DescribeInstances

# Test AWS connectivity
kubectl logs -n orca-system -l app.kubernetes.io/name=orca | grep -i aws
```

### Instance Not Creating

```bash
# Check pod annotations
kubectl get pod <pod-name> -o yaml | grep annotations -A 5

# Verify instance type is allowed
kubectl get configmap orca-config -n orca-system -o yaml | grep allowedInstanceTypes -A 10

# Check budget limits
kubectl logs -n orca-system -l app.kubernetes.io/name=orca | grep -i budget
```

## Uninstall

```bash
# Delete all ORCA resources
kubectl delete -f deploy/

# Or with Kustomize
kubectl delete -k deploy/

# Virtual node will be automatically removed
```

## AWS IAM Policy

ORCA requires these AWS permissions:

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
        "ec2:DescribeInstanceStatus",
        "ec2:DescribeInstanceTypes",
        "ec2:DescribeImages",
        "ec2:CreateTags",
        "ec2:DescribeTags"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeSubnets",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeVpcs"
      ],
      "Resource": "*"
    }
  ]
}
```

## Security Best Practices

1. **Use IRSA**: Prefer IAM Role for Service Accounts over static credentials
2. **Least Privilege**: Grant only required EC2 permissions
3. **Network Isolation**: Use dedicated VPC/subnet for burst instances
4. **Budget Limits**: Configure appropriate budget controls
5. **Audit Logging**: Enable CloudTrail for all EC2 API calls
6. **Pod Security**: Ensure pods run with appropriate security contexts
7. **Secrets Management**: Use AWS Secrets Manager for sensitive data

## Advanced Configuration

### Custom Node Labels

Add custom labels to the virtual node:

```yaml
node:
  labels:
    orca.research/provider: "aws"
    environment: "production"
    department: "research"
```

### Custom Node Taints

Control pod scheduling with taints:

```yaml
node:
  taints:
    - key: orca.research/burst-node
      value: "true"
      effect: NoSchedule
    - key: gpu
      value: "required"
      effect: NoExecute
```

### Multi-Region Deployment

Deploy multiple ORCA instances for different regions:

```bash
# Copy and modify for each region
cp -r deploy/ deploy-us-east-1/
# Edit deploy-us-east-1/configmap.yaml - change region to us-east-1
# Edit deploy-us-east-1/namespace.yaml - change to orca-system-us-east-1

kubectl apply -k deploy-us-east-1/
```

## Support

For issues, questions, or contributions:
- GitHub Issues: https://github.com/scttfrdmn/orca/issues
- Documentation: https://github.com/scttfrdmn/orca/blob/main/README.md
