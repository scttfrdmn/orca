// Package aws provides AWS EC2 integration for ORCA.
//
// It handles creating, managing, and terminating EC2 instances for pods.
// Supports both real AWS and LocalStack for testing.
//
// Example usage:
//
//	client, err := aws.NewClient(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	instanceID, err := client.CreateInstance(ctx, pod)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The client automatically tags instances for pod tracking and cleanup.
package aws
