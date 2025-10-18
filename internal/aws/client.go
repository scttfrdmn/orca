package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	corev1 "k8s.io/api/core/v1"

	orcaconfig "github.com/scttfrdmn/orca/pkg/config"
)

// Client manages AWS EC2 operations for ORCA.
type Client struct {
	ec2Client *ec2.Client
	config    *orcaconfig.Config
}

// NewClient creates a new AWS client.
// Supports both real AWS and LocalStack via endpoint override.
func NewClient(ctx context.Context, cfg *orcaconfig.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Build AWS config
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.AWS.Region),
	}

	// Use explicit credentials if provided
	if cfg.AWS.Credentials != nil {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWS.Credentials.AccessKeyID,
				cfg.AWS.Credentials.SecretAccessKey,
				"",
			),
		))
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create EC2 client
	ec2Opts := []func(*ec2.Options){}

	// Support LocalStack endpoint override
	if endpoint := cfg.AWS.LocalStackEndpoint; endpoint != "" {
		ec2Opts = append(ec2Opts, func(o *ec2.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	ec2Client := ec2.NewFromConfig(awsConfig, ec2Opts...)

	return &Client{
		ec2Client: ec2Client,
		config:    cfg,
	}, nil
}

// CreateInstance creates a new EC2 instance for a pod.
func (c *Client) CreateInstance(ctx context.Context, pod *corev1.Pod, instanceType string) (string, error) {
	if pod == nil {
		return "", fmt.Errorf("pod cannot be nil")
	}
	if instanceType == "" {
		return "", fmt.Errorf("instanceType cannot be empty")
	}

	// Determine launch type (on-demand or spot)
	launchType := c.getLaunchType(pod)

	// Build instance tags
	tags := c.buildTags(pod)

	// Build user data
	userData := c.buildUserData(pod)

	// Prepare run instances input
	input := &ec2.RunInstancesInput{
		ImageId:          aws.String(c.getAMI(pod)),
		InstanceType:     types.InstanceType(instanceType),
		MinCount:         aws.Int32(1),
		MaxCount:         aws.Int32(1),
		SubnetId:         aws.String(c.config.AWS.SubnetID),
		SecurityGroupIds: c.config.AWS.SecurityGroupIDs,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         tags,
			},
		},
		UserData: aws.String(userData),
	}

	// Handle spot instances
	if launchType == "spot" {
		maxPrice := c.getMaxSpotPrice(pod, instanceType)
		input.InstanceMarketOptions = &types.InstanceMarketOptionsRequest{
			MarketType: types.MarketTypeSpot,
			SpotOptions: &types.SpotMarketOptions{
				MaxPrice:         aws.String(maxPrice),
				SpotInstanceType: types.SpotInstanceTypeOneTime,
			},
		}
	}

	// Launch instance
	result, err := c.ec2Client.RunInstances(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to run instance: %w", err)
	}

	if len(result.Instances) == 0 {
		return "", fmt.Errorf("no instances created")
	}

	instanceID := *result.Instances[0].InstanceId

	// Wait for instance to be running (with timeout)
	if err := c.waitForInstanceRunning(ctx, instanceID); err != nil {
		// Attempt to terminate failed instance
		_ = c.TerminateInstance(ctx, instanceID)
		return "", fmt.Errorf("instance failed to start: %w", err)
	}

	return instanceID, nil
}

// TerminateInstance terminates an EC2 instance.
func (c *Client) TerminateInstance(ctx context.Context, instanceID string) error {
	if instanceID == "" {
		return fmt.Errorf("instanceID cannot be empty")
	}

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	}

	_, err := c.ec2Client.TerminateInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to terminate instance: %w", err)
	}

	return nil
}

// DescribeInstance retrieves information about an EC2 instance.
func (c *Client) DescribeInstance(ctx context.Context, instanceID string) (*Instance, error) {
	if instanceID == "" {
		return nil, fmt.Errorf("instanceID cannot be empty")
	}

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := c.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found: %s", instanceID)
	}

	ec2Instance := result.Reservations[0].Instances[0]

	return &Instance{
		ID:           *ec2Instance.InstanceId,
		Type:         string(ec2Instance.InstanceType),
		State:        string(ec2Instance.State.Name),
		PublicIP:     aws.ToString(ec2Instance.PublicIpAddress),
		PrivateIP:    aws.ToString(ec2Instance.PrivateIpAddress),
		LaunchTime:   aws.ToTime(ec2Instance.LaunchTime),
		InstanceType: string(ec2Instance.InstanceType),
	}, nil
}

// GetInstanceByPod finds the EC2 instance for a given pod.
func (c *Client) GetInstanceByPod(ctx context.Context, namespace, name string) (*Instance, error) {
	podKey := fmt.Sprintf("%s/%s", namespace, name)

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:orca.research/pod"),
				Values: []string{podKey},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"pending", "running", "stopping", "stopped"},
			},
		},
	}

	result, err := c.ec2Client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found for pod %s", podKey)
	}

	ec2Instance := result.Reservations[0].Instances[0]

	return &Instance{
		ID:           *ec2Instance.InstanceId,
		Type:         string(ec2Instance.InstanceType),
		State:        string(ec2Instance.State.Name),
		PublicIP:     aws.ToString(ec2Instance.PublicIpAddress),
		PrivateIP:    aws.ToString(ec2Instance.PrivateIpAddress),
		LaunchTime:   aws.ToTime(ec2Instance.LaunchTime),
		InstanceType: string(ec2Instance.InstanceType),
	}, nil
}

// waitForInstanceRunning waits for an instance to be in running state.
func (c *Client) waitForInstanceRunning(ctx context.Context, instanceID string) error {
	timeout := 5 * time.Minute
	interval := 10 * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for instance to be running")
		case <-ticker.C:
			instance, err := c.DescribeInstance(ctx, instanceID)
			if err != nil {
				return err
			}

			switch instance.State {
			case "running":
				return nil
			case "terminated", "shutting-down":
				return fmt.Errorf("instance entered terminal state: %s", instance.State)
			case "pending":
				// Continue waiting
				continue
			default:
				return fmt.Errorf("unexpected instance state: %s", instance.State)
			}
		}
	}
}

// getLaunchType determines whether to use on-demand or spot instances.
func (c *Client) getLaunchType(pod *corev1.Pod) string {
	if pod.Annotations == nil {
		return c.config.Instances.DefaultLaunchType
	}

	if launchType, ok := pod.Annotations["orca.research/launch-type"]; ok {
		return launchType
	}

	return c.config.Instances.DefaultLaunchType
}

// getMaxSpotPrice returns the maximum spot price for an instance type.
func (c *Client) getMaxSpotPrice(pod *corev1.Pod, instanceType string) string {
	// Check pod annotation first
	if pod.Annotations != nil {
		if maxPrice, ok := pod.Annotations["orca.research/max-spot-price"]; ok {
			return maxPrice
		}
	}

	// Check configuration
	if maxPrice, ok := c.config.Instances.MaxSpotPrices[instanceType]; ok {
		return maxPrice
	}

	// Default: don't set a max price (use on-demand price as max)
	return ""
}

// getAMI returns the AMI ID to use for the instance.
func (c *Client) getAMI(pod *corev1.Pod) string {
	// Check pod annotation first
	if pod.Annotations != nil {
		if ami, ok := pod.Annotations["orca.research/ami"]; ok && ami != "" {
			return ami
		}
	}

	// Use configured default AMI
	if c.config.AWS.AMIID != "" {
		return c.config.AWS.AMIID
	}

	// Fallback to Amazon Linux 2023
	// TODO: Make this region-specific
	return "ami-0c55b159cbfafe1f0"
}

// buildTags creates EC2 tags for pod tracking.
func (c *Client) buildTags(pod *corev1.Pod) []types.Tag {
	podKey := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)

	tags := []types.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(fmt.Sprintf("orca-%s-%s", pod.Namespace, pod.Name)),
		},
		{
			Key:   aws.String("orca.research/pod"),
			Value: aws.String(podKey),
		},
		{
			Key:   aws.String("orca.research/pod-uid"),
			Value: aws.String(string(pod.UID)),
		},
		{
			Key:   aws.String("orca.research/namespace"),
			Value: aws.String(pod.Namespace),
		},
		{
			Key:   aws.String("orca.research/provider"),
			Value: aws.String("orca"),
		},
		{
			Key:   aws.String("orca.research/created-at"),
			Value: aws.String(time.Now().Format(time.RFC3339)),
		},
	}

	// Add budget namespace tag if specified
	if pod.Annotations != nil {
		if budgetNS, ok := pod.Annotations["orca.research/budget-namespace"]; ok {
			tags = append(tags, types.Tag{
				Key:   aws.String("orca.research/budget-namespace"),
				Value: aws.String(budgetNS),
			})
		}
	}

	return tags
}

// buildUserData creates user data script for the instance.
func (c *Client) buildUserData(pod *corev1.Pod) string {
	// Check for custom user data
	if pod.Annotations != nil {
		if userData, ok := pod.Annotations["orca.research/user-data"]; ok {
			return userData
		}
	}

	// Default user data: basic setup
	return `#!/bin/bash
# ORCA instance initialization
echo "ORCA instance starting for pod: ` + pod.Namespace + `/` + pod.Name + `"
`
}
