package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	orcaconfig "github.com/scttfrdmn/orca/pkg/config"
)

// mockEC2Client is a mock implementation of the EC2 API for testing.
type mockEC2Client struct {
	runInstancesFunc       func(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
	terminateInstancesFunc func(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
	describeInstancesFunc  func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func (m *mockEC2Client) RunInstances(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error) {
	if m.runInstancesFunc != nil {
		return m.runInstancesFunc(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *mockEC2Client) TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error) {
	if m.terminateInstancesFunc != nil {
		return m.terminateInstancesFunc(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *mockEC2Client) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	if m.describeInstancesFunc != nil {
		return m.describeInstancesFunc(ctx, params, optFns...)
	}
	return nil, nil
}

// testConfig returns a test configuration.
func testConfig() *orcaconfig.Config {
	return &orcaconfig.Config{
		AWS: orcaconfig.AWSConfig{
			Region:           "us-east-1",
			VPCID:            "vpc-12345",
			SubnetID:         "subnet-12345",
			SecurityGroupIDs: []string{"sg-12345"},
			AMIID:            "ami-12345",
		},
		Instances: orcaconfig.InstancesConfig{
			DefaultLaunchType: "on-demand",
			MaxSpotPrices: map[string]string{
				"t3.small": "0.05",
			},
		},
	}
}

// testPod returns a test pod.
func testPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			UID:       "test-uid-12345",
			Annotations: map[string]string{
				"orca.research/instance-type": "t3.small",
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
}

func TestGetLaunchType(t *testing.T) {
	cfg := testConfig()
	client := &Client{config: cfg}

	tests := []struct {
		name     string
		pod      *corev1.Pod
		expected string
	}{
		{
			name:     "default launch type when no annotation",
			pod:      testPod(),
			expected: "on-demand",
		},
		{
			name: "spot from annotation",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						"orca.research/launch-type": "spot",
					},
				},
			},
			expected: "spot",
		},
		{
			name: "on-demand from annotation",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						"orca.research/launch-type": "on-demand",
					},
				},
			},
			expected: "on-demand",
		},
		{
			name: "nil annotations returns default",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			expected: "on-demand",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.getLaunchType(tt.pod)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetMaxSpotPrice(t *testing.T) {
	cfg := testConfig()
	client := &Client{config: cfg}

	tests := []struct {
		name         string
		pod          *corev1.Pod
		instanceType string
		expected     string
	}{
		{
			name: "from pod annotation",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"orca.research/max-spot-price": "0.10",
					},
				},
			},
			instanceType: "t3.small",
			expected:     "0.10",
		},
		{
			name:         "from config",
			pod:          &corev1.Pod{},
			instanceType: "t3.small",
			expected:     "0.05",
		},
		{
			name:         "default when not configured",
			pod:          &corev1.Pod{},
			instanceType: "t3.large",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.getMaxSpotPrice(tt.pod, tt.instanceType)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetAMI(t *testing.T) {
	cfg := testConfig()
	client := &Client{config: cfg}

	tests := []struct {
		name     string
		pod      *corev1.Pod
		expected string
	}{
		{
			name: "from pod annotation",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"orca.research/ami": "ami-custom",
					},
				},
			},
			expected: "ami-custom",
		},
		{
			name:     "from config",
			pod:      &corev1.Pod{},
			expected: "ami-12345",
		},
		{
			name: "empty annotation uses config",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"orca.research/ami": "",
					},
				},
			},
			expected: "ami-12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.getAMI(tt.pod)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBuildTags(t *testing.T) {
	cfg := testConfig()
	client := &Client{config: cfg}

	pod := testPod()
	tags := client.buildTags(pod)

	// Verify required tags exist
	requiredTags := map[string]bool{
		"Name":                        false,
		"orca.research/pod":           false,
		"orca.research/pod-uid":       false,
		"orca.research/namespace":     false,
		"orca.research/provider":      false,
		"orca.research/created-at":    false,
	}

	for _, tag := range tags {
		if _, ok := requiredTags[*tag.Key]; ok {
			requiredTags[*tag.Key] = true
		}
	}

	for key, found := range requiredTags {
		if !found {
			t.Errorf("required tag %s not found", key)
		}
	}

	// Verify specific tag values
	for _, tag := range tags {
		switch *tag.Key {
		case "Name":
			expected := "orca-default-test-pod"
			if *tag.Value != expected {
				t.Errorf("Name tag: expected %s, got %s", expected, *tag.Value)
			}
		case "orca.research/pod":
			expected := "default/test-pod"
			if *tag.Value != expected {
				t.Errorf("pod tag: expected %s, got %s", expected, *tag.Value)
			}
		case "orca.research/namespace":
			if *tag.Value != "default" {
				t.Errorf("namespace tag: expected default, got %s", *tag.Value)
			}
		case "orca.research/provider":
			if *tag.Value != "orca" {
				t.Errorf("provider tag: expected orca, got %s", *tag.Value)
			}
		}
	}
}

func TestBuildTags_WithBudgetNamespace(t *testing.T) {
	cfg := testConfig()
	client := &Client{config: cfg}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				"orca.research/budget-namespace": "research-team-1",
			},
		},
	}

	tags := client.buildTags(pod)

	// Find budget namespace tag
	found := false
	for _, tag := range tags {
		if *tag.Key == "orca.research/budget-namespace" {
			found = true
			if *tag.Value != "research-team-1" {
				t.Errorf("budget-namespace tag: expected research-team-1, got %s", *tag.Value)
			}
		}
	}

	if !found {
		t.Error("budget-namespace tag not found")
	}
}

func TestBuildUserData(t *testing.T) {
	cfg := testConfig()
	client := &Client{config: cfg}

	t.Run("default user data", func(t *testing.T) {
		pod := testPod()
		userData := client.buildUserData(pod)

		if userData == "" {
			t.Error("user data should not be empty")
		}

		// Check that it contains pod information
		if !contains(userData, "default/test-pod") {
			t.Error("user data should contain pod namespace/name")
		}
	})

	t.Run("custom user data from annotation", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
				Annotations: map[string]string{
					"orca.research/user-data": "#!/bin/bash\necho custom",
				},
			},
		}

		userData := client.buildUserData(pod)
		expected := "#!/bin/bash\necho custom"
		if userData != expected {
			t.Errorf("expected %s, got %s", expected, userData)
		}
	})
}

func TestTerminateInstance(t *testing.T) {
	cfg := testConfig()

	t.Run("successful termination", func(t *testing.T) {
		mock := &mockEC2Client{
			terminateInstancesFunc: func(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error) {
				return &ec2.TerminateInstancesOutput{}, nil
			},
		}

		client := &Client{
			ec2Client: mock,
			config:    cfg,
		}

		err := client.TerminateInstance(context.Background(), "i-12345")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty instance ID", func(t *testing.T) {
		client := &Client{
			config: cfg,
		}

		err := client.TerminateInstance(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty instance ID")
		}
	})
}

func TestDescribeInstance(t *testing.T) {
	cfg := testConfig()

	t.Run("successful describe", func(t *testing.T) {
		now := time.Now()
		mock := &mockEC2Client{
			describeInstancesFunc: func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
				return &ec2.DescribeInstancesOutput{
					Reservations: []types.Reservation{
						{
							Instances: []types.Instance{
								{
									InstanceId:       aws.String("i-12345"),
									InstanceType:     types.InstanceTypeT3Small,
									State:            &types.InstanceState{Name: types.InstanceStateNameRunning},
									PublicIpAddress:  aws.String("1.2.3.4"),
									PrivateIpAddress: aws.String("10.0.1.5"),
									LaunchTime:       &now,
								},
							},
						},
					},
				}, nil
			},
		}

		client := &Client{
			ec2Client: mock,
			config:    cfg,
		}

		instance, err := client.DescribeInstance(context.Background(), "i-12345")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if instance.ID != "i-12345" {
			t.Errorf("expected ID i-12345, got %s", instance.ID)
		}
		if instance.Type != "t3.small" {
			t.Errorf("expected type t3.small, got %s", instance.Type)
		}
		if instance.State != "running" {
			t.Errorf("expected state running, got %s", instance.State)
		}
		if instance.PublicIP != "1.2.3.4" {
			t.Errorf("expected public IP 1.2.3.4, got %s", instance.PublicIP)
		}
		if instance.PrivateIP != "10.0.1.5" {
			t.Errorf("expected private IP 10.0.1.5, got %s", instance.PrivateIP)
		}
	})

	t.Run("empty instance ID", func(t *testing.T) {
		client := &Client{
			config: cfg,
		}

		_, err := client.DescribeInstance(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty instance ID")
		}
	})

	t.Run("instance not found", func(t *testing.T) {
		mock := &mockEC2Client{
			describeInstancesFunc: func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
				return &ec2.DescribeInstancesOutput{
					Reservations: []types.Reservation{},
				}, nil
			},
		}

		client := &Client{
			ec2Client: mock,
			config:    cfg,
		}

		_, err := client.DescribeInstance(context.Background(), "i-99999")
		if err == nil {
			t.Error("expected error for non-existent instance")
		}
	})
}

func TestGetInstanceByPod(t *testing.T) {
	cfg := testConfig()

	t.Run("successful lookup", func(t *testing.T) {
		now := time.Now()
		mock := &mockEC2Client{
			describeInstancesFunc: func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
				// Verify filters
				if len(params.Filters) < 2 {
					t.Error("expected at least 2 filters")
				}

				return &ec2.DescribeInstancesOutput{
					Reservations: []types.Reservation{
						{
							Instances: []types.Instance{
								{
									InstanceId:       aws.String("i-12345"),
									InstanceType:     types.InstanceTypeT3Small,
									State:            &types.InstanceState{Name: types.InstanceStateNameRunning},
									PublicIpAddress:  aws.String("1.2.3.4"),
									PrivateIpAddress: aws.String("10.0.1.5"),
									LaunchTime:       &now,
								},
							},
						},
					},
				}, nil
			},
		}

		client := &Client{
			ec2Client: mock,
			config:    cfg,
		}

		instance, err := client.GetInstanceByPod(context.Background(), "default", "test-pod")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if instance.ID != "i-12345" {
			t.Errorf("expected ID i-12345, got %s", instance.ID)
		}
	})

	t.Run("pod not found", func(t *testing.T) {
		mock := &mockEC2Client{
			describeInstancesFunc: func(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
				return &ec2.DescribeInstancesOutput{
					Reservations: []types.Reservation{},
				}, nil
			},
		}

		client := &Client{
			ec2Client: mock,
			config:    cfg,
		}

		_, err := client.GetInstanceByPod(context.Background(), "default", "missing-pod")
		if err == nil {
			t.Error("expected error for non-existent pod")
		}
	})
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
