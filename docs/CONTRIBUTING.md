# Contributing to ORCA

Thank you for your interest in contributing to ORCA! We welcome contributions from the community.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker (for containerized testing)
- kubectl (for Kubernetes testing)
- AWS account (for integration testing)
- golangci-lint (for linting)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/scttfrdmn/orca.git
   cd orca
   ```

2. **Install dependencies**
   ```bash
   make mod-download
   ```

3. **Install development tools**
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

4. **Build the project**
   ```bash
   make build
   ```

5. **Run tests**
   ```bash
   make test
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write idiomatic Go code
   - Follow the project structure conventions
   - Add tests for new functionality
   - Update documentation as needed

3. **Format and lint**
   ```bash
   make fmt
   make lint
   make vet
   ```

4. **Run tests**
   ```bash
   make test
   make coverage
   ```

5. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Adding or updating tests
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `chore:` - Maintenance tasks

Examples:
```
feat: add explicit instance selection support
fix: handle pod deletion edge cases
docs: update instance selection guide
test: add unit tests for AWS client
```

### Pull Request Process

1. **Update documentation** - Ensure README.md and relevant docs are updated
2. **Add changelog entry** - Add your changes to CHANGELOG.md under [Unreleased]
3. **Ensure tests pass** - All tests must pass (`make test`)
4. **Ensure linting passes** - Code must pass linting (`make lint`)
5. **Create pull request** - Provide clear description of changes
6. **Request review** - Tag maintainers for review

## Code Style Guidelines

### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting (run `make fmt`)
- Keep functions focused and small
- Write self-documenting code with clear variable names
- Add comments for exported functions and complex logic

### Error Handling

```go
// Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create instance: %w", err)
}

// Bad: Generic error
if err != nil {
    return err
}
```

### Logging

```go
// Use structured logging
log.Info("creating pod",
    "namespace", pod.Namespace,
    "name", pod.Name,
    "instanceType", instanceType)
```

## Testing Guidelines

### Unit Tests

- Write tests for all new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies (AWS SDK, Kubernetes API)

Example:
```go
func TestSelectInstanceType(t *testing.T) {
    tests := []struct {
        name     string
        pod      *corev1.Pod
        expected string
    }{
        {
            name: "explicit instance type",
            pod: &corev1.Pod{
                ObjectMeta: metav1.ObjectMeta{
                    Annotations: map[string]string{
                        "orca.research/instance-type": "p5.48xlarge",
                    },
                },
            },
            expected: "p5.48xlarge",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := selectInstanceType(tt.pod)
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

### Integration Tests

- Test against real AWS services when possible
- Use AWS localstack for local testing
- Clean up resources after tests

## Documentation

### Code Documentation

- Document all exported functions, types, and constants
- Use godoc format

```go
// CreatePod creates a new pod by launching an EC2 instance.
// It returns an error if the instance cannot be created.
func (p *Provider) CreatePod(ctx context.Context, pod *corev1.Pod) error {
    // Implementation
}
```

### User Documentation

- Update README.md for user-facing changes
- Add examples in `examples/` directory
- Update relevant docs in `docs/` directory

## Architecture Guidelines

### Project Structure

```
cmd/orca/           # Main application entry point
pkg/                # Public libraries
  provider/         # Virtual Kubelet provider implementation
  config/           # Configuration management
  instances/        # Instance selection logic
internal/           # Private application code
  aws/              # AWS SDK integration
  container/        # Container runtime integration
  metrics/          # Metrics and monitoring
```

### Design Principles

1. **Explicit over Implicit** - Users should specify what they want
2. **Research-First** - Optimize for research computing workflows
3. **Production-Grade** - Write code as if it will run at scale
4. **Testability** - Design for testability from the start
5. **Observability** - Include metrics, logging, tracing

## Getting Help

- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - Questions and general discussion
- **Research Partners** - Reach out to NRP, SDSU contacts

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Welcome newcomers
- Assume good intentions

## License

By contributing to ORCA, you agree that your contributions will be licensed under the Apache License 2.0.
