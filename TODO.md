# TODO List - Technical Debt and Pending Work

This document tracks pending work, TODOs, and technical debt in the ORCA project.

## Critical - Required for MVP

### Virtual Kubelet Integration
- [ ] Add Virtual Kubelet library as dependency
- [ ] Implement node controller to register with Kubernetes API server
- [ ] Implement pod watch loop for pod lifecycle events
- [ ] Handle pod status updates and sync with EC2 instance state
- [ ] Implement node heartbeat/status updates

**Location:** `cmd/orca/main.go:95-98`

### Structured Logging
- [ ] Replace placeholder `log()` function with structured logging library
- [ ] Options: logrus, zap, or zerolog
- [ ] Implement proper log levels and formatting (JSON/text)
- [ ] Add context-aware logging throughout the codebase

**Location:** `cmd/orca/main.go:110-133`

### AWS Client Features
- [ ] Implement `GetContainerLogs` via CloudWatch Logs or Systems Manager
- [ ] Implement `RunInContainer` via AWS Systems Manager Session Manager
- [ ] Implement `GetStatsSummary` via CloudWatch metrics
- [ ] Add proper timeout handling for all AWS operations

**Locations:**
- `pkg/provider/provider.go:346-356` (GetContainerLogs)
- `pkg/provider/provider.go:352-356` (RunInContainer)
- `pkg/provider/provider.go:358-374` (GetStatsSummary)

## High Priority

### Testing
- [x] Add unit tests for `internal/aws/client.go` with AWS SDK mocks
- [x] Create integration tests using LocalStack
- [ ] Add unit tests for `pkg/provider` with mocks
- [ ] Add unit tests for `pkg/config` with validation tests
- [ ] Add end-to-end smoke tests
- [ ] Set up test coverage reporting

**Current Test Coverage:**
- `pkg/instances`: ✅ 21 unit tests passing (explicit, template, auto selectors)
- `internal/aws`: ✅ 30 unit tests passing (all helper methods + EC2 operations)
- `internal/aws`: ✅ 2 integration tests passing (LocalStack: create/describe/terminate, spot instances)
- `pkg/provider`: ❌ No tests
- `pkg/config`: ❌ No tests

### Deployment
- [ ] Create Kubernetes deployment manifest
- [ ] Create ServiceAccount, Role, RoleBinding manifests
- [ ] Create ConfigMap for configuration
- [ ] Add Helm chart for easy deployment
- [ ] Document deployment procedures

### Configuration Enhancements
- [ ] Add kubeconfig flag support when Virtual Kubelet is integrated
- [ ] Add support for multiple AWS profiles
- [ ] Add support for cross-account roles
- [ ] Validate all instance types against actual AWS availability

**Location:** `cmd/orca/main.go:33-34` (kubeconfig flag commented out)

## Medium Priority

### Monitoring & Metrics
- [ ] Implement Prometheus metrics endpoint
- [ ] Add metrics for instance lifecycle events
- [ ] Add metrics for pod scheduling decisions
- [ ] Add metrics for AWS API call latency and errors
- [ ] Add cost tracking metrics

**Location:** `pkg/provider/provider.go` (metrics integration points throughout)

### Cost Management
- [ ] Implement budget tracking and enforcement
- [ ] Add cost estimation before instance launch
- [ ] Implement spot instance price checking
- [ ] Add usage reports and cost analysis
- [ ] Implement automatic cleanup for exceeded budgets

**Related:** `pkg/config/config.go:77-82` (LimitsConfig)

### Multi-tenancy Features
- [ ] Implement namespace-level budget tracking
- [ ] Add department/team tagging and cost attribution
- [ ] Implement resource quotas per namespace
- [ ] Add usage dashboards per team

### Instance Management
- [ ] Add support for multiple subnets/AZs
- [ ] Implement instance type fallback (e.g., if p5.48xlarge unavailable, try p4de.24xlarge)
- [ ] Add spot instance interruption handling
- [ ] Implement graceful instance termination with pod eviction
- [ ] Add support for persistent EBS volumes

### Security
- [ ] Implement IAM role for service accounts (IRSA)
- [ ] Add support for pod security policies
- [ ] Implement secrets management integration (AWS Secrets Manager)
- [ ] Add audit logging for all AWS operations
- [ ] Add network policy support

## Low Priority

### Documentation
- [ ] Add architecture diagrams
- [ ] Add sequence diagrams for pod lifecycle
- [ ] Document all annotations and their usage
- [ ] Add troubleshooting guide
- [ ] Add FAQ section
- [ ] Document cost optimization strategies

### Developer Experience
- [ ] Add development scripts (setup, test, run)
- [ ] Add VS Code debug configurations
- [ ] Add GitHub issue templates
- [ ] Add pull request templates
- [ ] Set up GitHub Discussions

### Advanced Features
- [ ] Support for DaemonSets (instance per node in source cluster)
- [ ] Support for StatefulSets with persistent volumes
- [ ] Add auto-scaling based on cluster load
- [ ] Implement pod autoscaling based on metrics
- [ ] Add support for GPU sharing (NVIDIA MIG)

## Technical Debt

### Code Quality
- [ ] Add godoc comments for all exported functions
- [ ] Add examples in godoc
- [ ] Reduce function complexity where possible
- [ ] Add more descriptive error messages with context

### Build & CI/CD
- [ ] Add golangci-lint configuration to CI
- [ ] Add integration test workflow
- [ ] Add coverage reporting to CI
- [ ] Add security scanning (gosec, snyk)
- [ ] Add dependency update automation (Dependabot/Renovate)
- [ ] Add release automation

### Refactoring Opportunities
- [ ] Consider splitting `internal/aws/client.go` into multiple files by concern
- [ ] Consider extracting tag building logic into separate package
- [ ] Review error handling patterns for consistency
- [ ] Consider adding a separate package for Kubernetes operations

## Completed ✅

- [x] Create project structure
- [x] Implement instance selector with 3-tier strategy
- [x] Add unit tests for instance selector (21 tests)
- [x] Create configuration management
- [x] Implement AWS EC2 client
- [x] Add unit tests for AWS client with mocks (30 tests)
- [x] Add integration tests with LocalStack (2 tests)
- [x] Integrate provider with AWS client
- [x] Create main application entry point
- [x] Add LocalStack support and configuration
- [x] Create LocalStack setup and test scripts
- [x] Add LocalStack documentation
- [x] Achieve Go Report Card A+ grade
- [x] Create GitHub repository

## Notes

- **Go Report Card:** Must maintain A+ grade - check cyclomatic complexity, gofmt, go vet
- **Testing Philosophy:** All tests must provide functional and development value
- **LocalStack First:** Always test with LocalStack before real AWS
- **Explicit > Template > Auto:** Instance selection priority order

---

Last updated: 2025-10-18
Maintainer: Scott Friedman (@scttfrdmn)
