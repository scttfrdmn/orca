// Package provider implements the Virtual Kubelet provider interface for ORCA.
//
// The provider handles pod lifecycle operations by creating and managing
// EC2 instances in AWS. Each pod is mapped to a dedicated EC2 instance,
// with the instance type determined by pod annotations or templates.
//
// Example usage:
//
//	provider, err := provider.NewProvider(cfg, "orca-node", "kube-system", "v0.1.0")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	err = provider.CreatePod(ctx, pod)
//
// The provider supports explicit instance selection, template-based selection,
// and automatic selection based on pod resource requests.
package provider
