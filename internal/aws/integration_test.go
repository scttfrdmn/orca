// +build integration

package aws

import (
	"context"
	"os"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcaconfig "github.com/scttfrdmn/orca/pkg/config"
)

// TestCreateInstanceLocalStack tests instance creation against LocalStack.
// Run with: go test -tags=integration ./internal/aws/...
func TestCreateInstanceLocalStack(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Check if LocalStack endpoint is available
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	if endpoint == "" {
		t.Skip("AWS_ENDPOINT_URL not set - skipping LocalStack test")
	}

	// Load LocalStack config
	cfg, err := orcaconfig.LoadConfig("../../config.localstack.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create AWS client
	ctx := context.Background()
	client, err := NewClient(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	// Create test pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "integration-test-pod",
			Namespace: "default",
			UID:       "test-uid-integration",
			Annotations: map[string]string{
				"orca.research/instance-type": "t3.small",
				"orca.research/launch-type":   "on-demand",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test",
					Image: "busybox",
				},
			},
		},
	}

	// Test 1: Create instance
	t.Run("create instance", func(t *testing.T) {
		instanceID, err := client.CreateInstance(ctx, pod, "t3.small")
		if err != nil {
			t.Fatalf("Failed to create instance: %v", err)
		}

		if instanceID == "" {
			t.Fatal("Instance ID is empty")
		}

		t.Logf("Created instance: %s", instanceID)

		// Test 2: Describe instance
		t.Run("describe instance", func(t *testing.T) {
			instance, err := client.DescribeInstance(ctx, instanceID)
			if err != nil {
				t.Fatalf("Failed to describe instance: %v", err)
			}

			if instance.ID != instanceID {
				t.Errorf("Expected instance ID %s, got %s", instanceID, instance.ID)
			}

			if instance.Type != "t3.small" {
				t.Errorf("Expected instance type t3.small, got %s", instance.Type)
			}

			t.Logf("Instance state: %s", instance.State)
		})

		// Test 3: Get instance by pod
		t.Run("get instance by pod", func(t *testing.T) {
			instance, err := client.GetInstanceByPod(ctx, "default", "integration-test-pod")
			if err != nil {
				t.Fatalf("Failed to get instance by pod: %v", err)
			}

			if instance.ID != instanceID {
				t.Errorf("Expected instance ID %s, got %s", instanceID, instance.ID)
			}

			t.Logf("Found instance by pod: %s", instance.ID)
		})

		// Test 4: Terminate instance
		t.Run("terminate instance", func(t *testing.T) {
			err := client.TerminateInstance(ctx, instanceID)
			if err != nil {
				t.Fatalf("Failed to terminate instance: %v", err)
			}

			t.Logf("Terminated instance: %s", instanceID)
		})
	})
}

// TestSpotInstanceLocalStack tests spot instance creation.
func TestSpotInstanceLocalStack(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	if endpoint == "" {
		t.Skip("AWS_ENDPOINT_URL not set - skipping LocalStack test")
	}

	cfg, err := orcaconfig.LoadConfig("../../config.localstack.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	client, err := NewClient(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to create AWS client: %v", err)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spot-test-pod",
			Namespace: "default",
			UID:       "test-uid-spot",
			Annotations: map[string]string{
				"orca.research/instance-type":   "t3.small",
				"orca.research/launch-type":     "spot",
				"orca.research/max-spot-price":  "0.05",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test",
					Image: "busybox",
				},
			},
		},
	}

	instanceID, err := client.CreateInstance(ctx, pod, "t3.small")
	if err != nil {
		t.Fatalf("Failed to create spot instance: %v", err)
	}

	t.Logf("Created spot instance: %s", instanceID)

	// Cleanup
	defer func() {
		if err := client.TerminateInstance(ctx, instanceID); err != nil {
			t.Logf("Failed to cleanup instance: %v", err)
		}
	}()

	instance, err := client.DescribeInstance(ctx, instanceID)
	if err != nil {
		t.Fatalf("Failed to describe instance: %v", err)
	}

	if instance.ID != instanceID {
		t.Errorf("Expected instance ID %s, got %s", instanceID, instance.ID)
	}
}
