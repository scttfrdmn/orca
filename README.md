# ORCA - Orchestration for Research Cloud Access

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/orca)](https://goreportcard.com/report/github.com/scttfrdmn/orca)

**ORCA** enables research institutions to seamlessly burst Kubernetes workloads from on-premises clusters to AWS, with native support for GPU-intensive AI/ML computing.

---

## Overview

ORCA (Orchestration for Research Cloud Access) is a Kubernetes Virtual Kubelet provider that allows research computing clusters to dynamically extend capacity to AWS when local resources are exhausted.

### Key Features

- 🎓 **Research-First Design** - Built for academic and research workloads
- 🖥️ **GPU-Optimized** - First-class support for AI/ML with latest NVIDIA GPUs (P6, P5, P4d, G6e, G6, etc.)
- 🎯 **Explicit Control** - Users specify exact instance types, not guessed
- 💰 **Cost-Aware** - Budget controls, cost tracking, spot instance support
- 🔓 **Open Source** - Apache 2.0 licensed, community-driven

### ⚠️ Important: AWS Capacity Reservations

**Latest GPU instances (P6, P5, P4de, P4d, G6e) are virtually unavailable without AWS Capacity Reservations.** This is the reality of GPU availability on AWS in 2025, not a limitation of ORCA.

- **P6.48xlarge (Blackwell B200)**: Capacity Reservations required
- **P5.48xlarge (H100 80GB)**: Effectively requires Capacity Reservations
- **P5e.48xlarge (H200)**: Capacity Reservations required
- **P4de.24xlarge (A100 80GB)**: Reservations required in most regions
- **P4d.24xlarge (A100 40GB)**: Extremely limited on-demand availability
- **G6e (L40S)**: Better availability but still constrained

**Without Reservations**: Expect frequent `InsufficientInstanceCapacity` errors for modern GPUs.

ORCA v0.2.0+ will support targeting Capacity Reservations. See [docs/CAPACITY-RESERVATIONS.md](docs/CAPACITY-RESERVATIONS.md) for details.

### Use Cases

- **AI/ML Training** - Burst large model training to AWS GPUs
- **Batch Processing** - Scale compute-intensive jobs elastically
- **Research Computing** - Access specialized hardware on-demand
- **Multi-Tenant** - Support multiple departments with separate budgets

---

## Architecture

```text
┌──────────────────────────────────────────┐
│  Research Cluster (e.g., NRP Nautilus)  │
│                                          │
│  ┌────────────────────────────────────┐ │
│  │  ORCA Virtual Node                 │ │
│  │  (Virtual Kubelet Provider)        │ │
│  └──────────────┬─────────────────────┘ │
└─────────────────┼────────────────────────┘
                  │
                  │ Secure Connection
                  │
┌─────────────────▼────────────────────────┐
│  AWS Account                             │
│  ┌────────────────────────────────────┐  │
│  │  VPC                               │  │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ │  │
│  │  │ EC2    │ │ EC2    │ │ EC2    │ │  │
│  │  │ (GPU)  │ │ (GPU)  │ │ (GPU)  │ │  │
│  │  └────────┘ └────────┘ └────────┘ │  │
│  └────────────────────────────────────┘  │
└──────────────────────────────────────────┘
```

When a pod with specific tolerations is scheduled:
1. ORCA intercepts the pod creation
2. Launches a right-sized EC2 instance (explicit user choice)
3. Runs the pod workload on the instance
4. Terminates the instance when the pod completes

---

## Quick Start

### Prerequisites

- Kubernetes cluster (1.28+)
- AWS account with appropriate permissions
- Go 1.21+ (for building from source)

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/scttfrdmn/orca.git
cd orca

# 2. Edit deploy/configmap.yaml with your AWS settings
#    - region, vpcID, subnetID, securityGroupIDs

# 3. Configure AWS credentials (IRSA recommended, or edit deploy/secret.yaml)

# 4. Deploy to your cluster
kubectl apply -k deploy/

# 5. Verify installation
kubectl get nodes  # Should see orca-aws-node
kubectl get pods -n orca-system
```

See [deploy/README.md](deploy/README.md) for detailed installation instructions.

### Example: Burst a GPU Job to AWS

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: gpu-training-job
  annotations:
    # Explicit instance selection - you control exactly what you get
    orca.research/instance-type: "p5.48xlarge"  # 8x H100 GPUs
    orca.research/launch-type: "spot"           # 70% cost savings
spec:
  nodeSelector:
    orca.research/provider: "aws"
  tolerations:
  - key: orca.research/burst-node
    operator: Equal
    value: "true"
    effect: NoSchedule
  containers:
  - name: trainer
    image: pytorch/pytorch:2.1.0-cuda12.1-cudnn8-runtime
    command: ["python", "train.py"]
    resources:
      limits:
        nvidia.com/gpu: 8
        memory: "1024Gi"
        cpu: "192"
```

Apply and watch it burst to AWS:
```bash
kubectl apply -f my-gpu-job.yaml
kubectl get pods -w
```

---

## Project Status

**Current Phase**: Design & Planning

- ✅ Architecture design
- ✅ Requirements gathering
- 🚧 Core implementation (in progress)
- ⏳ Alpha release (planned)

This is a **greenfield project** built from scratch with modern tools and patterns, learning from but not forking existing Virtual Kubelet providers.

---

## Design Principles

### 1. Explicit Over Implicit
Users specify exactly what they want (instance types, launch types) rather than relying on "smart" auto-selection that often gets it wrong.

### 2. Research-First
Optimized for research computing workflows, not generic enterprise workloads. GPU support, batch jobs, cost awareness.

### 3. Production-Grade
Idiomatic Go, comprehensive testing, proper error handling, observability built-in from day one.

### 4. Learn from History
Studied issues and limitations of prior Virtual Kubelet providers (especially Kip) to avoid known pitfalls.

### 5. NRP Integration
Native support for National Research Platform features like Ceph storage, multi-site awareness, quota integration.

---

## Why ORCA?

### vs. Elotl Kip
- **Kip is EOL** (last updated 2021) - stuck on K8s 1.18, AWS SDK v1
- **ORCA is modern** - K8s 1.31, AWS SDK v2, Go 1.21, latest instance types (P6, G6e)
- **Kip auto-selects instances** - often wrong for research workloads
- **ORCA prioritizes explicit control** - users know their requirements

### vs. AWS Fargate Virtual Kubelet
- **Fargate provider is unmaintained** and doesn't support GPUs
- **ORCA is GPU-first** - built for AI/ML research

### vs. Building on Managed K8s
- **EKS requires managing worker nodes** or using Fargate (no GPUs)
- **ORCA extends existing clusters** - research institutions already have K8s
- **No migration needed** - burst workloads, keep existing infrastructure

---

## Roadmap

### Phase 1: MVP (Months 1-3)
- Core Virtual Kubelet provider
- Basic pod-to-EC2 mapping
- Explicit instance selection
- Simple lifecycle management

### Phase 2: Production Features (Months 4-6)
- GPU support (all NVIDIA types)
- Container runtime integration
- kubectl logs and exec
- Spot instance handling

### Phase 3: NRP Integration (Months 7-9)
- Ceph storage auto-mounting
- NRP namespace awareness
- Multi-tenancy support
- Cost tracking and budgets

### Phase 4: Advanced Features (Months 9+)
- Advanced scheduling algorithms
- Capacity planning
- Compliance features (HIPAA, FedRAMP)
- Multi-region support

---

## Contributing

ORCA is open source and welcomes contributions! We especially value:

- 🐛 Bug reports and fixes
- 📖 Documentation improvements
- 🧪 Testing and validation
- 💡 Feature requests from research users
- 🎓 University partnerships

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## Community

- **GitHub Issues**: Bug reports and feature requests
- **Discussions**: Design discussions and questions
- **Research Partners**: NRP, SDSU, and other universities

---

## Architecture & Development

See detailed documentation:
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - System design
- [DEVELOPMENT.md](docs/DEVELOPMENT.md) - Building and testing
- [INSTANCE_SELECTION.md](docs/INSTANCE_SELECTION.md) - How instance selection works
- [KIP_LESSONS.md](docs/KIP_LESSONS.md) - Lessons learned from Kip

---

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

Copyright 2025 Scott Friedman and ORCA Contributors

## Disclaimer

This project is not officially associated with or supported by Kip, Amazon Web Services (AWS), the National Research Platform (NRP), or San Diego State University (SDSU). ORCA is an independent open-source project developed for research computing purposes

---

## Acknowledgments

- **Elotl Kip** - Reference architecture and inspiration
- **Virtual Kubelet** - Provider framework
- **National Research Platform** - Requirements and testing
- **San Diego State University** - Partnership and feedback

---

## Contact

- **Maintainer**: Scott Friedman (AWS)
- **Research Partners**: SDSU, NRP
- **Issues**: https://github.com/scttfrdmn/orca/issues

---

*Built with 🌊 for research computing*
