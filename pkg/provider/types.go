package provider

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// Provider is the Virtual Kubelet provider interface for ORCA.
// It handles pod lifecycle operations by creating and managing EC2 instances.
type Provider interface {
	// Pod Lifecycle Methods

	// CreatePod creates a new pod by launching an EC2 instance.
	// The instance type is determined by pod annotations.
	CreatePod(ctx context.Context, pod *corev1.Pod) error

	// UpdatePod updates an existing pod.
	// Currently limited operations are supported (mainly status updates).
	UpdatePod(ctx context.Context, pod *corev1.Pod) error

	// DeletePod deletes a pod by terminating its EC2 instance.
	DeletePod(ctx context.Context, pod *corev1.Pod) error

	// GetPod retrieves a pod by namespace and name.
	GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error)

	// GetPodStatus retrieves the status of a pod.
	GetPodStatus(ctx context.Context, namespace, name string) (*corev1.PodStatus, error)

	// GetPods retrieves all pods managed by this provider.
	GetPods(ctx context.Context) ([]*corev1.Pod, error)

	// Node Information Methods

	// ConfigureNode sets up the virtual node with appropriate capacity and labels.
	ConfigureNode(ctx context.Context, node *corev1.Node)

	// GetNodeStatus returns the current status of the virtual node.
	GetNodeStatus(ctx context.Context) (*corev1.NodeStatus, error)

	// Container Operations (for kubectl logs, exec)

	// GetContainerLogs retrieves logs from a container.
	GetContainerLogs(ctx context.Context, namespace, podName, containerName string, opts ContainerLogOpts) ([]byte, error)

	// RunInContainer executes a command in a container (for kubectl exec).
	RunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string, attach AttachIO) error

	// GetStatsSummary retrieves resource usage statistics.
	GetStatsSummary(ctx context.Context) (*StatsSummary, error)
}

// ContainerLogOpts specifies options for retrieving container logs.
type ContainerLogOpts struct {
	Tail         int
	Follow       bool
	Previous     bool
	SinceSeconds int
	SinceTime    time.Time
	Timestamps   bool
}

// AttachIO represents streams for interactive container operations.
type AttachIO interface {
	Stdin() []byte
	Stdout() chan []byte
	Stderr() chan []byte
	TTY() bool
}

// StatsSummary represents resource usage statistics for the node.
type StatsSummary struct {
	Node NodeStats  `json:"node"`
	Pods []PodStats `json:"pods"`
}

// NodeStats represents resource usage for the virtual node.
type NodeStats struct {
	NodeName  string       `json:"nodeName"`
	StartTime time.Time    `json:"startTime"`
	CPU       *CPUStats    `json:"cpu,omitempty"`
	Memory    *MemoryStats `json:"memory,omitempty"`
}

// PodStats represents resource usage for a pod.
type PodStats struct {
	PodRef    PodReference `json:"podRef"`
	StartTime time.Time    `json:"startTime"`
	CPU       *CPUStats    `json:"cpu,omitempty"`
	Memory    *MemoryStats `json:"memory,omitempty"`
}

// PodReference identifies a pod.
type PodReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	UID       string `json:"uid"`
}

// CPUStats represents CPU usage statistics.
type CPUStats struct {
	Time                 time.Time `json:"time"`
	UsageNanoCores       *uint64   `json:"usageNanoCores,omitempty"`
	UsageCoreNanoSeconds *uint64   `json:"usageCoreNanoSeconds,omitempty"`
}

// MemoryStats represents memory usage statistics.
type MemoryStats struct {
	Time            time.Time `json:"time"`
	UsageBytes      *uint64   `json:"usageBytes,omitempty"`
	WorkingSetBytes *uint64   `json:"workingSetBytes,omitempty"`
}
