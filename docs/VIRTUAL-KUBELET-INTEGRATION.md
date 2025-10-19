# Virtual Kubelet Integration

This document describes ORCA's Virtual Kubelet integration, which enables ORCA to register as a Kubernetes node and handle pod lifecycle events.

## Overview

ORCA uses the [Virtual Kubelet](https://virtual-kubelet.io/) framework to present itself as a node in a Kubernetes cluster. When pods are scheduled to the ORCA node, they are executed as EC2 instances on AWS rather than as containers on a physical node.

## Architecture

```
┌─────────────────────────────────────────┐
│ Kubernetes Control Plane                │
│ ┌─────────────┐                         │
│ │  Scheduler  │──────────────┐          │
│ └─────────────┘              │          │
│                               ▼          │
│                        ┌──────────────┐ │
│                        │  API Server  │ │
│                        └──────┬───────┘ │
└───────────────────────────────┼─────────┘
                                │
                                │ Watch Pods
                                │ Update Status
                                │
                    ┌───────────▼───────────┐
                    │   ORCA Virtual Node   │
                    │ ┌───────────────────┐ │
                    │ │ Node Controller   │ │
                    │ │ - Register Node   │ │
                    │ │ - Watch Pods      │ │
                    │ │ - Update Status   │ │
                    │ │ - Manage Lease    │ │
                    │ └────────┬──────────┘ │
                    │          │            │
                    │ ┌────────▼──────────┐ │
                    │ │ VK Adapter        │ │
                    │ │ - CreatePod       │ │
                    │ │ - DeletePod       │ │
                    │ │ - GetPodStatus    │ │
                    │ └────────┬──────────┘ │
                    │          │            │
                    │ ┌────────▼──────────┐ │
                    │ │ ORCA Provider     │ │
                    │ │ - Instance Mgmt   │ │
                    │ │ - AWS Client      │ │
                    │ └────────┬──────────┘ │
                    └──────────┼────────────┘
                               │
                               │ EC2 API
                               │
                    ┌──────────▼────────────┐
                    │   AWS EC2             │
                    │ ┌────┐ ┌────┐ ┌────┐ │
                    │ │ i1 │ │ i2 │ │ i3 │ │
                    │ └────┘ └────┘ └────┘ │
                    └───────────────────────┘
```

## Components

### pkg/node/controller.go

The `Controller` manages the Virtual Kubelet node lifecycle:

- **Initialization**: Creates Kubernetes client, ORCA provider, and Virtual Kubelet components
- **Node Registration**: Registers the virtual node with the Kubernetes API server
- **Lease Management**: Maintains node heartbeat via Kubernetes lease mechanism (40s intervals)
- **Graceful Shutdown**: Handles SIGTERM/SIGINT with proper cleanup

**Key Methods:**
```go
// NewController creates a new node controller
func NewController(cfg *config.Config, kubeconfigPath, namespace, version string, logger zerolog.Logger) (*Controller, error)

// Run starts the Virtual Kubelet node controller
func (c *Controller) Run(ctx context.Context) error

// Shutdown gracefully shuts down the controller
func (c *Controller) Shutdown(ctx context.Context) error
```

### pkg/node/adapter.go

The `VirtualKubeletAdapter` implements the Virtual Kubelet `PodLifecycleHandler` interface, bridging between Virtual Kubelet and the ORCA provider:

**Pod Lifecycle Methods:**
- `CreatePod`: Called when Kubernetes schedules a pod to this node
- `UpdatePod`: Called when pod spec is updated
- `DeletePod`: Called when pod is deleted
- `GetPod`: Retrieves pod by namespace/name
- `GetPods`: Lists all pods on this node
- `GetPodStatus`: Gets current pod status

**Exec/Logs Methods:**
- `GetContainerLogs`: Retrieves container logs (TODO: implement via CloudWatch)
- `RunInContainer`: Executes commands in container (TODO: implement via SSM)

**Node Methods:**
- `ConfigureNode`: Sets node capacity, labels, taints
- `NotifyNodeStatus`: Callback for node status updates
- `Ping`: Health check for provider responsiveness

## Pod Lifecycle Flow

### 1. Pod Creation

```
User submits pod → Scheduler assigns to orca-aws-node → 
CreatePod called → Instance selector chooses type →
EC2 instance launched → Pod status updated to Running
```

**Code Flow:**
```go
// 1. Virtual Kubelet calls CreatePod
adapter.CreatePod(ctx, pod)
  ↓
// 2. ORCA provider handles creation
provider.CreatePod(ctx, pod)
  ↓
// 3. Select instance type from annotations
selector.Select(pod) // Returns "p5.48xlarge"
  ↓
// 4. Launch EC2 instance
awsClient.CreateInstance(ctx, pod, instanceType)
  ↓
// 5. Tag instance with pod metadata
buildInstanceTags(pod, instanceType)
  ↓
// 6. Wait for instance running state
waiter.Wait(ctx, instanceID, 5*time.Minute)
  ↓
// 7. Update pod status
pod.Status.Phase = corev1.PodRunning
pod.Status.HostIP = instance.PublicIP
```

### 2. Pod Monitoring

ORCA continuously syncs pod status with EC2 instance state:

```go
// Virtual Kubelet periodically calls GetPodStatus
status := provider.GetPodStatus(ctx, namespace, name)
  ↓
// Query EC2 instance state
instance := awsClient.GetInstanceByPod(ctx, namespace, name)
  ↓
// Map EC2 state to Pod phase
switch instance.State {
case "running":  → corev1.PodRunning
case "pending":  → corev1.PodPending
case "stopped":  → corev1.PodFailed
}
```

### 3. Pod Deletion

```
kubectl delete pod → DeletePod called →
Find EC2 instance by tags → Terminate instance →
Pod removed from tracking
```

**Code Flow:**
```go
adapter.DeletePod(ctx, pod)
  ↓
provider.DeletePod(ctx, pod)
  ↓
// Find instance by pod tags
instance := awsClient.GetInstanceByPod(ctx, pod.Namespace, pod.Name)
  ↓
// Terminate instance
awsClient.TerminateInstance(ctx, instance.ID)
  ↓
// Remove from internal tracking
delete(provider.pods, pod.UID)
```

## Node Configuration

The virtual node is configured via `config.yaml`:

```yaml
node:
  name: orca-aws-node
  
  labels:
    orca.research/provider: "aws"
    orca.research/region: "us-west-2"
    type: virtual-kubelet
  
  taints:
    - key: orca.research/burst-node
      value: "true"
      effect: NoSchedule
  
  operatingSystem: Linux
  
  # Aggregate capacity (upper bounds)
  cpu: "1000"      # 1000 vCPUs
  memory: "4Ti"    # 4 TiB
  pods: "1000"     # Max 1000 pods
  gpu: "100"       # Max 100 GPUs
```

**Node Labels:**
- `orca.research/provider: aws` - Identifies ORCA nodes
- `orca.research/region: us-west-2` - AWS region
- `type: virtual-kubelet` - Standard Virtual Kubelet label

**Node Taints:**
- `orca.research/burst-node=true:NoSchedule` - Prevents regular pods from scheduling
  - Pods must explicitly tolerate this taint to burst to AWS

## Pod Annotations

Pods control their EC2 instance configuration via annotations:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: gpu-training
  annotations:
    # Required: Explicit instance type selection
    orca.research/instance-type: "p5.48xlarge"
    
    # Optional: Launch type (on-demand or spot)
    orca.research/launch-type: "spot"
    
    # Optional: Max spot price ($/hour)
    orca.research/max-spot-price: "35.00"
    
    # Optional: Custom AMI
    orca.research/ami: "ami-0123456789abcdef0"
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
      resources:
        limits:
          nvidia.com/gpu: 8
```

## Kubernetes Connection

ORCA supports three methods for connecting to Kubernetes:

### 1. In-Cluster Configuration (Production)

When running as a pod in Kubernetes:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: orca
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: orca-role
rules:
  - apiGroups: [""]
    resources: ["nodes", "pods", "pods/status"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: ["coordination.k8s.io"]
    resources: ["leases"]
    verbs: ["get", "create", "update", "patch", "delete"]
```

ORCA automatically uses in-cluster config when `--kubeconfig` is not specified.

### 2. Kubeconfig File (Development)

For local development:

```bash
orca --kubeconfig ~/.kube/config --config config.yaml
```

### 3. Default Kubeconfig

If no `--kubeconfig` specified and not in-cluster, ORCA tries `~/.kube/config`.

## Node Heartbeat & Lease

ORCA maintains its node registration via Kubernetes leases:

- **Lease Duration**: 40 seconds
- **Renew Interval**: Automatic (handled by Virtual Kubelet)
- **Namespace**: Same as ORCA deployment (typically `kube-system`)

If the lease expires (e.g., ORCA crashes), Kubernetes marks the node as `NotReady` and evicts pods.

## Node Status Reporting

ORCA reports node conditions to Kubernetes:

```go
node.Status.Conditions = []corev1.NodeCondition{
    {
        Type:   corev1.NodeReady,
        Status: corev1.ConditionTrue,
        Reason: "OrcaProviderReady",
    },
    {
        Type:   corev1.NodeMemoryPressure,
        Status: corev1.ConditionFalse,
    },
    // ... other conditions
}
```

**Node Info:**
```go
node.Status.NodeInfo = corev1.NodeSystemInfo{
    Architecture:            "amd64",
    OperatingSystem:         "Linux",
    KubeletVersion:          "v1.0.0-orca",
    ContainerRuntimeVersion: "orca://1.0.0",
    OSImage:                 "AWS EC2",
}
```

## Logging

ORCA uses structured logging (zerolog) throughout:

```bash
# JSON format (production)
{"level":"info","time":"2025-10-18T18:00:00Z","message":"Starting ORCA Virtual Kubelet node","node_name":"orca-aws-node"}

# Console format (development)
2025-10-18T18:00:00Z INF Starting ORCA Virtual Kubelet node node_name=orca-aws-node
```

**Log Levels:**
- `debug`: Detailed debugging information (AWS API calls, pod events)
- `info`: General operational information (node registered, pods created)
- `warn`: Warning conditions (spot instance interruption, quota limits)
- `error`: Error conditions (instance launch failed, AWS API errors)

Configure via config.yaml:
```yaml
logging:
  level: info     # debug, info, warn, error
  format: json    # json, text
```

## Error Handling

ORCA handles errors at multiple levels:

### Node Controller Errors

```go
// Startup errors (fatal - exit immediately)
if err := controller.Run(ctx); err != nil {
    logger.Fatal().Err(err).Msg("Controller error")
}

// Runtime errors (logged, retry automatically)
node.WithNodeStatusUpdateErrorHandler(func(ctx context.Context, err error) error {
    logger.Error().Err(err).Msg("Node status update failed")
    return err // Virtual Kubelet will retry
})
```

### Pod Creation Errors

```go
// Instance launch failure
if err := awsClient.CreateInstance(ctx, pod, instanceType); err != nil {
    // Update pod status to Failed
    pod.Status.Phase = corev1.PodFailed
    pod.Status.Conditions = append(pod.Status.Conditions, corev1.PodCondition{
        Type:    corev1.PodReady,
        Status:  corev1.ConditionFalse,
        Reason:  "InstanceCreationFailed",
        Message: fmt.Sprintf("Failed to create EC2 instance: %v", err),
    })
    return err
}
```

Common errors:
- `InsufficientInstanceCapacity`: No capacity available (especially for GPU instances)
- `InstanceLimitExceeded`: AWS account limits reached
- `UnauthorizedOperation`: IAM permissions issue
- `InvalidParameterValue`: Configuration error (bad AMI, subnet, etc.)

## Testing

### Unit Tests

```bash
# Test Virtual Kubelet adapter
go test ./pkg/node/... -v

# Test full integration
go test ./... -v
```

### Manual Testing

```bash
# 1. Start ORCA with local kubeconfig
./orca --kubeconfig ~/.kube/config --config config.yaml --log-level debug

# 2. Verify node registered
kubectl get nodes
# Should show: orca-aws-node   Ready   <none>   10s   v1.0.0-orca

# 3. Deploy test pod
kubectl apply -f examples/gpu-training-pod.yaml

# 4. Watch pod status
kubectl get pods -w

# 5. Check ORCA logs
# Should see: CreatePod called, instance launching, pod running

# 6. Verify EC2 instance created
aws ec2 describe-instances --filters "Name=tag:ManagedBy,Values=ORCA"

# 7. Delete pod
kubectl delete pod gpu-training

# 8. Verify instance terminated
aws ec2 describe-instances --instance-ids <instance-id>
```

## Troubleshooting

### Node Not Appearing

**Symptoms:** `kubectl get nodes` doesn't show orca-aws-node

**Causes:**
1. ORCA not running
2. Kubeconfig not configured
3. RBAC permissions missing
4. Network connectivity issues

**Solution:**
```bash
# Check ORCA logs
./orca --kubeconfig ~/.kube/config --log-level debug

# Verify RBAC
kubectl auth can-i create nodes --as=system:serviceaccount:kube-system:orca

# Test Kubernetes connectivity
kubectl cluster-info
```

### Pods Stuck in Pending

**Symptoms:** Pod scheduled to orca-aws-node but stays Pending

**Causes:**
1. Missing instance-type annotation
2. AWS credentials not configured
3. EC2 instance launch failure
4. Subnet/security group misconfiguration

**Solution:**
```bash
# Check pod events
kubectl describe pod <pod-name>

# Check ORCA logs
# Look for: "Failed to create instance" errors

# Verify AWS credentials
AWS_PROFILE=orca aws sts get-caller-identity

# Test EC2 instance launch manually
AWS_PROFILE=orca aws ec2 run-instances \
  --image-id ami-xxx \
  --instance-type t3.micro \
  --subnet-id subnet-xxx \
  --security-group-ids sg-xxx
```

### Pod Stuck in Unknown State

**Symptoms:** Pod shows "Unknown" or "NodeLost" status

**Causes:**
1. ORCA crashed or killed
2. Node lease expired
3. Kubernetes API connectivity lost

**Solution:**
```bash
# Check if ORCA is running
ps aux | grep orca

# Check node lease
kubectl get lease -n kube-system orca-aws-node

# Restart ORCA
./orca --kubeconfig ~/.kube/config --config config.yaml
```

## Security Considerations

1. **RBAC**: ORCA requires permissions to create/update nodes, pods, and leases
2. **AWS Credentials**: Use IRSA (IAM Roles for Service Accounts) in production
3. **Network Policy**: Ensure ORCA can reach Kubernetes API server
4. **Pod Security**: ORCA runs as non-root user (UID 65532) in container

## Performance

- **Node Registration**: ~1-2 seconds
- **Pod Creation**: ~60-90 seconds (EC2 instance boot time)
- **Pod Status Sync**: Every 10 seconds (Virtual Kubelet default)
- **Node Heartbeat**: Every 40 seconds (lease renewal)

## Limitations

1. **Container Runtime**: Pods run as full EC2 instances, not containers
   - Cannot use Docker/containerd features
   - No shared node resources
   - Higher overhead than container pods

2. **Exec/Logs**: Not yet implemented
   - `kubectl logs` returns "not implemented"
   - `kubectl exec` returns "not implemented"
   - Future: Will use CloudWatch Logs and SSM Session Manager

3. **Volume Mounts**: Not yet implemented
   - EmptyDir, HostPath not supported
   - Future: Will support EBS volumes and EFS

4. **Networking**: Simplified model
   - Each pod gets its own EC2 instance with public/private IP
   - No pod-to-pod networking within node
   - No CNI plugin integration

## Future Enhancements

- [ ] CloudWatch Logs integration for container logs
- [ ] SSM Session Manager for kubectl exec
- [ ] EBS volume support
- [ ] EFS volume support
- [ ] Multi-container pod support (multiple processes on single instance)
- [ ] Init container support
- [ ] Ephemeral container support
- [ ] Resource metrics via CloudWatch
- [ ] Custom metrics API integration

## References

- [Virtual Kubelet Documentation](https://virtual-kubelet.io/)
- [Virtual Kubelet GitHub](https://github.com/virtual-kubelet/virtual-kubelet)
- [Kubernetes Node Documentation](https://kubernetes.io/docs/concepts/architecture/nodes/)
- [Kubernetes Pod Lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/)
