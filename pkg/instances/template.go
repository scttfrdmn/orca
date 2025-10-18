package instances

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/scttfrdmn/orca/pkg/config"
)

const (
	// AnnotationWorkloadTemplate is the pod annotation for workload template.
	AnnotationWorkloadTemplate = "orca.research/workload-template"
)

// TemplateSelector selects instance type from predefined templates.
// Templates provide convenient defaults for common workload types.
type TemplateSelector struct {
	templates map[string]config.WorkloadTemplate
}

// NewTemplateSelector creates a new template selector.
func NewTemplateSelector(templates map[string]config.WorkloadTemplate) *TemplateSelector {
	return &TemplateSelector{
		templates: templates,
	}
}

// Select returns the instance type from a workload template.
func (s *TemplateSelector) Select(pod *corev1.Pod) (string, error) {
	if pod.Annotations == nil {
		return "", fmt.Errorf("pod has no annotations")
	}

	templateName, ok := pod.Annotations[AnnotationWorkloadTemplate]
	if !ok {
		return "", fmt.Errorf("pod missing annotation: %s", AnnotationWorkloadTemplate)
	}

	if templateName == "" {
		return "", fmt.Errorf("annotation %s is empty", AnnotationWorkloadTemplate)
	}

	template, ok := s.templates[templateName]
	if !ok {
		return "", fmt.Errorf("unknown workload template: %s", templateName)
	}

	if template.InstanceType == "" {
		return "", fmt.Errorf("template %s has no instance type defined", templateName)
	}

	return template.InstanceType, nil
}
