package instances

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAutoSelector(t *testing.T) {
	tests := []struct {
		name     string
		pod      *corev1.Pod
		expected string
	}{
		{
			name: "single GPU requested",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gpu-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "gpu-container",
							Image: "pytorch:latest",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"nvidia.com/gpu": resource.MustParse("1"),
								},
							},
						},
					},
				},
			},
			expected: "g5.xlarge", // 1x A10G
		},
		{
			name: "2 GPUs requested",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gpu-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "gpu-container",
							Image: "pytorch:latest",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"nvidia.com/gpu": resource.MustParse("2"),
								},
							},
						},
					},
				},
			},
			expected: "g5.2xlarge",
		},
		{
			name: "4 GPUs requested",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gpu-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "gpu-container",
							Image: "pytorch:latest",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"nvidia.com/gpu": resource.MustParse("4"),
								},
							},
						},
					},
				},
			},
			expected: "g5.12xlarge",
		},
		{
			name: "8 GPUs requested",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gpu-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "gpu-container",
							Image: "pytorch:latest",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"nvidia.com/gpu": resource.MustParse("8"),
								},
							},
						},
					},
				},
			},
			expected: "p5.48xlarge", // 8x H100
		},
		{
			name: "small CPU workload",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cpu-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "cpu-container",
							Image: "busybox",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("1"),
									corev1.ResourceMemory: resource.MustParse("2Gi"),
								},
							},
						},
					},
				},
			},
			expected: "t3.small",
		},
		{
			name: "medium CPU workload",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cpu-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "cpu-container",
							Image: "busybox",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("3"),
									corev1.ResourceMemory: resource.MustParse("6Gi"),
								},
							},
						},
					},
				},
			},
			expected: "t3.large",
		},
		{
			name: "large CPU workload",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cpu-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "cpu-container",
							Image: "busybox",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("6"),
									corev1.ResourceMemory: resource.MustParse("12Gi"),
								},
							},
						},
					},
				},
			},
			expected: "c7i.2xlarge",
		},
		{
			name: "multiple containers aggregate resources",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "multi-container-pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "container1",
							Image: "busybox",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("2"),
									corev1.ResourceMemory: resource.MustParse("4Gi"),
								},
							},
						},
						{
							Name:  "container2",
							Image: "busybox",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("2"),
									corev1.ResourceMemory: resource.MustParse("4Gi"),
								},
							},
						},
					},
				},
			},
			expected: "t3.large", // Total: 4 CPUs, 8GB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := NewAutoSelector()
			result, err := selector.Select(tt.pod)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
