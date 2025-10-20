# Project-Template-Go Implementation Guide for ORCA

## Overview

This guide provides step-by-step instructions to apply the project-template-go standards and conventions to the ORCA project.

**Reference Documents:**
- `PROJECT_TEMPLATE_GO_ANALYSIS.md` - Comprehensive analysis (15 sections, 837 lines)
- `TEMPLATE_QUICK_REFERENCE.md` - Quick lookup guide with tables and checklists

---

## Phase 1: Directory Structure Setup

### Step 1.1: Organize Existing Code

Move existing ORCA code to follow the template structure:

```bash
# Current state -> Target state
src/orca/
├── cmd/                        # NEW: Create this
│   ├── orca/                   # NEW: Main CLI application
│   │   └── main.go            # Move/create entry point
│   └── <other-services>/      # If multiple binaries
├── internal/                   # NEW: Move internal packages here
│   ├── config/                # Configuration management
│   ├── handler/               # HTTP/gRPC handlers
│   ├── service/               # Business logic
│   ├── model/                 # Data models
│   └── <domain>/              # Domain packages
├── pkg/                        # NEW: Public APIs (if applicable)
├── docs/                       # ENHANCE: Existing docs
├── .github/                    # ENHANCE: Existing GitHub config
│   ├── workflows/
│   ├── ISSUE_TEMPLATE/
│   └── scripts/
├── scripts/                    # NEW: Setup/utility scripts
├── tests/                      # NEW: Test organization
│   ├── unit/
│   ├── integration/
│   ├── e2e/
│   └── fixtures/
└── Makefile                    # NEW: Build automation
```

### Step 1.2: Create Missing Directories

```bash
mkdir -p cmd/orca
mkdir -p internal/{config,handler,service,model}
mkdir -p scripts
mkdir -p tests/{unit,integration,e2e,fixtures}
mkdir -p .github/ISSUE_TEMPLATE
mkdir -p .github/scripts
```

---

## Phase 2: Configuration Files

### Step 2.1: Create Makefile

Start with the template's Makefile and customize:

```makefile
# Copy from template and modify:
# - BINARY_NAME: "orca" (or service name)
# - Build commands for ORCA structure
# - Custom targets for ORCA-specific build steps
```

### Step 2.2: Set Up Pre-Commit

Copy and customize these files:

```bash
cp /path/to/template/.pre-commit-config.yaml .
cp /path/to/template/.golangci.yml .
cp /path/to/template/.editorconfig .
cp /path/to/template/.air.toml .
```

**Customization needed:**
- `.air.toml`: Update build commands if ORCA has multiple binaries
- `.golangci.yml`: Adjust linting rules for ORCA's code style
- `.editorconfig`: Keep as-is for cross-editor consistency

### Step 2.3: Update .gitignore

Copy template and add ORCA-specific entries:

```bash
cp /path/to/template/.gitignore .
# Add ORCA-specific patterns:
# - ORCA-generated files
# - Build artifacts specific to ORCA
# - Local development files
```

---

## Phase 3: GitHub Configuration

### Step 3.1: Setup GitHub Labels

Copy label configuration:

```bash
cp /path/to/template/.github/labels.yml .github/
```

Customize for ORCA:

```yaml
# Update area labels to match ORCA structure:
# Instead of:
# - area: cmd
# - area: pkg
# - area: internal

# Use ORCA-specific areas:
# - area: cli
# - area: api
# - area: core
# - area: orchestration
# - area: plugins
# etc.
```

Apply labels to repository:

```bash
gh label sync --repo <org>/orca --labels .github/labels.yml
```

### Step 3.2: Setup Issue Templates

Copy templates:

```bash
cp -r /path/to/template/.github/ISSUE_TEMPLATE/ .github/
```

Customize component dropdown in bug report to match ORCA structure:

```yaml
# bug_report.yml - Update component list:
options:
  - cli (command-line interface)
  - api (REST/gRPC API)
  - core (orchestration engine)
  - plugins (plugin system)
  - Build System / Makefile
  - Tests
  - Documentation
  - Dependencies
  - Other
```

### Step 3.3: Setup Pull Request Template

Copy template:

```bash
cp /path/to/template/.github/PULL_REQUEST_TEMPLATE.md .github/
```

No changes needed (use as-is, generic to all Go projects).

### Step 3.4: Setup Dependabot

Copy configuration:

```bash
cp /path/to/template/.github/dependabot.yml .github/
```

Customize if needed:

```yaml
# Update reviewer and schedule for ORCA workflow
updates:
  - package-ecosystem: "gomod"
    reviewers:
      - <orca-maintainer>  # Change from scttfrdmn
    schedule:
      interval: "weekly"
```

### Step 3.5: Setup Funding

Create/update:

```yaml
# .github/FUNDING.yml
ko_fi: <orca-funding>  # Or alternative funding sources
github_sponsors: <org>
```

---

## Phase 4: CI/CD Pipeline

### Step 4.1: Create GitHub Actions Workflow

Copy and customize:

```bash
cp /path/to/template/.github/workflows/ci.yml .github/workflows/
```

Customization points:

```yaml
# Adjust for ORCA:
name: ORCA CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    strategy:
      matrix:
        # Update Go versions if ORCA has different support
        go-version: [1.22.x, 1.23.x]
        
        # Keep or adjust OS matrix
        os: [ubuntu-latest, macos-latest, windows-latest]
```

### Step 4.2: Add Service-Specific Jobs

If ORCA has multiple services or special build requirements:

```yaml
# Add additional jobs:
build-docker:
  runs-on: ubuntu-latest
  needs: [test, lint]
  # Build and push ORCA Docker images

integration-tests:
  runs-on: ubuntu-latest
  # Integration tests specific to ORCA

performance-tests:
  runs-on: ubuntu-latest
  # Performance benchmarking if needed
```

---

## Phase 5: Documentation

### Step 5.1: Update README.md

Enhance existing README with template structure:

```markdown
# ORCA - [Full Description]

## Features
- [Feature list]

## Quick Start
- Installation
- Basic usage
- Environment variables

## Development
- Prerequisites (Go 1.23+, Make, Docker)
- Make targets
- Running tests

## API Documentation
- Endpoints
- Examples

## Contributing
- See CONTRIBUTING.md

## License
- MIT or appropriate license
```

### Step 5.2: Create CONTRIBUTING.md

Copy and customize:

```bash
cp /path/to/template/docs/CONTRIBUTING.md docs/
```

Customize for ORCA:

```markdown
# Contributing to ORCA

## Getting Started
# [ORCA-specific setup]

## Development Setup
# [ORCA-specific instructions]

## Code Style
# [ORCA conventions]

## [ORCA-specific sections]
```

### Step 5.3: Create CODE_OF_CONDUCT.md

Copy as-is (generic):

```bash
cp /path/to/template/docs/CODE_OF_CONDUCT.md docs/
```

### Step 5.4: Create SECURITY.md

Copy and customize:

```bash
cp /path/to/template/docs/SECURITY.md docs/
```

Update email and supported versions:

```markdown
# Security Policy

## Reporting a Vulnerability

Please report to: [orca-security@example.com]

## Supported Versions

| Version | Supported |
|---------|-----------|
| 2.x.x   | ✅        |
| 1.x.x   | ✅        |
| < 1.0   | ❌        |
```

---

## Phase 6: Setup Scripts

### Step 6.1: Create Install Hooks Script

Copy and customize:

```bash
cp /path/to/template/scripts/install-hooks.sh scripts/
```

No changes needed (generic Go setup).

### Step 6.2: Create Setup Script (Optional)

Create `scripts/setup-environment.sh` for ORCA-specific setup:

```bash
#!/bin/bash
# Setup ORCA development environment

set -e

echo "Setting up ORCA development environment..."

# Install Go if needed
# Install Docker if needed
# Setup databases if needed
# Run setup scripts

echo "ORCA environment ready!"
```

---

## Phase 7: Go Module Setup

### Step 7.1: Update go.mod

If moving code around:

```bash
# Update module name and imports
go mod edit -module github.com/orca-project/orca

# Update all imports in code
find . -name "*.go" -type f -exec sed -i '' \
  's|github.com/scttfrdmn/project-name|github.com/orca-project/orca|g' {} +

# Verify and tidy
go mod tidy
```

### Step 7.2: Create/Update go.sum

```bash
go mod download
go mod verify
```

---

## Phase 8: Code Quality Setup

### Step 8.1: Run Initial Linting

```bash
# Ensure all tools are installed
make install-tools

# Run linting
make lint

# Fix formatting
make fmt

# Resolve any issues
make vet
make staticcheck
make gosec
```

### Step 8.2: Setup Test Structure

Ensure tests follow template patterns:

```bash
# Unit tests co-located with code:
internal/config/config_test.go
internal/service/service_test.go

# Additional tests in tests/ directory:
tests/integration/api_test.go
tests/e2e/workflow_test.go

# Test fixtures:
tests/fixtures/sample-config.json
```

### Step 8.3: Configure Coverage

```bash
# Generate coverage report
make coverage

# Review coverage.html
# Target: >80% coverage on critical code paths
```

---

## Phase 9: Docker Setup

### Step 9.1: Create Dockerfile

Copy and customize:

```bash
cp /path/to/template/Dockerfile .
```

Customize:

```dockerfile
# If using different entry point:
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o orca \
    ./cmd/orca

# Update CMD if ORCA has different entry
ENTRYPOINT ["/orca"]
```

### Step 9.2: Create docker-compose.yml (Optional)

If ORCA needs local services:

```yaml
version: '3.8'
services:
  orca:
    build: .
    environment:
      - PORT=8080
      # ORCA-specific env vars
    ports:
      - "8080:8080"
    # Additional services (databases, cache, etc.)
```

---

## Phase 10: Verification and Testing

### Step 10.1: Run Full Quality Check

```bash
# Install all tools
./scripts/install-hooks.sh

# Run complete check
make all

# If successful: ✅ All checks pass
# If failures: Fix issues and re-run
```

### Step 10.2: Test Pre-Commit Hooks

```bash
# Setup hooks
pre-commit install

# Run hooks manually
pre-commit run --all-files

# If successful: Ready for git commits
```

### Step 10.3: Test CI/CD Locally (Optional)

```bash
# Install act (GitHub Actions local runner)
# Run workflows locally
act -j test
act -j lint
```

---

## Phase 11: Repository Configuration

### Step 11.1: Branch Protection Rules

In GitHub repository settings:

1. **Branch protection for main:**
   - Require PR reviews before merge
   - Dismiss stale PR approvals
   - Require status checks to pass:
     - test (all matrix combinations)
     - lint
     - build
     - docker

2. **Require CODEOWNERS** (optional)
   - Create `.github/CODEOWNERS`
   - Assign ownership to team members

### Step 11.2: Enable Features

- Enable "Automatically delete head branches"
- Enable "Require branches to be up to date before merging"
- Enable "Dismiss stale pull request approvals"

### Step 11.3: Setup Codecov (Optional)

```bash
# Add Codecov token as GitHub secret
# CI will automatically upload coverage
```

---

## Phase 12: Documentation Updates

### Step 12.1: Create ARCHITECTURE.md (Optional)

Document ORCA's architecture:

```markdown
# ORCA Architecture

## Overview
- System diagram
- Component breakdown

## Layers
- Handler layer (HTTP/gRPC)
- Service layer (business logic)
- Data layer
- Configuration

## Design Decisions
- Why this structure
- Key trade-offs
- Future scalability
```

### Step 12.2: Create API_DOCS.md (If Applicable)

```markdown
# ORCA API Documentation

## Endpoints
- [Detailed endpoint documentation]

## Authentication
- [Auth mechanisms]

## Examples
- [Usage examples]
```

---

## Customization Checklist

- [ ] **Module Path**: Updated `go.mod` with ORCA module path
- [ ] **GitHub Labels**: Customized area labels for ORCA structure
- [ ] **Environment Variables**: Added ORCA-specific env vars to README
- [ ] **Makefile**: Customized for ORCA build process
- [ ] **API Endpoints**: Updated documentation for ORCA APIs
- [ ] **Issue Templates**: Customized component dropdowns
- [ ] **CI/CD**: Adjusted Go versions and job matrix if needed
- [ ] **Docker**: Updated Dockerfile for ORCA entry point
- [ ] **Documentation**: All docs customized for ORCA
- [ ] **Funding**: Updated FUNDING.yml with ORCA sponsor
- [ ] **CODEOWNERS**: Created (optional) with team assignments
- [ ] **Branch Protection**: Configured in GitHub

---

## Post-Implementation Validation

### Checklist

1. **Run full quality suite:**
   ```bash
   make all
   ```

2. **Verify pre-commit hooks:**
   ```bash
   git add .
   git commit -m "initial: apply project template structure"
   ```

3. **Test CI/CD pipeline:**
   ```bash
   git push origin main
   # Check GitHub Actions
   ```

4. **Test Docker build:**
   ```bash
   make docker-build
   docker run --rm orca:latest --help
   ```

5. **Verify documentation:**
   - README.md readable and complete
   - CONTRIBUTING.md clear
   - API docs (if applicable) functional

6. **Test label sync:**
   ```bash
   gh label sync --repo <org>/orca --labels .github/labels.yml
   ```

7. **Create test PR:**
   - Create feature branch
   - Make dummy change
   - Verify templates work
   - Verify checks run

---

## Common Issues and Solutions

### Issue: Pre-commit hooks fail

**Solution:**
```bash
# Verify tools installed
which golangci-lint staticcheck gosec

# Install if missing
make install-tools

# Run manually
pre-commit run --all-files --show-diff
```

### Issue: CI/CD matrix failures

**Solution:**
- Check Go version compatibility
- Verify cross-platform code
- Check for platform-specific issues in tests

### Issue: Docker build fails

**Solution:**
- Verify Dockerfile entry point matches binary name
- Check binary actually exists after build
- Verify multi-stage build paths

### Issue: Label sync fails

**Solution:**
```bash
# Verify GitHub CLI installed
gh --version

# Authenticate
gh auth login

# Try again
gh label sync --repo <org>/orca --labels .github/labels.yml
```

---

## Timeline Estimate

- **Phase 1-2** (Directory + Config): 1-2 hours
- **Phase 3-4** (GitHub + CI/CD): 1-2 hours
- **Phase 5-6** (Documentation + Scripts): 1-2 hours
- **Phase 7-9** (Go Module + Docker): 1-2 hours
- **Phase 10-12** (Verification + Polish): 1-2 hours

**Total: 5-10 hours** for complete implementation

---

## Support and References

- **Template Repository**: `/Users/scttfrdmn/src/project-template-go`
- **Analysis Document**: `PROJECT_TEMPLATE_GO_ANALYSIS.md`
- **Quick Reference**: `TEMPLATE_QUICK_REFERENCE.md`

For specific implementations, refer to the template files directly:
- Makefile patterns
- Workflow configurations
- Linting rules
- Documentation examples

---

## Next Steps After Implementation

1. **Monitor first PR**: Ensure all checks work
2. **Gather feedback**: Ask team about workflow changes
3. **Iterate on labels**: Refine based on actual usage
4. **Document patterns**: Create ORCA-specific conventions
5. **Training**: Share with team about new structure
6. **Gradual migration**: Move existing code incrementally if needed
7. **Continuous improvement**: Update as ORCA evolves

