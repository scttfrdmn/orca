# ORCA Quick Start Guide

This guide will help you get ORCA up and running in your Kubernetes cluster.

## Prerequisites

- Kubernetes cluster (1.28+)
- AWS account with appropriate permissions
- kubectl configured to access your cluster
- Go 1.21+ (if building from source)

## Installation Steps

### 1. Build ORCA

```bash
# Clone the repository
git clone https://github.com/scttfrdmn/orca.git
cd orca

# Build the binary
go build -o orca ./cmd/orca

# Or build with version info
VERSION=$(git describe --tags --always --dirty)
GIT_COMMIT=$(git rev-parse HEAD)
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build \
  -ldflags="-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT} -X main.buildDate=${BUILD_DATE}" \
  -o orca \
  ./cmd/orca
```

### 2. Configure AWS Credentials

ORCA needs AWS credentials to create EC2 instances. Three options:

**Option A: AWS Profile (Development)**
```bash
export AWS_PROFILE=orca
./orca --config config.yaml --kubeconfig ~/.kube/config
```

**Option B: Environment Variables**
```bash
export AWS_ACCESS_KEY_ID=AKIA...
export AWS_SECRET_ACCESS_KEY=...
export AWS_REGION=us-west-2
```

**Option C: IRSA (Production - Recommended)**
Create IAM role and service account (see deploy/README.md)

### 3. Create Configuration File

Create `config.yaml`:

```yaml
aws:
  region: us-west-2
  vpcID: vpc-xxxxx
  subnetID: subnet-xxxxx
  securityGroupIDs:
    - sg-xxxxx
  tags:
    Environment: production
    Project: orca

node:
  name: orca-aws-node
  labels:
    orca.research/provider: "aws"
    orca.research/region: "us-west-2"
  taints:
    - key: orca.research/burst-node
      value: "true"
      effect: NoSchedule
  cpu: "1000"
  memory: "4Ti"
  pods: "1000"
  gpu: "100"

instances:
  selectionMode: explicit
  defaultLaunchType: on-demand

logging:
  level: info
  format: json

metrics:
  enabled: true
  port: 8080
  path: /metrics
```

### 4. Run ORCA Locally (Testing)

```bash
# Start ORCA
./orca \
  --config config.yaml \
  --kubeconfig ~/.kube/config \
  --namespace kube-system \
  --log-level debug

# You should see:
# {"level":"info","time":"...","message":"Starting ORCA","version":"..."}
# {"level":"info","message":"Starting HTTP server","port":8080}
# {"level":"info","message":"Starting ORCA Virtual Kubelet node"}
# {"level":"info","message":"ORCA is running. Press Ctrl+C to stop.","http_port":8080}
```

### 5. Verify Node Registration

In another terminal:

```bash
# Check that orca-aws-node appears
kubectl get nodes

# Should show:
# NAME              STATUS   ROLES    AGE   VERSION
# orca-aws-node     Ready    <none>   10s   v1.0.0-orca
# ...

# Check node details
kubectl describe node orca-aws-node
```

### 6. Deploy a Test Pod

Create `test-pod.yaml`:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-burst
  annotations:
    orca.research/instance-type: "t3.small"
spec:
  nodeSelector:
    orca.research/provider: "aws"
  tolerations:
    - key: orca.research/burst-node
      operator: Equal
      value: "true"
      effect: NoSchedule
  containers:
    - name: test
      image: nginx:latest
      ports:
        - containerPort: 80
```

Deploy and watch:

```bash
# Deploy pod
kubectl apply -f test-pod.yaml

# Watch pod status
kubectl get pods -w

# You should see:
# test-burst   0/1     Pending   0          1s
# test-burst   0/1     Pending   0          5s
# test-burst   0/1     Running   0          65s  # After EC2 instance starts
```

### 7. Verify EC2 Instance Created

```bash
# List ORCA-managed instances
aws ec2 describe-instances \
  --filters "Name=tag:ManagedBy,Values=ORCA" \
  --query 'Reservations[*].Instances[*].[InstanceId,InstanceType,State.Name,Tags[?Key==`Name`].Value|[0]]' \
  --output table

# Should show:
# |  i-0123456789abcdef  |  t3.small  |  running  |  orca-default-test-burst  |
```

### 8. Test Health Checks

```bash
# Liveness check
curl http://localhost:8080/healthz
# {"status":"ok","service":"orca"}

# Readiness check  
curl http://localhost:8080/readyz
# {"status":"ready","service":"orca"}

# Prometheus metrics
curl http://localhost:8080/metrics
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
# go_goroutines 42
# ...
```

### 9. Clean Up

```bash
# Delete the test pod
kubectl delete pod test-burst

# The EC2 instance will be automatically terminated

# Stop ORCA
# Press Ctrl+C in the ORCA terminal

# Verify instance terminated
aws ec2 describe-instances \
  --filters "Name=tag:ManagedBy,Values=ORCA" \
  --query 'Reservations[*].Instances[*].State.Name'
```

## Production Deployment

For production deployment as a Kubernetes Deployment:

### 1. Create Namespace

```bash
kubectl create namespace orca-system
```

### 2. Create RBAC Resources

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: orca
  namespace: orca-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: orca-role
rules:
  - apiGroups: [""]
    resources: ["nodes", "nodes/status"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: [""]
    resources: ["pods", "pods/status"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: ["coordination.k8s.io"]
    resources: ["leases"]
    verbs: ["get", "create", "update", "patch", "delete", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: orca-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: orca-role
subjects:
  - kind: ServiceAccount
    name: orca
    namespace: orca-system
```

### 3. Create ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: orca-config
  namespace: orca-system
data:
  config.yaml: |
    aws:
      region: us-west-2
      vpcID: vpc-xxxxx
      subnetID: subnet-xxxxx
      securityGroupIDs:
        - sg-xxxxx
      tags:
        Environment: production
        Project: orca
    node:
      name: orca-aws-node
      # ... rest of config
```

### 4. Create Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: orca
  namespace: orca-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: orca
  template:
    metadata:
      labels:
        app: orca
    spec:
      serviceAccountName: orca
      containers:
        - name: orca
          image: orca:latest  # Build and push your image
          command:
            - /orca
          args:
            - --config=/config/config.yaml
            - --namespace=orca-system
          volumeMounts:
            - name: config
              mountPath: /config
          ports:
            - name: http
              containerPort: 8080
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 30
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
      volumes:
        - name: config
          configMap:
            name: orca-config
```

### 5. Deploy

```bash
kubectl apply -f deploy/
```

## Troubleshooting

### Node Not Appearing

```bash
# Check ORCA logs
kubectl logs -n orca-system deployment/orca

# Check RBAC permissions
kubectl auth can-i create nodes --as=system:serviceaccount:orca-system:orca
```

### Pods Stuck in Pending

```bash
# Check pod events
kubectl describe pod <pod-name>

# Check ORCA logs for instance creation errors
kubectl logs -n orca-system deployment/orca | grep "CreateInstance"

# Verify AWS credentials
kubectl exec -n orca-system deployment/orca -- env | grep AWS
```

### EC2 Instances Not Terminating

```bash
# Check ORCA logs
kubectl logs -n orca-system deployment/orca | grep "DeletePod"

# Manually check instances
aws ec2 describe-instances --filters "Name=tag:ManagedBy,Values=ORCA"

# Manually terminate if needed
aws ec2 terminate-instances --instance-ids i-xxxxx
```

## Next Steps

- [GPU Training Example](../examples/gpu-training-pod.yaml)
- [Spot Instance Example](../examples/spot-instance-pod.yaml)
- [Architecture Documentation](../architecture/overview.md)
- [Virtual Kubelet Integration](../architecture/virtual-kubelet.md)
- [Instance Selection Guide](../architecture/instance-selection.md)

## Getting Help

- GitHub Issues: https://github.com/scttfrdmn/orca/issues
- Documentation: https://github.com/scttfrdmn/orca/tree/main/docs
