# ORCA Testing Guide

This document describes ORCA's testing approach: **pragmatic, functional, and value-driven**. Every test serves a clear purpose - catching real bugs, preventing regressions, or improving development speed.

## Testing Philosophy

**No testing for testing's sake.** Every test must provide:
1. **Bug prevention** - Catches real issues before production
2. **Regression protection** - Prevents broken code from merging
3. **Development velocity** - Makes refactoring safe and fast
4. **Documentation** - Shows how code should work

If a test doesn't provide clear value, we don't write it.

---

## Test Types

### 1. Lint (Code Quality)
**Purpose**: Catch common mistakes and enforce consistency
**When**: Pre-commit, CI on every PR
**Value**: Prevents bugs, ensures idiomatic Go code

```bash
# Run linters
make lint

# Auto-fix issues
golangci-lint run --fix
```

**What we check**:
- ✅ `go vet` - Suspicious constructs
- ✅ `staticcheck` - Go best practices
- ✅ `errcheck` - Unchecked errors
- ✅ `gosec` - Security issues
- ✅ `gofmt` - Code formatting
- ✅ `gocritic` - Style issues

**Example catches**:
```go
// BAD: Uncaught error
ec2Client.TerminateInstance(instanceID)

// GOOD: Error handling
if err := ec2Client.TerminateInstance(instanceID); err != nil {
    return fmt.Errorf("failed to terminate: %w", err)
}
```

---

### 2. Unit Tests (Fast Feedback)
**Purpose**: Test individual functions in isolation
**When**: TDD during development, CI on every PR
**Value**: Fast feedback, safe refactoring, clear API contracts

```bash
# Run unit tests
make test

# Run specific package
go test -v ./pkg/instances/...

# Run specific test
go test -v ./pkg/instances -run TestExplicitSelector
```

**What to unit test**:
- ✅ **Instance selection logic** - Core functionality
- ✅ **Configuration parsing** - Prevents config bugs
- ✅ **Pod annotation extraction** - Common bug source
- ✅ **Budget calculations** - Critical for cost control
- ✅ **Tag generation** - Required for pod tracking
- ✅ **Error handling** - Ensures graceful failures

**What NOT to unit test**:
- ❌ Trivial getters/setters
- ❌ Third-party library behavior
- ❌ Code that's just wiring

#### Example: Instance Selector

```go
// pkg/instances/selector_test.go

func TestExplicitSelector(t *testing.T) {
    tests := []struct {
        name        string
        pod         *corev1.Pod
        expected    string
        expectError bool
    }{
        {
            name: "explicit p5.48xlarge annotation",
            pod: &corev1.Pod{
                ObjectMeta: metav1.ObjectMeta{
                    Annotations: map[string]string{
                        AnnotationInstanceType: "p5.48xlarge",
                    },
                },
            },
            expected:    "p5.48xlarge",
            expectError: false,
        },
        {
            name: "missing annotation returns error",
            pod: &corev1.Pod{
                ObjectMeta: metav1.ObjectMeta{
                    Annotations: map[string]string{},
                },
            },
            expected:    "",
            expectError: true,
        },
        {
            name: "invalid instance type returns error",
            pod: &corev1.Pod{
                ObjectMeta: metav1.ObjectMeta{
                    Annotations: map[string]string{
                        AnnotationInstanceType: "invalid-type",
                    },
                },
            },
            expected:    "",
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            selector := NewExplicitSelector()
            result, err := selector.Select(tt.pod)

            if tt.expectError && err == nil {
                t.Error("expected error but got nil")
            }
            if !tt.expectError && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

**Value**: This test catches 90% of selector bugs in milliseconds.

#### Example: Budget Enforcement

```go
// pkg/budget/enforcer_test.go

func TestBudgetEnforcement(t *testing.T) {
    tests := []struct {
        name           string
        dailyBudget    float64
        currentSpend   float64
        instanceCost   float64
        shouldAllow    bool
    }{
        {
            name:         "under budget allows instance",
            dailyBudget:  50.00,
            currentSpend: 10.00,
            instanceCost: 5.00,
            shouldAllow:  true,
        },
        {
            name:         "at budget limit denies instance",
            dailyBudget:  50.00,
            currentSpend: 48.00,
            instanceCost: 5.00,
            shouldAllow:  false,
        },
        {
            name:         "zero budget denies all",
            dailyBudget:  0.00,
            currentSpend: 0.00,
            instanceCost: 0.01,
            shouldAllow:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            enforcer := &BudgetEnforcer{
                DailyLimit:   tt.dailyBudget,
                CurrentSpend: tt.currentSpend,
            }

            allowed := enforcer.CanLaunchInstance(tt.instanceCost)

            if allowed != tt.shouldAllow {
                t.Errorf("expected %v, got %v", tt.shouldAllow, allowed)
            }
        })
    }
}
```

**Value**: Prevents cost overruns - critical for production.

---

### 3. Integration Tests (Real Interactions)
**Purpose**: Test components working together
**When**: After unit tests pass, CI on merge to main
**Value**: Catches integration bugs, validates AWS interactions

```bash
# Run integration tests (requires AWS credentials)
make integration-test

# Or run with specific tag
go test -v -tags=integration ./...
```

**What to integration test**:
- ✅ **AWS SDK calls** - Ensure EC2 APIs work
- ✅ **Kubernetes API** - Verify pod operations
- ✅ **Instance lifecycle** - Create → Run → Terminate
- ✅ **Configuration loading** - Test full config chain
- ✅ **Error scenarios** - Network failures, AWS throttling

#### Example: AWS Client Integration

```go
// internal/aws/client_test.go
// +build integration

func TestEC2CreateInstance(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Setup
    cfg := loadTestConfig(t)
    client := NewClient(cfg)

    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "test-pod",
            Namespace: "default",
            Annotations: map[string]string{
                AnnotationInstanceType: "t3.small",
            },
        },
        Spec: corev1.PodSpec{
            Containers: []corev1.Container{
                {
                    Name:  "test",
                    Image: "busybox",
                },
            },
        },
    }

    // Execute
    ctx := context.Background()
    instanceID, err := client.CreateInstance(ctx, pod)
    if err != nil {
        t.Fatalf("failed to create instance: %v", err)
    }

    // Verify
    instance, err := client.DescribeInstance(ctx, instanceID)
    if err != nil {
        t.Fatalf("failed to describe instance: %v", err)
    }
    if instance.State != "running" && instance.State != "pending" {
        t.Errorf("expected running/pending, got %s", instance.State)
    }

    // Cleanup
    defer func() {
        if err := client.TerminateInstance(ctx, instanceID); err != nil {
            t.Errorf("cleanup failed: %v", err)
        }
    }()
}
```

**Value**: Catches AWS API changes, permission issues, network problems.

**Cost control for integration tests**:
```go
// Use smallest instance types
const testInstanceType = "t3.nano"  // $0.0052/hour

// Set short timeouts
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

// Always cleanup
defer cleanupTestResources(t, client)
```

---

### 4. Smoke Tests (Quick Sanity Check)
**Purpose**: Verify basic functionality works end-to-end
**When**: After deployment, before releasing
**Value**: Catches deployment issues, validates basic workflows

```bash
# Run smoke tests against deployed ORCA
make smoke-test CLUSTER=kind-orca-dev

# Or manually
./scripts/smoke-test.sh
```

**What smoke tests check**:
- ✅ ORCA pod is running
- ✅ Virtual node is registered
- ✅ Simple pod can be created
- ✅ Instance launches successfully
- ✅ Pod reaches Running state
- ✅ Pod can be deleted
- ✅ Instance terminates

#### Example: Smoke Test Script

```bash
#!/bin/bash
# scripts/smoke-test.sh

set -e

echo "🔍 Running ORCA smoke tests..."

# 1. Check ORCA is running
echo "Checking ORCA pod..."
kubectl get pods -n kube-system -l app=orca | grep Running

# 2. Check virtual node exists
echo "Checking virtual node..."
kubectl get node -l orca.research/provider=aws

# 3. Deploy test pod
echo "Deploying test pod..."
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: smoke-test-pod
  annotations:
    orca.research/instance-type: "t3.small"
spec:
  nodeSelector:
    orca.research/provider: "aws"
  containers:
  - name: test
    image: busybox
    command: ["sh", "-c", "echo 'Smoke test passed' && sleep 60"]
EOF

# 4. Wait for pod to run
echo "Waiting for pod to run..."
kubectl wait --for=condition=Ready pod/smoke-test-pod --timeout=5m

# 5. Verify instance exists
echo "Checking EC2 instance..."
aws ec2 describe-instances \
  --filters "Name=tag:orca.research/pod,Values=default/smoke-test-pod" \
  --query 'Reservations[0].Instances[0].State.Name' \
  --output text | grep running

# 6. Cleanup
echo "Cleaning up..."
kubectl delete pod smoke-test-pod --wait=true

# 7. Verify instance terminated
echo "Verifying cleanup..."
sleep 30
aws ec2 describe-instances \
  --filters "Name=tag:orca.research/pod,Values=default/smoke-test-pod" \
  --query 'Reservations[0].Instances[0].State.Name' \
  --output text | grep -E 'terminated|shutting-down'

echo "✅ Smoke tests passed!"
```

**Value**: 5-minute test catches 80% of deployment issues.

---

### 5. Regression Tests (Prevent Known Bugs)
**Purpose**: Ensure fixed bugs stay fixed
**When**: CI on every PR, before release
**Value**: Prevents bugs from reappearing

**Process**:
1. Bug is reported
2. Write failing test that reproduces bug
3. Fix bug
4. Test now passes
5. Test prevents regression forever

#### Example: Regression for Issue #42

```go
// pkg/provider/provider_test.go

// TestIssue42_PodWithNoAnnotations tests the fix for:
// https://github.com/scttfrdmn/orca/issues/42
// Bug: ORCA crashed when pod had no annotations
func TestIssue42_PodWithNoAnnotations(t *testing.T) {
    provider := NewProvider(testConfig())

    pod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:        "no-annotations",
            Namespace:   "default",
            Annotations: nil,  // This was causing panic
        },
        Spec: corev1.PodSpec{
            Containers: []corev1.Container{{Name: "test", Image: "busybox"}},
        },
    }

    // Should not panic
    err := provider.CreatePod(context.Background(), pod)

    // Should return error, not panic
    if err == nil {
        t.Error("expected error for missing annotations, got nil")
    }
    if !strings.Contains(err.Error(), "missing required annotation") {
        t.Errorf("expected helpful error message, got: %v", err)
    }
}
```

**Value**: Bug #42 can never come back.

---

## Test Coverage

We track coverage but **don't worship it**. 80%+ coverage is good. 100% is wasteful.

```bash
# Generate coverage report
make coverage

# View in browser
open coverage.html

# Fail CI if coverage drops below 80%
go test -coverprofile=coverage.txt ./...
go tool cover -func=coverage.txt | grep total | awk '{print $3}' | sed 's/%//' | awk '$1 < 80 {exit 1}'
```

**Coverage priorities**:
1. **Critical paths**: 100% (budget enforcement, instance selection)
2. **Common paths**: 80%+ (pod creation, deletion)
3. **Error handling**: 70%+ (various failure modes)
4. **Happy paths**: 60%+ (basic workflows)

**Don't cover**:
- Generated code
- Third-party integrations (test with mocks)
- Trivial code (simple getters)

---

## Testing Tools

### Mocking AWS SDK

```go
// internal/aws/mock.go

type MockEC2Client struct {
    CreateInstanceFunc    func(context.Context, *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error)
    TerminateInstanceFunc func(context.Context, *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error)
}

func (m *MockEC2Client) RunInstances(ctx context.Context, input *ec2.RunInstancesInput, opts ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error) {
    if m.CreateInstanceFunc != nil {
        return m.CreateInstanceFunc(ctx, input)
    }
    return &ec2.RunInstancesOutput{
        Instances: []types.Instance{
            {InstanceId: aws.String("i-mock123")},
        },
    }, nil
}
```

### Test Fixtures

```go
// pkg/testing/fixtures.go

// CreateTestPod creates a pod for testing
func CreateTestPod(name, instanceType string) *corev1.Pod {
    return &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: "default",
            Annotations: map[string]string{
                AnnotationInstanceType: instanceType,
            },
        },
        Spec: corev1.PodSpec{
            Containers: []corev1.Container{
                {Name: "test", Image: "busybox"},
            },
        },
    }
}
```

### Table-Driven Tests

```go
func TestInstanceSelection(t *testing.T) {
    tests := []struct {
        name     string
        input    *corev1.Pod
        expected string
        wantErr  bool
    }{
        // Test cases here
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic here
        })
    }
}
```

---

## CI/CD Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/test.yml

name: Tests

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  unit-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: make test
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: ./coverage.txt

  integration-test:
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Configure AWS
        uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-west-2
      - name: Run integration tests
        run: make integration-test

  smoke-test:
    runs-on: ubuntu-latest
    needs: [unit-test]
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - name: Create kind cluster
        uses: helm/kind-action@v1
      - name: Build and deploy ORCA
        run: |
          make docker-build
          kind load docker-image orca:latest
          kubectl apply -f deploy/kubernetes/
      - name: Run smoke tests
        run: ./scripts/smoke-test.sh
```

---

## Test Organization

```
orca/
├── pkg/
│   ├── provider/
│   │   ├── provider.go
│   │   ├── provider_test.go      # Unit tests
│   │   └── provider_integration_test.go  # Integration tests
│   ├── instances/
│   │   ├── selector.go
│   │   └── selector_test.go
│   └── budget/
│       ├── enforcer.go
│       └── enforcer_test.go
├── internal/
│   └── aws/
│       ├── client.go
│       ├── client_test.go
│       ├── client_integration_test.go
│       └── mock.go               # Mock implementations
├── scripts/
│   ├── smoke-test.sh             # Smoke tests
│   └── regression-test.sh        # Regression suite
└── tests/
    ├── fixtures/                 # Test fixtures
    └── e2e/                      # End-to-end tests
```

---

## Testing Checklist

Before merging a PR:

- [ ] `make lint` passes
- [ ] `make test` passes with >80% coverage
- [ ] Integration tests pass (if touching AWS code)
- [ ] Smoke test passes (if changing core logic)
- [ ] Added regression test (if fixing bug)
- [ ] Updated test documentation (if adding new test patterns)

Before releasing:

- [ ] All CI tests pass
- [ ] Smoke tests pass on kind
- [ ] Smoke tests pass on EKS (if available)
- [ ] Manual testing of new features
- [ ] Regression suite passes
- [ ] Performance tests pass (if applicable)

---

## Performance Testing

For features impacting performance:

```go
func BenchmarkInstanceSelection(b *testing.B) {
    selector := NewExplicitSelector()
    pod := CreateTestPod("test", "p5.48xlarge")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := selector.Select(pod)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

Run with:
```bash
go test -bench=. -benchmem ./pkg/instances/
```

---

## Summary

Our testing approach:

1. **Lint**: Catch mistakes early (seconds)
2. **Unit tests**: Fast feedback (milliseconds)
3. **Integration tests**: Validate AWS (minutes, limited runs)
4. **Smoke tests**: Deployment validation (5 minutes)
5. **Regression tests**: Keep bugs fixed (continuous)

Every test provides clear value. No busywork. Focus on preventing real bugs and enabling safe refactoring.

**Test pyramid**:
```
        /\
       /  \      E2E Tests (few, slow, expensive)
      /    \
     /------\    Integration Tests (some, moderate)
    /        \
   /----------\  Unit Tests (many, fast, cheap)
  /------------\
 /   Linting   \ Static Analysis (everywhere, instant)
/________________\
```

Most value comes from the bottom. Write more unit tests, fewer integration tests, minimal E2E tests.
