# Project-Template-Go Quick Reference for ORCA

## Key Files Location Reference

| Component | File/Directory | Purpose |
|-----------|---|---|
| **Build** | `Makefile` | All build automation and common tasks |
| **Config** | `.editorconfig` | Cross-editor code formatting |
| **Config** | `.gitignore` | Version control exclusions |
| **Config** | `.air.toml` | Hot-reload development configuration |
| **Container** | `Dockerfile` | Multi-stage Docker build (Alpine -> scratch) |
| **Quality Assurance** | `.pre-commit-config.yaml` | Git hooks for automated checks |
| **Quality Assurance** | `.golangci.yml` | 28-linter configuration |
| **Testing** | `tests/` | Test organization structure |
| **CI/CD** | `.github/workflows/ci.yml` | GitHub Actions pipeline |
| **GitHub** | `.github/labels.yml` | 22 issue labels (types, priorities, areas) |
| **GitHub** | `.github/ISSUE_TEMPLATE/` | Bug and feature request forms |
| **GitHub** | `.github/PULL_REQUEST_TEMPLATE.md` | PR submission form |
| **GitHub** | `.github/dependabot.yml` | Automated dependency updates |
| **Documentation** | `README.md` | Main project overview |
| **Documentation** | `docs/CONTRIBUTING.md` | Contributor guidelines |
| **Documentation** | `docs/CODE_OF_CONDUCT.md` | Community standards |
| **Documentation** | `docs/SECURITY.md` | Security policy |
| **Code** | `cmd/app/main.go` | Application entry point |
| **Code** | `internal/config/` | Configuration management |
| **Code** | `internal/service/` | Business logic layer |
| **Code** | `internal/handler/` | HTTP handlers/routing |
| **Scripts** | `scripts/install-hooks.sh` | Initial setup automation |

---

## Essential Make Targets

```bash
# Development
make dev                # Hot-reload development
make run               # Build and run
make build             # Compile binary

# Quality (Run before commit)
make pre-commit        # Run all checks
make fmt               # Format code
make lint              # Run linting
make vet               # Run go vet
make staticcheck       # Static analysis
make gosec             # Security scan

# Testing
make test              # Run all tests
make test-race         # With race detector
make coverage          # Coverage report

# Maintenance
make clean             # Remove artifacts
make deps              # Update dependencies
make install-tools    # Install dev tools
make docker-build      # Build container
```

---

## GitHub Configuration Structure

### Labels (by category)

**Types (5):** bug, enhancement, documentation, question, technical-debt

**Priority (3):** priority:high, priority:medium, priority:low

**Areas (7):** area:cmd, area:pkg, area:internal, area:build, area:tests, area:docs, area:deps

**Status (4):** triage, ready, in-progress, blocked

**Resolution (3):** duplicate, wontfix, invalid

**Special (6):** good first issue, help wanted, breaking-change, security, performance, dependencies

### Issue Templates

1. **Bug Report** - Component, steps to reproduce, expected/actual behavior, environment info
2. **Feature Request** - Problem statement, proposed solution, examples, priority

### Pull Request Sections

- Description, related issues, type of change, user impact, testing, checklist, breaking changes

---

## Pre-Commit Hooks (Auto-run on Commit)

| Hook | Action | Files |
|------|--------|-------|
| go-fmt | Format code | *.go |
| go-mod-tidy | Tidy dependencies | go.mod/go.sum |
| golangci-lint | Lint with fixes | *.go |
| go-vet | Check suspicious code | *.go |
| staticcheck | Static analysis | *.go |
| gosec | Security scan | *.go |
| go-test | Run tests | *.go |
| Generic | Trailing whitespace, EOF, YAML, merge conflicts | All |

---

## Directory Structure Pattern

```
project-template-go/
├── cmd/app/              ← Executable entry point
├── internal/             ← Private implementation
│   ├── config/           ← Configuration management
│   ├── handler/          ← HTTP/API handlers
│   └── service/          ← Business logic
├── pkg/                  ← Public packages (optional)
├── tests/                ← Test organization (fixtures, etc.)
├── .github/
│   ├── workflows/        ← CI/CD pipelines
│   ├── ISSUE_TEMPLATE/   ← Issue forms
│   ├── scripts/          ← GitHub automation
│   ├── labels.yml        ← Label definitions
│   └── PULL_REQUEST_TEMPLATE.md
├── scripts/              ← Dev/setup scripts
├── docs/                 ← Documentation
├── Dockerfile            ← Container config
├── Makefile              ← Build automation
├── .pre-commit-config.yaml ← Git hooks
├── .golangci.yml         ← Linting rules
├── .air.toml             ← Hot-reload config
├── .editorconfig         ← Editor formatting
├── .gitignore            ← Git exclusions
└── go.mod / go.sum       ← Dependencies
```

---

## Critical Configuration Settings

### Dockerfile
- Base: `golang:1.23-alpine` → `scratch`
- Security: Non-root user, static linking
- Health check: `/app health`

### CI/CD Matrix
- Go versions: 1.22.x, 1.23.x
- OS: ubuntu-latest, macos-latest, windows-latest
- Separate jobs: test, lint, build, docker

### Makefile Build Info
```makefile
VERSION = git describe --tags --always --dirty
BUILD_TIME = date -u +"%Y-%m-%dT%H:%M:%SZ"
GIT_COMMIT = git rev-parse --short HEAD
```
Injected via `-ldflags` at build time

### Go Version
- Minimum: 1.23
- Uses modern Go features
- Module-based dependency management

---

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| PORT | 8080 | Server port |
| HOST | 0.0.0.0 | Server host |
| ENVIRONMENT | development | Environment name |
| LOG_LEVEL | info | Logging level |
| DEBUG | false | Debug mode flag |

---

## API Endpoints (Example)

```
GET  /health          → Health check
GET  /info           → Service information
POST /api/process    → Data processing
```

---

## Testing Patterns

**Table-driven tests:**
```go
tests := []struct {
    name        string
    input       string
    expected    string
    expectError bool
}{...}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {...})
}
```

**Coverage:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Key Implementation Details

### Three-Layer Architecture
1. **Handler** - HTTP/API layer
2. **Service** - Business logic
3. **Config** - Configuration management

### Graceful Shutdown
- Context-based control
- Signal handling (SIGTERM, SIGINT)
- 30-second timeout for cleanup

### Error Handling
- Context-aware execution
- Structured error returns
- HTTP error code mapping

---

## Tools Required

| Tool | Purpose | Installation |
|------|---------|--------------|
| go | Language | Native |
| pre-commit | Git hooks | pip/brew |
| golangci-lint | Linting | go install |
| staticcheck | Analysis | go install |
| gosec | Security | go install |
| air | Hot-reload | go install |
| gh CLI | GitHub management | brew/apt |
| Docker | Containers | brew/apt |

---

## Setup Checklist

1. Clone repository
2. Run `./scripts/install-hooks.sh`
3. Run `make deps` to install dependencies
4. Run `make all` to verify setup
5. Run `gh label sync --labels .github/labels.yml` to setup labels
6. Update `go.mod` with project module path
7. Customize `.github/labels.yml` for your areas
8. Update FUNDING.yml with your sponsor

---

## Common Workflows

### Adding a New Feature
```bash
git checkout -b feature/description
make dev              # Hot-reload development
make test             # Run tests
make pre-commit       # Final quality check
git commit -m "feat: description"
git push origin feature/description
# Create PR using template
```

### Making a Release
```bash
# Update CHANGELOG.md with new version
make clean
make all              # Full validation
git tag v1.2.3
git push origin main --tags
# GitHub Actions automatically builds and creates release
```

### Updating Dependencies
```bash
make deps             # go mod download && go mod tidy
make test             # Verify tests pass
git commit -m "deps: update dependencies"
```

---

## Performance Considerations

- **Docker:** Multi-stage build reduces image size
- **Testing:** Race detector enabled in CI
- **Caching:** Go module cache in GitHub Actions
- **Linting:** Parallelized linter with 5-minute timeout

---

## Security Practices

1. Non-root container execution
2. Static binary linking (no external deps)
3. Security scanning with gosec
4. Dependency verification
5. Vulnerability reporting policy (email, not public)
6. Code review via PR workflow

---

## Customization for ORCA

**Immediate changes needed:**
- Module path in `go.mod`
- `scttfrdmn` → your organization
- Area labels in `.github/labels.yml`
- Environment variables for ORCA
- API endpoints for ORCA services
- FUNDING.yml sponsor link
- Go version support matrix

**Optional enhancements:**
- CodeQL security scanning
- Branch protection rules
- CODEOWNERS file
- GitHub Pages documentation
- Release automation
