# EKS Cluster: NRP Development Replica

This Terraform configuration creates an EKS cluster that replicates key aspects of the NRP (National Research Platform) environment for ORCA development and testing.

## What This Replicates

### NRP Nautilus Characteristics

The NRP Nautilus cluster has several unique characteristics we want to replicate:

1. **Multi-tenant by design** - Multiple departments/research groups sharing infrastructure
2. **GPU-heavy workloads** - Primarily AI/ML research with various GPU types
3. **Distributed storage** - Ceph S3-compatible storage
4. **Heterogeneous nodes** - Mix of CPU and GPU nodes across sites
5. **Resource quotas** - Per-namespace limits and budgets
6. **Federation** - Multiple clusters across universities

### What We'll Replicate in EKS

```
┌─────────────────────────────────────────────────────┐
│  EKS Control Plane                                  │
└─────────────────────────────────────────────────────┘
              │
              │
    ┌─────────┼─────────┐
    │         │         │
┌───▼────┬───▼────┬───▼────┐
│ Node   │ Node   │ Node   │  Regular worker nodes
│ Group  │ Group  │ Group  │  (t3.large - like NRP CPU nodes)
│   1    │   2    │   3    │
└────────┴────────┴────────┘

┌────────────────────────────────────────────┐
│  ORCA Virtual Nodes (not real EKS nodes)   │
│  ├── orca-gpu-node    → AWS GPU instances  │
│  ├── orca-cpu-node    → AWS CPU instances  │
│  └── orca-spot-node   → AWS Spot instances │
└────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  Namespaces (Multi-tenant)                  │
│  ├── biology-dept    (Budget: $500/mo)      │
│  ├── cs-dept         (Budget: $1000/mo)     │
│  ├── physics-dept    (Budget: $750/mo)      │
│  └── research-lab-x  (Budget: $300/mo)      │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│  Storage                                    │
│  ├── S3 (replicates Ceph)                   │
│  ├── EFS (shared filesystem)                │
│  └── EBS (persistent volumes)               │
└─────────────────────────────────────────────┘
```

## Features

- ✅ **Multi-tenant namespaces** - Separate departments with quotas
- ✅ **Mixed workload types** - CPU, GPU, batch jobs
- ✅ **S3 storage** - Simulates NRP Ceph S3 gateway
- ✅ **Cost tracking** - Per-namespace budget tracking
- ✅ **GPU node groups** - For testing without ORCA (optional)
- ✅ **Spot instances** - Cost optimization like NRP
- ✅ **Resource quotas** - Per-namespace limits
- ✅ **RBAC** - Multi-tenant access control

## Cost Estimate

### Minimal Configuration (Development)
```
EKS Control Plane:         $73/month
Worker Nodes (3x t3.large): $150/month (spot: $60/month)
NAT Gateway:               $33/month
Data Transfer:             ~$10/month
─────────────────────────────────────────
Total (on-demand):         ~$266/month
Total (with spot):         ~$176/month
```

### With GPU Nodes (Optional - for testing)
```
Above +
1x g5.xlarge GPU node:     $730/month (spot: $290/month)
─────────────────────────────────────────
Total with GPU:            ~$996/month (on-demand)
Total with GPU (spot):     ~$466/month
```

### Recommended: CPU-only for ORCA testing
```
Use ORCA to burst to GPU instances when needed
Base cluster: ~$176/month (spot)
GPU testing: Pay-per-use through ORCA (~$10-50/month)
─────────────────────────────────────────
Total monthly cost:        ~$200-250/month
```

## Quick Start

### 1. Prerequisites

```bash
# Install required tools
brew install terraform kubectl aws-cli

# Configure AWS credentials
aws configure

# Set your region
export AWS_REGION=us-west-2
```

### 2. Deploy the Cluster

```bash
cd deploy/terraform/eks-nrp-replica

# Review configuration
cat terraform.tfvars.example

# Copy and customize
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars

# Deploy
terraform init
terraform plan
terraform apply
```

### 3. Configure kubectl

```bash
# Get kubeconfig
aws eks update-kubeconfig \
  --region us-west-2 \
  --name orca-nrp-dev

# Verify
kubectl get nodes
kubectl get namespaces
```

### 4. Deploy ORCA

```bash
# From project root
cd ../../../

# Build and push ORCA image to ECR
make docker-build
aws ecr get-login-password | docker login ...
docker push ...

# Deploy ORCA to EKS
kubectl apply -f deploy/kubernetes/

# Verify ORCA is running
kubectl get pods -n kube-system | grep orca
kubectl get nodes  # Should see ORCA virtual nodes
```

### 5. Test Multi-Tenant Setup

```bash
# Deploy test workloads for each department
kubectl apply -f examples/nrp-replica/biology-dept/
kubectl apply -f examples/nrp-replica/cs-dept/
kubectl apply -f examples/nrp-replica/physics-dept/

# Watch pods burst to AWS
kubectl get pods -A -o wide --watch
```

## Configuration Options

### terraform.tfvars

```hcl
# Basic Configuration
cluster_name = "orca-nrp-dev"
region       = "us-west-2"

# Worker Nodes
node_instance_type = "t3.large"     # Similar to NRP CPU nodes
node_count_min     = 2
node_count_max     = 5
node_count_desired = 3
node_use_spot      = true           # 60-70% cost savings

# Optional: GPU Nodes (for testing without ORCA)
enable_gpu_nodes   = false          # Set to true if you want real GPU nodes
gpu_instance_type  = "g5.xlarge"
gpu_count          = 1

# Namespaces (departments)
departments = [
  {
    name          = "biology-dept"
    daily_budget  = 20.00
    monthly_budget = 500.00
    max_gpus      = 4
  },
  {
    name          = "cs-dept"
    daily_budget  = 40.00
    monthly_budget = 1000.00
    max_gpus      = 8
  },
  {
    name          = "physics-dept"
    daily_budget  = 30.00
    monthly_budget = 750.00
    max_gpus      = 6
  }
]

# Storage
enable_efs = true  # Shared filesystem
enable_s3_csi = true  # S3 access (like Ceph)

# Cost Controls
enable_cluster_autoscaler = true
enable_node_termination_handler = true  # Graceful spot termination
```

## NRP-Specific Features

### 1. Ceph Storage Simulation

NRP uses Ceph with S3-compatible gateway. We'll use AWS S3:

```yaml
# Pod with S3 mount (like NRP Ceph)
apiVersion: v1
kind: Pod
metadata:
  name: research-job
  namespace: biology-dept
spec:
  containers:
  - name: training
    image: pytorch:latest
    env:
    - name: S3_BUCKET
      value: orca-nrp-biology-data
    - name: AWS_REGION
      value: us-west-2
    volumeMounts:
    - name: data
      mountPath: /data
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: s3-data-pvc
```

### 2. Multi-Tenant Isolation

Each department gets:
- Dedicated namespace
- Resource quotas
- Budget limits
- RBAC policies
- Cost tracking tags

```yaml
# Auto-created by Terraform
apiVersion: v1
kind: Namespace
metadata:
  name: biology-dept
  labels:
    department: biology
    orca.research/budget-daily: "20.00"
    orca.research/budget-monthly: "500.00"
---
apiVersion: v1
kind: ResourceQuota
metadata:
  name: biology-dept-quota
  namespace: biology-dept
spec:
  hard:
    requests.cpu: "100"
    requests.memory: "500Gi"
    requests.nvidia.com/gpu: "4"
    persistentvolumeclaims: "10"
```

### 3. Workload Templates

NRP users often run similar workloads. We'll provide templates:

```bash
# LLM Training (like biology protein folding)
kubectl apply -f examples/nrp-replica/templates/llm-training.yaml

# Computer Vision (like medical imaging)
kubectl apply -f examples/nrp-replica/templates/vision-training.yaml

# Batch Processing (like genomics pipelines)
kubectl apply -f examples/nrp-replica/templates/batch-job.yaml
```

## Testing Scenarios

### Scenario 1: Department Bursting

Simulate biology department bursting to AWS:

```bash
# Biology dept hits local resource limits
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: protein-folding-job
  namespace: biology-dept
  annotations:
    orca.research/instance-type: "p5.48xlarge"
    orca.research/launch-type: "spot"
spec:
  nodeSelector:
    orca.research/provider: "aws"
  containers:
  - name: alphafold
    image: alphafold:latest
    resources:
      limits:
        nvidia.com/gpu: 8
EOF
```

### Scenario 2: Multi-Department Concurrent Usage

```bash
# Multiple departments burst simultaneously
./scripts/test-multi-tenant.sh

# Expected:
# - Biology: 2 pods on p5.48xlarge
# - CS: 3 pods on g5.4xlarge
# - Physics: 1 pod on c7i.8xlarge
# - Budget tracking per department
# - No interference between departments
```

### Scenario 3: Budget Enforcement

```bash
# CS dept exceeds daily budget
kubectl apply -f examples/nrp-replica/cs-dept/exceed-budget.yaml

# Expected:
# - First 5 jobs succeed
# - 6th job rejected with budget error
# - Budget resets next day
```

### Scenario 4: Spot Interruption

```bash
# Launch spot instances
kubectl apply -f examples/nrp-replica/spot-workload.yaml

# Simulate interruption
./scripts/simulate-spot-interruption.sh

# Expected:
# - ORCA detects interruption
# - Pod marked as failed
# - Optional: Retry on new spot instance
```

## Comparison: NRP vs EKS Replica

| Feature | NRP Nautilus | EKS Replica | Notes |
|---------|--------------|-------------|-------|
| Multi-tenant | ✅ | ✅ | Namespaces + quotas |
| GPU workloads | ✅ | ✅ | Via ORCA bursting |
| Ceph storage | ✅ | ~S3~ | S3 similar enough |
| Federation | ✅ | ❌ | Single cluster (simpler) |
| Spot instances | ~Limited~ | ✅ | Better in AWS |
| Cost tracking | ~Basic~ | ✅ | Per-namespace tags |
| Geographic distribution | ✅ | ❌ | Multi-region possible later |

## Maintenance

### Daily Checks

```bash
# Check cluster health
kubectl get nodes
kubectl top nodes

# Check ORCA status
kubectl logs -n kube-system -l app=orca --tail=50

# Check costs (AWS Cost Explorer)
aws ce get-cost-and-usage \
  --time-period Start=$(date -d yesterday +%Y-%m-%d),End=$(date +%Y-%m-%d) \
  --granularity DAILY \
  --metrics BlendedCost
```

### Weekly Cleanup

```bash
# Remove completed pods
kubectl delete pods --field-selector status.phase=Succeeded -A

# Check for orphaned instances
./scripts/check-orphaned-instances.sh

# Review cost by namespace
./scripts/cost-report.sh
```

### Weekend Shutdown (Optional - Save ~$120/weekend)

```bash
# Friday 5pm
terraform apply -var="node_count_desired=0"

# Monday 9am
terraform apply -var="node_count_desired=3"
```

## Differences from Production NRP

This is a **development replica**, not production:

**Not included:**
- ❌ Geographic federation (multiple sites)
- ❌ Bare metal nodes
- ❌ InfiniBand networking
- ❌ LDAP/Active Directory integration
- ❌ Legacy workload support

**Included (for ORCA testing):**
- ✅ Multi-tenant namespaces
- ✅ GPU workload patterns
- ✅ S3 storage (similar to Ceph)
- ✅ Budget tracking
- ✅ Resource quotas
- ✅ RBAC
- ✅ Spot instances

## Cost Optimization

### Approach 1: Spot Only (Cheapest)
```hcl
node_use_spot = true
enable_gpu_nodes = false
```
**Cost**: ~$176/month base + ORCA usage

### Approach 2: On-Demand Workers
```hcl
node_use_spot = false
enable_gpu_nodes = false
```
**Cost**: ~$266/month base + ORCA usage

### Approach 3: GPU Node Included
```hcl
node_use_spot = true
enable_gpu_nodes = true
gpu_use_spot = true
```
**Cost**: ~$466/month (includes permanent GPU node)

**Recommendation**: Approach 1 (spot workers, no GPU nodes)
- Use ORCA for GPU bursting (pay-per-use)
- Most realistic for ORCA testing
- Lowest cost (~$200-250/month total)

## Terraform Modules

```
eks-nrp-replica/
├── main.tf              # Main cluster configuration
├── vpc.tf              # VPC, subnets, NAT
├── eks.tf              # EKS cluster
├── node-groups.tf      # Worker and GPU nodes
├── namespaces.tf       # Department namespaces
├── storage.tf          # S3, EFS setup
├── iam.tf              # Roles for ORCA, pods
├── variables.tf        # Input variables
├── outputs.tf          # Cluster info
├── terraform.tfvars.example
└── README.md           # This file
```

## Next Steps

After deploying:

1. **Deploy ORCA** - See main deployment docs
2. **Create test workloads** - Use examples in `examples/nrp-replica/`
3. **Monitor costs** - Set up billing alerts
4. **Test scenarios** - Run multi-tenant tests
5. **Iterate** - Adjust based on findings

## Getting Help

- EKS Issues: Check `kubectl logs` and CloudWatch
- ORCA Issues: See main project docs
- Terraform Issues: Run `terraform plan` to debug

## Teardown

```bash
# Remove all pods first (clean up burst instances)
kubectl delete pods --all -A

# Destroy cluster
terraform destroy

# Verify no orphaned resources
aws ec2 describe-instances --filters "Name=tag:orca.research/cluster,Values=orca-nrp-dev"
```

---

This setup gives you a realistic NRP-like environment for ~$200-250/month, perfect for ORCA development and validation before deploying to actual NRP infrastructure.
