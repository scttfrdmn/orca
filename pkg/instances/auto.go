package instances

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// AutoSelector automatically selects instance type based on pod resource requests.
// This is the fallback selector when no explicit instance type or template is specified.
type AutoSelector struct{}

// NewAutoSelector creates a new auto selector.
func NewAutoSelector() *AutoSelector {
	return &AutoSelector{}
}

// Select returns an instance type based on pod resource requirements.
func (s *AutoSelector) Select(pod *corev1.Pod) (string, error) {
	// Calculate total resource requests
	var totalCPU resource.Quantity
	var totalMemory resource.Quantity
	var gpuCount int64

	for _, container := range pod.Spec.Containers {
		if cpu, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
			totalCPU.Add(cpu)
		}
		if memory, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
			totalMemory.Add(memory)
		}
		if gpu, ok := container.Resources.Requests["nvidia.com/gpu"]; ok {
			gpuCount += gpu.Value()
		}
	}

	// If GPU requested, select GPU instance
	if gpuCount > 0 {
		return s.selectGPUInstance(gpuCount), nil
	}

	// Otherwise select CPU instance
	return s.selectCPUInstance(totalCPU.MilliValue(), totalMemory.Value()), nil
}

// selectGPUInstance selects an appropriate GPU instance type.
func (s *AutoSelector) selectGPUInstance(gpuCount int64) string {
	// Simple heuristic: smallest instance that fits
	switch gpuCount {
	case 1:
		return "g5.xlarge" // 1x A10G
	case 2:
		return "g5.2xlarge" // 2x A10G
	case 4:
		return "g5.12xlarge" // 4x A10G
	case 8:
		return "p5.48xlarge" // 8x H100
	default:
		// Default to single GPU instance
		return "g5.xlarge"
	}
}

// selectCPUInstance selects an appropriate CPU instance type.
func (s *AutoSelector) selectCPUInstance(cpuMilli int64, memoryBytes int64) string {
	// Convert to vCPUs and GB
	vCPUs := cpuMilli / 1000
	memoryGB := memoryBytes / (1024 * 1024 * 1024)

	// Simple heuristic: t3 family for small workloads, c7i for larger
	if vCPUs <= 2 && memoryGB <= 4 {
		return "t3.small" // 2 vCPU, 2 GB
	} else if vCPUs <= 4 && memoryGB <= 8 {
		return "t3.large" // 2 vCPU, 8 GB
	} else if vCPUs <= 8 && memoryGB <= 16 {
		return "c7i.2xlarge" // 8 vCPU, 16 GB
	} else if vCPUs <= 16 && memoryGB <= 32 {
		return "c7i.4xlarge" // 16 vCPU, 32 GB
	} else if vCPUs <= 32 && memoryGB <= 64 {
		return "c7i.8xlarge" // 32 vCPU, 64 GB
	} else {
		return "c7i.16xlarge" // 64 vCPU, 128 GB
	}
}
