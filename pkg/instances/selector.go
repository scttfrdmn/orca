package instances

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/scttfrdmn/orca/pkg/config"
)

// Selector chooses the appropriate EC2 instance type for a pod.
type Selector interface {
	// Select returns the EC2 instance type to use for the given pod.
	Select(pod *corev1.Pod) (string, error)
}

// ChainSelector implements a priority chain of selectors.
// It tries each selector in order until one succeeds.
type ChainSelector struct {
	selectors []Selector
}

// NewSelector creates a new instance selector based on configuration.
func NewSelector(cfg config.InstancesConfig) (Selector, error) {
	// Build selector chain based on selection mode
	var selectors []Selector

	switch cfg.SelectionMode {
	case "explicit":
		// Explicit only - no fallback
		selectors = []Selector{
			NewExplicitSelector(),
		}

	case "template":
		// Template with explicit fallback
		selectors = []Selector{
			NewExplicitSelector(),
			NewTemplateSelector(cfg.Templates),
		}

	case "auto":
		// Full chain: explicit → template → auto
		selectors = []Selector{
			NewExplicitSelector(),
			NewTemplateSelector(cfg.Templates),
			NewAutoSelector(),
		}

	default:
		return nil, fmt.Errorf("invalid selection mode: %s (must be explicit, template, or auto)", cfg.SelectionMode)
	}

	return &ChainSelector{
		selectors: selectors,
	}, nil
}

// Select tries each selector in the chain until one succeeds.
func (c *ChainSelector) Select(pod *corev1.Pod) (string, error) {
	var lastErr error

	for _, selector := range c.selectors {
		instanceType, err := selector.Select(pod)
		if err == nil {
			return instanceType, nil
		}
		lastErr = err
	}

	// All selectors failed
	return "", fmt.Errorf("no selector could determine instance type: %w", lastErr)
}
