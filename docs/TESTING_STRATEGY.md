# ORCA Testing Strategy

This document outlines the progressive testing strategy for ORCA, from local development to production-scale validation.

## Testing Phases

### Phase 1: Local Development (kind + Mock AWS)
**Goal**: Develop core functionality without AWS costs
**Duration**: Weeks 1-4
**Cost**: $0

#### Setup
```bash
# Install kind
brew install kind

# Create local cluster
./scripts/setup-kind-cluster.sh

# Run ORCA in mock mode
./scripts/test-local.sh
```

#### What to Test
- ✅ Virtual Kubelet provider interface
- ✅ Pod lifecycle management (create, update, delete)
- ✅ Instance selection logic (explicit, template, auto)
- ✅ Configuration parsing and validation
- ✅ Kubernetes API integration
- ✅ Error handling and edge cases
- ✅ Metrics and logging

#### Advantages
- Fast iteration cycle (seconds)
- No AWS charges
- Easy cleanup
- Multiple test scenarios quickly

---

### Phase 2: Local Development (kind + Real AWS)
**Goal**: Validate AWS integration with minimal costs
**Duration**: Weeks 4-8
**Cost**: ~$10-50/month

#### Setup
```bash
# Use kind cluster from Phase 1
# Update config.dev.yaml with real AWS credentials

# Run ORCA with real AWS
./bin/orca --config config.dev.yaml

# Deploy test pod (t3.small = ~$0.02/hour)
kubectl apply -f examples/test-pod.yaml
```

#### What to Test
- ✅ EC2 instance creation and termination
- ✅ Security groups and networking
- ✅ IAM permissions
- ✅ Instance tagging and pod mapping
- ✅ Spot instance handling
- ✅ Cost tracking and budget enforcement
- ✅ GPU instance types (brief tests)
- ✅ Instance lifecycle edge cases

#### Cost Controls
```yaml
# config.dev.yaml safeguards
limits:
  maxConcurrentInstances: 5
  dailyBudget: 50.00
  monthlyBudget: 500.00
  maxInstanceLifetime: 4h
```

#### Recommended Test Instances
| Instance Type | Use Case | Cost/Hour | Test Duration |
|--------------|----------|-----------|---------------|
| t3.small | Basic functionality | $0.021 | 30 min = $0.01 |
| t3.medium | CPU workloads | $0.042 | 30 min = $0.02 |
| g5.xlarge | GPU basic test | $1.006 | 10 min = $0.17 |
| p5.48xlarge | GPU validation | $54.00 | 5 min = $4.50 |

**Total estimated cost for Phase 2**: $20-50 with careful testing

---

### Phase 3: EKS Integration Testing
**Goal**: Test ORCA on real AWS Kubernetes (EKS) for production-like environment
**Duration**: Weeks 8-12
**Cost**: ~$150-300/month

#### Why EKS Testing Matters

EKS provides a more realistic test environment:
- **Real AWS networking** - VPC, subnets, security groups as they'll exist in production
- **IAM for Service Accounts** - IRSA for secure credential management
- **Scale testing** - Test with multiple nodes, namespaces, users
- **Production patterns** - LoadBalancers, persistent volumes, service mesh
- **Multi-tenant validation** - Department isolation, cost allocation
- **Networking complexity** - Test pod-to-pod, pod-to-instance communication

#### EKS Setup

```bash
# Create EKS cluster (use Terraform)
cd deploy/terraform/eks-test

# Provision EKS cluster
terraform init
terraform plan
terraform apply

# Outputs:
# - EKS cluster (2 t3.medium nodes)
# - VPC with public/private subnets
# - IAM roles for ORCA
# - Security groups
```

#### EKS Cost Breakdown

**EKS Control Plane**: $73/month (0.10/hour × 730 hours)
**Worker Nodes**:
- 2x t3.medium = $61/month (2 × $0.042/hour × 730 hours)
**Total Base**: ~$134/month

**Burst Testing** (additional):
- Test instances: ~$20-50/month
- Data transfer: ~$5-10/month

**Total Phase 3 Cost**: ~$150-300/month

#### Cost Optimization for Testing

```hcl
# terraform.tfvars - Minimal EKS for testing
cluster_version = "1.28"
node_instance_type = "t3.medium"
node_count = 2  # Minimum for HA
disk_size = 20  # GB

# Enable spot instances for worker nodes (60-70% savings)
node_use_spot = true
```

**With spot instances**: ~$100-150/month

#### What to Test on EKS

**Functional Testing**
- ✅ ORCA deployment via Helm/kubectl
- ✅ Multiple virtual nodes
- ✅ Multi-namespace pod bursting
- ✅ IRSA (IAM Roles for Service Accounts)
- ✅ VPC networking and security groups
- ✅ Persistent volume handling
- ✅ Service discovery (DNS)
- ✅ Load balancer integration

**Scale Testing**
- ✅ 10-50 concurrent burst pods
- ✅ Multiple departments/namespaces
- ✅ Budget enforcement across teams
- ✅ Instance pool management
- ✅ Spot instance interruptions
- ✅ Fast scale-up (0 → 10 instances)
- ✅ Fast scale-down (cleanup)

**Production Validation**
- ✅ High availability (controller restarts)
- ✅ Upgrade testing (rolling updates)
- ✅ Monitoring and alerting
- ✅ Log aggregation (CloudWatch)
- ✅ Cost tracking accuracy
- ✅ Security scanning (IAM, network policies)

#### EKS Test Scenarios

**Scenario 1: Basic Burst**
```bash
# Deploy 5 pods across different instance types
kubectl apply -f examples/eks-test/mixed-workloads.yaml

# Expected: 5 EC2 instances created, pods running
# Cost: ~$0.50 for 30-minute test
```

**Scenario 2: Multi-Department**
```bash
# Create namespaces for biology, cs, physics departments
kubectl apply -f examples/eks-test/multi-tenant.yaml

# Test budget isolation
# Expected: Department budgets enforced separately
```

**Scenario 3: Spot Interruption**
```bash
# Launch spot instances
kubectl apply -f examples/eks-test/spot-workload.yaml

# Simulate interruption (terminate instance manually)
# Expected: ORCA handles gracefully, marks pod as failed
```

**Scenario 4: GPU Workload**
```bash
# Deploy GPU training job
kubectl apply -f examples/eks-test/gpu-training.yaml

# Expected: p5.48xlarge launched, NVIDIA drivers work
# Cost: ~$4.50 for 5-minute validation
```

---

### Phase 4: NRP/SDSU Production Testing
**Goal**: Validate ORCA on actual NRP infrastructure
**Duration**: Weeks 12-16
**Cost**: Covered by NRP/SDSU AWS accounts

#### Setup
- Deploy ORCA to NRP Nautilus cluster
- Configure ORCA to use SDSU AWS account
- Work with SDSU researchers for real workloads

#### What to Test
- ✅ Ceph storage integration
- ✅ NRP namespace/quota integration
- ✅ Real research workloads (AI/ML training)
- ✅ Multi-site awareness
- ✅ Long-running jobs (days)
- ✅ Large-scale GPU jobs (8x H100)
- ✅ User feedback and usability

---

## Testing Matrix

| Test Type | Phase 1 | Phase 2 | Phase 3 | Phase 4 |
|-----------|---------|---------|---------|---------|
| **Environment** | kind | kind + AWS | EKS | NRP Nautilus |
| **Cost** | $0 | $10-50/mo | $150-300/mo | Covered |
| **Duration** | Weeks 1-4 | Weeks 4-8 | Weeks 8-12 | Weeks 12-16 |
| Unit tests | ✅ | ✅ | ✅ | ✅ |
| Integration tests | Mock | Real AWS | Real AWS | Real AWS |
| Scale (pods) | 1-5 | 1-10 | 10-50 | 50-500 |
| GPU testing | Mock | Brief | Full | Production |
| Multi-tenancy | N/A | Single | Multiple | Multiple |
| NRP features | N/A | N/A | N/A | ✅ |

---

## EKS vs kind: When to Use Each

### Use kind for:
- ✅ Daily development
- ✅ Unit/integration tests
- ✅ CI/CD pipeline tests
- ✅ Quick iterations
- ✅ Feature development
- ✅ Pre-commit validation

### Use EKS for:
- ✅ Production-like validation
- ✅ Scale testing
- ✅ Multi-tenancy validation
- ✅ Networking complexity
- ✅ IAM/IRSA testing
- ✅ Performance benchmarking
- ✅ Pre-release validation

---

## Terraform Setup for EKS Testing

We'll create Terraform modules for easy EKS test cluster provisioning:

```
deploy/terraform/
├── eks-test/
│   ├── main.tf           # EKS cluster definition
│   ├── variables.tf      # Configurable parameters
│   ├── outputs.tf        # Cluster endpoints, kubeconfig
│   ├── vpc.tf           # VPC for testing
│   └── iam.tf           # IAM roles for ORCA
```

### Quick Commands

```bash
# Create EKS test cluster
cd deploy/terraform/eks-test
terraform apply

# Deploy ORCA to EKS
kubectl apply -f ../../kubernetes/

# Run tests
kubectl apply -f ../../../examples/eks-test/

# Destroy when done
terraform destroy
```

### Cost-Saving Tips

**1. Use Spot Instances for Workers**
```hcl
node_use_spot = true  # 60-70% savings on worker nodes
```

**2. Auto-Shutdown Overnight**
```bash
# Scale to 0 nodes overnight (save ~$20/day)
kubectl scale deployment --all --replicas=0 -n kube-system
```

**3. Weekend Teardown**
```bash
# Destroy cluster on Friday, recreate Monday
terraform destroy  # Friday 5pm
terraform apply    # Monday 9am
# Saves ~$67/weekend
```

**4. On-Demand Only When Needed**
```bash
# Use EKS only for specific tests
# - Pre-release validation (1 week before release)
# - Scale testing (once per sprint)
# - Production issue reproduction
```

**Realistic EKS Testing Cost**: ~$50-100/month with disciplined usage

---

## Continuous Integration Strategy

### GitHub Actions (Phase 1-2)
```yaml
# .github/workflows/ci.yml
- Unit tests on every PR (kind)
- Integration tests on merge to main (kind + mock AWS)
- Nightly: Integration tests with real AWS (small instances)
```

### EKS Testing (Phase 3)
```yaml
# Manual trigger or weekly schedule
- Deploy to EKS test cluster
- Run full test suite
- Performance benchmarks
- Cost validation
- Teardown after tests
```

---

## Testing Checklist

### Before Each Release

- [ ] Phase 1: All unit tests pass
- [ ] Phase 1: All integration tests pass (mock AWS)
- [ ] Phase 2: AWS integration tests pass (t3.small instances)
- [ ] Phase 2: GPU validation (g5.xlarge, 10 min test)
- [ ] Phase 3: EKS deployment successful
- [ ] Phase 3: Multi-tenant validation
- [ ] Phase 3: Scale test (20 concurrent pods)
- [ ] Phase 3: Cost tracking accuracy verified
- [ ] Documentation updated
- [ ] CHANGELOG updated

---

## Cost Tracking

Track testing costs across all phases:

```bash
# AWS Cost Explorer query
aws ce get-cost-and-usage \
  --time-period Start=2025-01-01,End=2025-01-31 \
  --granularity MONTHLY \
  --metrics BlendedCost \
  --group-by Type=TAG,Key=orca:testing

# Expected monthly costs:
# - Phase 1: $0
# - Phase 2: $10-50
# - Phase 3: $150-300 (or $50-100 with optimizations)
# - Phase 4: Covered by partners
```

---

## Timeline Summary

**Month 1**: Phase 1 (kind + mock AWS) - $0
**Month 2**: Phase 2 (kind + real AWS) - $10-50
**Month 3-4**: Phase 3 (EKS testing) - $150-300 ($50-100 optimized)
**Month 4+**: Phase 4 (NRP production) - Partner covered

**Total estimated testing cost**: $200-400 over 4 months
**With optimizations**: $100-200 over 4 months

---

## Conclusion

The progressive testing strategy allows you to:
1. ✅ Develop quickly and cheaply (kind + mock)
2. ✅ Validate AWS integration (kind + real AWS)
3. ✅ Test production scenarios (EKS)
4. ✅ Deploy to real users (NRP)

EKS testing (Phase 3) is valuable but not required until you're confident in the core functionality. You can develop the entire MVP in Phase 1-2 for under $50, then invest in EKS testing when preparing for production deployment.

**Recommendation**: Start with Phase 1-2, add EKS testing when ready to engage NRP/SDSU for production pilots.
