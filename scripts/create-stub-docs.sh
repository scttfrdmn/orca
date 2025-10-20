#!/bin/bash
# Script to create stub documentation pages for ORCA

set -e

DOCS_DIR="docs"

# Create stub pages with basic content
create_stub() {
    local file="$1"
    local title="$2"
    local description="$3"

    if [ ! -f "$DOCS_DIR/$file" ]; then
        mkdir -p "$(dirname "$DOCS_DIR/$file")"
        cat > "$DOCS_DIR/$file" << EOF
# $title

$description

!!! info "Documentation In Progress"
    This page is under construction. Check back soon for detailed content!

## Coming Soon

This section will cover:

- Key concepts and features
- Step-by-step instructions
- Examples and best practices
- Troubleshooting tips

## Need Help Now?

- Check our [GitHub Issues](https://github.com/scttfrdmn/orca/issues)
- Join the [Discussions](https://github.com/scttfrdmn/orca/discussions)
- Read the [README](https://github.com/scttfrdmn/orca)
EOF
        echo "✅ Created: $file"
    else
        echo "⏭️  Exists: $file"
    fi
}

# Getting Started
create_stub "getting-started/installation.md" "Installation" "Detailed instructions for installing ORCA in your Kubernetes cluster."
create_stub "getting-started/configuration.md" "Configuration" "Configure ORCA for your AWS environment and research workloads."
create_stub "getting-started/first-pod.md" "Your First Bursted Pod" "Deploy your first pod to AWS using ORCA."

# User Guide
create_stub "user-guide/index.md" "User Guide" "Comprehensive guide to using ORCA for research computing workloads."
create_stub "user-guide/instance-selection.md" "Instance Selection" "How to select the right EC2 instance type for your workload."
create_stub "user-guide/gpu-workloads.md" "GPU Workloads" "Running GPU-accelerated AI/ML workloads on AWS."
create_stub "user-guide/spot-instances.md" "Spot Instances" "Using AWS Spot instances for cost-optimized computing."
create_stub "user-guide/cost-management.md" "Cost Management" "Managing costs and budgets for bursted workloads."
create_stub "user-guide/troubleshooting.md" "Troubleshooting" "Common issues and how to resolve them."

# Architecture
create_stub "architecture/index.md" "Architecture" "ORCA's architecture and design principles."
create_stub "architecture/overview.md" "Architecture Overview" "High-level overview of how ORCA works."
create_stub "architecture/instance-selection.md" "Instance Selection Architecture" "How ORCA selects EC2 instances (explicit, template, auto)."
create_stub "architecture/pod-lifecycle.md" "Pod Lifecycle" "How ORCA manages pod lifecycle from creation to termination."
create_stub "architecture/aws-integration.md" "AWS Integration" "How ORCA integrates with AWS EC2 and other services."
create_stub "architecture/design-decisions.md" "Design Decisions" "Key architectural decisions and trade-offs."

# Development
create_stub "development/index.md" "Development Guide" "Guide for contributing to ORCA development."
create_stub "development/building.md" "Building ORCA" "How to build ORCA from source."
create_stub "development/contributing.md" "Contributing" "How to contribute to ORCA."
create_stub "development/code-style.md" "Code Style" "Coding standards and style guide for ORCA."

# API Reference
create_stub "api/index.md" "API Reference" "Reference documentation for ORCA annotations, configuration, and metrics."
create_stub "api/annotations.md" "Pod Annotations" "Complete reference of ORCA pod annotations."
create_stub "api/configuration.md" "Configuration Reference" "Complete configuration options for ORCA."
create_stub "api/metrics.md" "Metrics Reference" "Prometheus metrics exposed by ORCA."

# Community
create_stub "community/index.md" "Community" "Join the ORCA community."
create_stub "community/support.md" "Support" "How to get help with ORCA."
create_stub "community/roadmap.md" "Roadmap" "ORCA's development roadmap and future plans."
create_stub "community/license.md" "License" "ORCA's Apache 2.0 license."

# Copy CONTRIBUTING.md to docs
if [ -f "CONTRIBUTING.md" ]; then
    cp CONTRIBUTING.md "$DOCS_DIR/CONTRIBUTING.md"
    echo "✅ Copied: CONTRIBUTING.md"
fi

echo ""
echo "✅ All stub documentation pages created!"
echo "Run 'make docs-serve' to preview the site."
