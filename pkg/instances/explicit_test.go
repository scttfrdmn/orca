package instances

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationInstanceType: "p5.48xlarge",
					},
				},
			},
			expected:    "p5.48xlarge",
			expectError: false,
		},
		{
			name: "explicit g5.4xlarge annotation",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationInstanceType: "g5.4xlarge",
					},
				},
			},
			expected:    "g5.4xlarge",
			expectError: false,
		},
		{
			name: "explicit t3.small annotation",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						AnnotationInstanceType: "t3.small",
					},
				},
			},
			expected:    "t3.small",
			expectError: false,
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
			name: "nil annotations returns error",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-pod",
					Namespace:   "default",
					Annotations: nil,
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
						AnnotationInstanceType: "",
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
