package instances

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

const (
	// AnnotationInstanceType is the pod annotation for explicit instance type.
	AnnotationInstanceType = "orca.research/instance-type"
)

// ExplicitSelector selects instance type from explicit pod annotation.
// This is the highest priority selector (user knows exactly what they want).
type ExplicitSelector struct{}

// NewExplicitSelector creates a new explicit selector.
func NewExplicitSelector() *ExplicitSelector {
	return &ExplicitSelector{}
}

// Select returns the instance type from pod annotation.
func (s *ExplicitSelector) Select(pod *corev1.Pod) (string, error) {
	if pod.Annotations == nil {
		return "", fmt.Errorf("pod has no annotations")
	}

	instanceType, ok := pod.Annotations[AnnotationInstanceType]
	if !ok {
		return "", fmt.Errorf("pod missing annotation: %s", AnnotationInstanceType)
	}

	if instanceType == "" {
		return "", fmt.Errorf("annotation %s is empty", AnnotationInstanceType)
	}

	// TODO: Validate instance type format (e.g., "p5.48xlarge")
	// TODO: Check against allowed instance types list

	return instanceType, nil
}
