#!/bin/bash
# Script to create GitHub issues for ORCA based on TODO.md
# Organized by milestones/phases like CloudWorkstation

set -e

REPO="scttfrdmn/orca"

echo "🎯 Creating ORCA Project Issues from TODO.md"
echo "=============================================="
echo ""

# Phase 0.1: MVP - Core Functionality
echo "📋 Phase 0.1: MVP - Core Functionality"
echo ""

gh issue create --repo "$REPO" \
  --title "[MVP] Container Runtime Integration" \
  --label "enhancement,priority: critical,area: provider,milestone: v0.1-mvp" \
  --milestone "v0.1-mvp" \
  --body "## Context
ORCA creates EC2 instances but doesn't yet run actual containers. Need to integrate containerd to execute pod workloads.

## Tasks from TODO.md
Located at: \`pkg/provider/provider.go:346-374\`

- [ ] Add user data script to install containerd on launch
- [ ] Configure containerd with proper CNI networking
- [ ] Pull container images from registry
- [ ] Run containers with resource limits (CPU, memory, GPU)
- [ ] Stream container logs to CloudWatch
- [ ] Monitor container health
- [ ] Report container state back to Kubernetes

## Acceptance Criteria
- [ ] EC2 instances have containerd pre-installed
- [ ] Containers run with correct env vars and mounts
- [ ] GPU pass-through works for GPU instances
- [ ] Container state syncs with pod status

## Priority
**Critical** - Blocking MVP functionality

## Related
Part of TODO.md \"Critical - Required for MVP\" section"

gh issue create --repo "$REPO" \
  --title "[MVP] kubectl logs Support via CloudWatch" \
  --label "enhancement,priority: critical,area: provider,milestone: v0.1-mvp" \
  --milestone "v0.1-mvp" \
  --body "## Context
Implement GetContainerLogs provider method to enable \`kubectl logs\`.

## Tasks from TODO.md
Located at: \`pkg/provider/provider.go:346-356\`

- [ ] Implement GetContainerLogs via CloudWatch Logs
- [ ] Stream container stdout/stderr to CloudWatch
- [ ] Support --tail, --since, --follow flags
- [ ] Handle multi-container pods
- [ ] Implement log streaming for -f flag

## Implementation
- Use CloudWatch Logs agent on EC2
- Log group: \`/orca/{namespace}/{pod}\`
- Log stream per container
- CloudWatch Insights for querying

## Acceptance Criteria
- [ ] \`kubectl logs <pod>\` returns logs
- [ ] \`kubectl logs -f <pod>\` streams in real-time
- [ ] \`kubectl logs --tail=100 <pod>\` works
- [ ] Logs persist after pod deletion (configurable)

## Priority
**Critical** - Essential for debugging

## Location
\`pkg/provider/provider.go:346-356\`"

gh issue create --repo "$REPO" \
  --title "[MVP] kubectl exec Support via Systems Manager" \
  --label "enhancement,priority: high,area: provider,milestone: v0.1-mvp" \
  --milestone "v0.1-mvp" \
  --body "## Context
Implement RunInContainer provider method to enable \`kubectl exec\`.

## Tasks from TODO.md
Located at: \`pkg/provider/provider.go:352-356\`

- [ ] Implement RunInContainer via AWS Systems Manager
- [ ] Install SSM agent on EC2 instances
- [ ] Execute commands in running containers
- [ ] Support interactive sessions (TTY)
- [ ] Handle stdin/stdout/stderr streaming

## Implementation
- AWS Systems Manager Session Manager
- SSM agent in EC2 user data
- Bridge to containerd exec API
- Proper IAM permissions

## Acceptance Criteria
- [ ] \`kubectl exec <pod> -- command\` works
- [ ] \`kubectl exec -it <pod> -- bash\` provides interactive shell
- [ ] stdin/stdout/stderr work correctly
- [ ] TTY and non-TTY modes both supported

## Priority
**High** - Important for debugging and development

## Location
\`pkg/provider/provider.go:352-356\`"

gh issue create --repo "$REPO" \
  --title "[MVP] Metrics and Statistics via CloudWatch" \
  --label "enhancement,priority: medium,area: provider,milestone: v0.1-mvp" \
  --milestone "v0.1-mvp" \
  --body "## Context
Implement GetStatsSummary for resource usage metrics.

## Tasks from TODO.md
Located at: \`pkg/provider/provider.go:358-374\`

- [ ] Implement GetStatsSummary via CloudWatch metrics
- [ ] Collect CPU, memory, network, disk metrics
- [ ] Expose GPU utilization for GPU instances
- [ ] Integrate with \`kubectl top\` command
- [ ] Add custom ORCA metrics

## Implementation
- CloudWatch agent on EC2 instances
- Publish to CloudWatch Metrics
- Query via AWS SDK
- Cache for performance

## Acceptance Criteria
- [ ] \`kubectl top pod <pod>\` returns metrics
- [ ] \`kubectl top node orca-aws-node\` works
- [ ] GPU metrics visible for GPU pods
- [ ] Metrics update every 60 seconds

## Priority
**Medium** - Useful for monitoring

## Location
\`pkg/provider/provider.go:358-374\`"

echo ""
echo "📋 Phase 0.2: GPU Features"
echo ""

gh issue create --repo "$REPO" \
  --title "[GPU] AWS Capacity Reservations Support" \
  --label "enhancement,priority: critical,area: aws,milestone: v0.2-gpu,aws: capacity-reservations" \
  --milestone "v0.2-gpu" \
  --body "## Context
🚨 **ESSENTIAL FOR GPU VIABILITY**

Modern NVIDIA GPU instances (P5, P4d, P4de, G6e) are virtually unavailable without Capacity Reservations. This is THE #1 blocker for GPU research.

## Reality Check (October 2025)
- **P6.48xlarge** (B200): Virtually impossible without reservations
- **P5.48xlarge** (H100): Virtually impossible without reservations
- **P5e.48xlarge** (H200): Reservations required
- **P4de.24xlarge** (A100 80GB): Reservations required in most regions
- **P4d.24xlarge** (A100 40GB): Extremely limited on-demand
- **G6e** (L40S): Constrained during peak

**Without this feature, researchers CANNOT use ORCA for modern GPU workloads.**

## v0.2.0 Tasks (CRITICAL)
- [ ] Add \`orca.research/capacity-reservation-id\` annotation
- [ ] Add \`orca.research/capacity-reservation-preference\` annotation (open|targeted)
- [ ] Pass CapacityReservationSpecification to RunInstances
- [ ] Add config option for default reservation preference
- [ ] Fail gracefully with clear error messages
- [ ] Document ODCR creation and configuration

## Example Usage
\`\`\`yaml
metadata:
  annotations:
    orca.research/instance-type: \"p5.48xlarge\"
    orca.research/capacity-reservation-id: \"cr-0123456789abcdef\"
\`\`\`

## v0.3.0 Future Tasks
- [ ] Automatic ODCR discovery
- [ ] Prefer reserved capacity automatically
- [ ] Capacity Block for ML support
- [ ] Reservation utilization metrics

## Priority
**CRITICAL** - Without this, ORCA cannot support modern GPU research

## References
- TODO.md: \"Critical - GPU/ML Features\" section
- docs/CAPACITY-RESERVATIONS.md"

gh issue create --repo "$REPO" \
  --title "[GPU] Instance Type Fallback for GPU Availability" \
  --label "enhancement,priority: high,area: instances,milestone: v0.2-gpu" \
  --milestone "v0.2-gpu" \
  --body "## Context
When a GPU instance type is unavailable, automatically try alternatives.

## Tasks from TODO.md
Located in \"Instance Management\" section

- [ ] Add fallback instance type configuration
- [ ] Example: p5.48xlarge → p4de.24xlarge → p4d.24xlarge
- [ ] Respect GPU count requirements
- [ ] Log fallback attempts
- [ ] Report which instance type was actually used

## Configuration Example
\`\`\`yaml
instances:
  fallbacks:
    p5.48xlarge: [p4de.24xlarge, p4d.24xlarge]
    g6e.48xlarge: [g5.48xlarge, p3.16xlarge]
\`\`\`

## Acceptance Criteria
- [ ] Fallback configured per instance type
- [ ] Tries alternatives in order
- [ ] Respects resource requirements (GPU count, memory)
- [ ] Logs fallback attempts clearly
- [ ] Pod annotation shows actual instance used

## Priority
**High** - Improves GPU availability success rate"

echo ""
echo "📋 Phase 0.3: Production Readiness"
echo ""

gh issue create --repo "$REPO" \
  --title "[Testing] Add Provider Unit Tests" \
  --label "technical-debt,priority: high,area: tests,milestone: v0.3-production" \
  --milestone "v0.3-production" \
  --body "## Context
The provider package has no unit tests. Need comprehensive test coverage.

## Tasks from TODO.md
Testing section: \"pkg/provider: ❌ No tests\"

- [ ] Test pod creation lifecycle
- [ ] Test pod deletion and cleanup
- [ ] Test status updates
- [ ] Test error handling
- [ ] Mock AWS client interactions
- [ ] Test concurrent pod operations
- [ ] Test pod state transitions

## Target Coverage
Aim for >80% coverage on pkg/provider

## Acceptance Criteria
- [ ] CreatePod tests with mocks
- [ ] DeletePod tests with cleanup verification
- [ ] GetPodStatus tests with state sync
- [ ] Error case coverage
- [ ] Race condition tests with -race flag

## Priority
**High** - Required for production confidence"

gh issue create --repo "$REPO" \
  --title "[CI/CD] GitHub Actions for Testing and Building" \
  --label "technical-debt,priority: high,area: build,milestone: v0.3-production" \
  --milestone "v0.3-production" \
  --body "## Context
Set up comprehensive CI/CD pipeline.

## Tasks from TODO.md
Build & CI/CD section

- [ ] Add test workflow (unit + integration)
- [ ] Add lint workflow (golangci-lint)
- [ ] Add build workflow
- [ ] Add container image build/push
- [ ] Add coverage reporting (Codecov)
- [ ] Add security scanning (gosec, snyk)
- [ ] Matrix testing (Go 1.23, 1.24, 1.25)
- [ ] Integration with GitHub releases

## Workflows Needed
1. \`.github/workflows/test.yml\`
2. \`.github/workflows/lint.yml\`
3. \`.github/workflows/build.yml\`
4. \`.github/workflows/release.yml\` (GoReleaser)

## Acceptance Criteria
- [ ] Tests run on every PR
- [ ] Linting enforced
- [ ] Container images built on main
- [ ] Releases automated via tags

## Priority
**High** - Essential for production quality"

gh issue create --repo "$REPO" \
  --title "[Security] Implement IRSA for AWS Credentials" \
  --label "enhancement,priority: high,area: aws,security,milestone: v0.3-production" \
  --milestone "v0.3-production" \
  --body "## Context
Use IAM Roles for Service Accounts instead of static credentials.

## Tasks from TODO.md
Security section

- [ ] Create IAM role with ORCA permissions
- [ ] Configure trust policy for IRSA
- [ ] Update deployment to use IRSA
- [ ] Remove static credential support (deprecated)
- [ ] Document IRSA setup process
- [ ] Add troubleshooting guide

## Benefits
- No static credentials in cluster
- Automatic credential rotation
- Least privilege access
- AWS CloudTrail attribution

## Acceptance Criteria
- [ ] ORCA pod uses IRSA
- [ ] No AWS credentials in secrets
- [ ] Documentation updated
- [ ] Example terraform for IAM role

## Priority
**High** - Security best practice"

echo ""
echo "📋 Phase 0.4: NRP Integration"
echo ""

gh issue create --repo "$REPO" \
  --title "[NRP] Namespace-Level Budget Tracking" \
  --label "enhancement,priority: medium,area: provider,milestone: v0.4-nrp" \
  --milestone "v0.4-nrp" \
  --body "## Context
Enable multi-tenant cost tracking for research groups.

## Tasks from TODO.md
Multi-tenancy and Cost Management sections

- [ ] Track costs per namespace
- [ ] Add budget limits per namespace
- [ ] Alert when approaching budget
- [ ] Block new pods when over budget
- [ ] Cost attribution via tags
- [ ] Generate monthly reports

## Configuration
\`\`\`yaml
limits:
  namespaceQuotas:
    research-group-a:
      monthlyBudget: 5000.00
      maxInstances: 10
      maxGPUs: 8
\`\`\`

## Acceptance Criteria
- [ ] Per-namespace cost tracking
- [ ] Budget enforcement
- [ ] Alerts via annotations or events
- [ ] Usage dashboards

## Priority
**Medium** - Important for NRP multi-tenancy"

echo ""
echo "📋 Documentation & Good First Issues"
echo ""

gh issue create --repo "$REPO" \
  --title "[Docs] Add Architecture Diagrams" \
  --label "documentation,priority: medium,area: docs,good first issue" \
  --milestone "v1.0" \
  --body "## Context
Add visual diagrams to explain ORCA architecture.

## Tasks from TODO.md
Documentation section

- [ ] Pod lifecycle sequence diagram
- [ ] EC2 instance provisioning flow
- [ ] Network architecture diagram
- [ ] Virtual Kubelet integration diagram
- [ ] Cost tracking flow
- [ ] Multi-tenancy architecture

## Tools
- Mermaid diagrams (already supported in MkDocs)
- Draw.io/Lucidchart for complex diagrams
- Place in \`docs/architecture/\`

## Acceptance Criteria
- [ ] At least 4 diagrams added
- [ ] Diagrams render correctly in docs
- [ ] Diagrams are clear and helpful
- [ ] Source files included for editing

## Priority
**Medium** - Helps users understand ORCA

## Good First Issue
Great for visual-minded contributors!"

gh issue create --repo "$REPO" \
  --title "[Docs] Complete API Reference Documentation" \
  --label "documentation,priority: medium,area: docs,good first issue" \
  --milestone "v1.0" \
  --body "## Context
Document all ORCA annotations and configuration options.

## Tasks from TODO.md
Documentation section: \"Document all annotations\"

- [ ] Complete \`docs/api/annotations.md\` with all pod annotations
- [ ] Complete \`docs/api/configuration.md\` with all config options
- [ ] Complete \`docs/api/metrics.md\` with all Prometheus metrics
- [ ] Add examples for each option
- [ ] Add validation rules

## Annotations to Document
- \`orca.research/instance-type\`
- \`orca.research/launch-type\`
- \`orca.research/capacity-reservation-id\`
- \`orca.research/spot-max-price\`
- \`orca.research/workload-template\`

## Acceptance Criteria
- [ ] All annotations documented
- [ ] All config fields documented
- [ ] Examples provided
- [ ] Validation rules explained

## Priority
**Medium** - Essential reference

## Good First Issue
No code required, just technical writing!"

gh issue create --repo "$REPO" \
  --title "[Quality] Add godoc Comments to All Exported Functions" \
  --label "technical-debt,priority: low,area: tests,good first issue" \
  --milestone "v1.0" \
  --body "## Context
Improve code documentation with comprehensive godoc comments.

## Tasks from TODO.md
Code Quality section

- [ ] Add godoc to pkg/provider exports
- [ ] Add godoc to pkg/instances exports
- [ ] Add godoc to pkg/config exports
- [ ] Add godoc to pkg/node exports
- [ ] Add examples in godoc
- [ ] Generate godoc.org documentation

## Example Format
\`\`\`go
// CreatePod creates a new pod by launching an EC2 instance.
// It selects an appropriate instance type, creates the instance,
// and monitors its status until running.
//
// Example:
//    err := provider.CreatePod(ctx, pod)
func (p *OrcaProvider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
\`\`\`

## Acceptance Criteria
- [ ] All exported functions have godoc
- [ ] Examples included where helpful
- [ ] godoc.org renders correctly
- [ ] Go Report Card passes

## Priority
**Low** - Quality improvement

## Good First Issue
Great for learning the codebase!"

echo ""
echo "✅ Issues created successfully!"
echo ""
echo "📊 Summary:"
echo "   - Phase 0.1 (MVP): 4 issues"
echo "   - Phase 0.2 (GPU): 2 issues"
echo "   - Phase 0.3 (Production): 3 issues"
echo "   - Phase 0.4 (NRP): 1 issue"
echo "   - Documentation: 2 issues"
echo "   - Good First Issues: 1 issue"
echo ""
echo "🔗 View all issues: https://github.com/$REPO/issues"
echo "🎯 View project: https://github.com/users/scttfrdmn/projects/3"
