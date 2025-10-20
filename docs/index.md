# ORCA - Orchestration for Research Cloud Access

<div align="center">
  <img src="images/orca-400.png" alt="ORCA Logo" width="400"/>
</div>

**ORCA** enables research institutions to seamlessly burst Kubernetes workloads from on-premises Kubernetes clusters to AWS, with native support for GPU-intensive AI/ML computing.

## What is ORCA?

ORCA (Orchestration for Research Cloud Access) is a Kubernetes Virtual Kubelet provider that allows research Kubernetes clusters to dynamically extend capacity to AWS when local resources are exhausted.

## Key Features

- 🎓 **Research-First Design** - Built for academic and research workloads
- 🖥️ **AI/ML Accelerators** - Support for NVIDIA GPUs (P6, P5, P4d, G6e), AWS Trainium, Inferentia, and FPGAs
- 🎯 **Explicit Control** - Users specify exact instance types, not guessed
- 💰 **Cost-Aware** - Budget controls, cost tracking, spot instance support
- 🔓 **Open Source** - Apache 2.0 licensed, community-driven

## Quick Links

<div class="grid cards" markdown>

-   :material-rocket-launch:{ .lg .middle } __Getting Started__

    ---

    Get ORCA up and running in minutes

    [:octicons-arrow-right-24: Quick Start](getting-started/quick-start.md)

-   :material-book-open-variant:{ .lg .middle } __User Guide__

    ---

    Learn how to use ORCA for your workloads

    [:octicons-arrow-right-24: User Guide](user-guide/index.md)

-   :material-sitemap:{ .lg .middle } __Architecture__

    ---

    Understand how ORCA works

    [:octicons-arrow-right-24: Architecture](architecture/index.md)

-   :material-code-braces:{ .lg .middle } __Development__

    ---

    Contribute to ORCA development

    [:octicons-arrow-right-24: Development](development/index.md)

</div>

## Architecture Overview

```mermaid
graph TB
    subgraph "Research Cluster"
        K8S[Kubernetes API]
        VK[ORCA Virtual Kubelet]
        POD[Pod with GPU Request]
    end

    subgraph "AWS"
        EC2[EC2 Instance<br/>P5.48xlarge<br/>8x H100 GPUs]
        SPOT[Spot Instances]
        CR[Capacity Reservations]
    end

    POD -->|Schedule| VK
    VK -->|Register| K8S
    VK -->|Launch| EC2
    VK -.->|Optional| SPOT
    VK -.->|Preferred| CR

    style VK fill:#4285f4,stroke:#333,stroke-width:2px,color:#fff
    style EC2 fill:#ff9900,stroke:#333,stroke-width:2px
    style POD fill:#326ce5,stroke:#333,stroke-width:2px,color:#fff
```

## Use Cases

### AI/ML Training
Burst large model training to AWS GPUs, Trainium, or Inferentia when local clusters are full.

### Cost-Optimized Computing
Use Trainium for 50% lower training costs or Inferentia for 70% lower inference costs compared to GPUs.

### Research Computing
Access specialized hardware on-demand: FPGAs for genomics, latest GPUs for deep learning.

### Multi-Tenant Research
Support multiple departments with separate budgets and cost tracking.

## Why ORCA?

### vs. Elotl Kip
- **Kip is EOL** (last updated 2021) - stuck on K8s 1.18, AWS SDK v1
- **ORCA is modern** - K8s 1.34, AWS SDK v2, Go 1.25, latest instance types (P6, G6e)
- **ORCA prioritizes explicit control** - users know their requirements

### vs. AWS Fargate Virtual Kubelet
- **Fargate provider is unmaintained** and doesn't support GPUs
- **ORCA is GPU-first** - built for AI/ML research

### vs. Building on Managed K8s
- **ORCA extends existing Kubernetes clusters** - research institutions already have K8s
- **No migration needed** - burst workloads, keep existing infrastructure

## Project Status

**Current Phase**: Active Development (v0.1.0-dev)

- ✅ Core architecture designed and implemented
- ✅ AWS EC2 integration complete
- ✅ Instance selection (explicit, template, auto)
- ✅ Virtual Kubelet integration
- 🚧 Container runtime integration (in progress)
- ⏳ GPU capacity reservations (v0.2.0)
- ⏳ kubectl logs/exec (v0.2.0)

## Roadmap

ORCA development follows a phased approach aligned with research computing needs. Track our progress on the [GitHub project board](https://github.com/scttfrdmn/orca/projects/3).

### [Phase 1: MVP](https://github.com/scttfrdmn/orca/milestone/1) ✅ Complete
**Months 1-3**

Core Virtual Kubelet provider with basic pod-to-EC2 mapping and explicit instance selection. Simple lifecycle management.

**Status**: Implementation complete, metrics in progress

### [Phase 2: Production Features](https://github.com/scttfrdmn/orca/milestone/2) 🚧 In Progress
**Months 4-6**

Production-ready features including:
- GPU support for all NVIDIA instance types (P6, P5, P4d, G6e)
- Container runtime integration with containerd
- kubectl logs and exec via CloudWatch and Systems Manager
- Spot instance support for cost optimization

[View Phase 2 issues →](https://github.com/scttfrdmn/orca/milestone/2)

### [Phase 3: NRP Integration](https://github.com/scttfrdmn/orca/milestone/4) ⏳ Planned
**Months 7-9**

National Research Platform integration:
- Automatic Ceph storage mounting
- NRP namespace awareness and identity
- Multi-tenancy with per-namespace quotas
- Cost tracking and budget enforcement

[View Phase 3 issues →](https://github.com/scttfrdmn/orca/milestone/4)

### [Phase 4: Advanced Features](https://github.com/scttfrdmn/orca/milestone/6) ⏳ Future
**Months 9+**

Enterprise and advanced capabilities:
- Intelligent scheduling algorithms
- Capacity planning and forecasting
- Compliance features (HIPAA, FedRAMP)
- Multi-region support

[View Phase 4 issues →](https://github.com/scttfrdmn/orca/milestone/6)

[View all milestones →](https://github.com/scttfrdmn/orca/milestones)

## Community

- **Website**: [orcapod.dev](https://orcapod.dev)
- **GitHub**: [scttfrdmn/orca](https://github.com/scttfrdmn/orca)
- **Issues**: [Report bugs or request features](https://github.com/scttfrdmn/orca/issues)
- **Discussions**: [Questions and ideas](https://github.com/scttfrdmn/orca/discussions)
- **License**: [Apache 2.0](https://github.com/scttfrdmn/orca/blob/main/LICENSE)

## Getting Help

- 📖 [Read the docs](getting-started/index.md)
- 🐛 [Report issues](https://github.com/scttfrdmn/orca/issues)
- 💬 [Discussions](https://github.com/scttfrdmn/orca/discussions)
- 🤝 [Contributing guide](CONTRIBUTING.md)

---

*Built with 🌊 for research computing*
