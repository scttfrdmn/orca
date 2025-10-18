// Package instances provides instance type selection for ORCA pods.
//
// It implements a three-tier selection strategy:
//
//  1. Explicit: User specifies exact instance type via annotation
//     Example: orca.research/instance-type: "p5.48xlarge"
//
//  2. Template: Named workload templates for common use cases
//     Example: orca.research/workload-template: "llm-training"
//
//  3. Auto: Automatic selection based on pod resource requests
//     Example: 8 GPUs requested -> p5.48xlarge
//
// The selector chain tries each strategy in order until one succeeds.
//
// Example usage:
//
//	selector, err := instances.NewSelector(cfg.Instances)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	instanceType, err := selector.Select(pod)
package instances
