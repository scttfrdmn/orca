package instances

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/scttfrdmn/orca/pkg/config"
)

func TestTemplateSelector(t *testing.T) {
	templates := map[string]config.WorkloadTemplate{
		"llm-training": {
			InstanceType: "p5.48xlarge",
			LaunchType:   "spot",
			MaxSpotPrice: "50.00",
		},
		"vision-training": {
			InstanceType: "g5.4xlarge",
			LaunchType:   "spot",
			MaxSpotPrice: "2.00",
		},
		"inference": {
			InstanceType: "g6.2xlarge",
			LaunchType:   "on-demand",
		},
		"empty-template": {
			// Intentionally no instance type
			LaunchType: "on-demand",
		},
	}

	tests := []struct {
		name        string
		pod         *corev1.Pod
		expected    string
		expectError bool
	}{
		{
			name: "llm-training template",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationWorkloadTemplate: "llm-training",
					},
				},
			},
			expected:    "p5.48xlarge",
			expectError: false,
		},
		{
			name: "vision-training template",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationWorkloadTemplate: "vision-training",
					},
				},
			},
			expected:    "g5.4xlarge",
			expectError: false,
		},
		{
			name: "inference template",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationWorkloadTemplate: "inference",
					},
				},
			},
			expected:    "g6.2xlarge",
			expectError: false,
		},
		{
			name: "unknown template returns error",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationWorkloadTemplate: "unknown-template",
					},
				},
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "missing annotation returns error",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-pod",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "empty annotation value returns error",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationWorkloadTemplate: "",
					},
				},
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "template with no instance type returns error",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationWorkloadTemplate: "empty-template",
					},
				},
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := NewTemplateSelector(templates)
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
