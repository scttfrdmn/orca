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

const (
	// ORCA annotation constants (duplicated here to avoid import cycle)
	annotationLaunchType = "orca.research/launch-type"
)

// Client is the AWS EC2 client for ORCA operations.
type Client struct {
	ec2Client *ec2.Client
	config    *orcaconfig.Config
}

// NewClient creates a new AWS client.
func NewClient(ctx context.Context, cfg *orcaconfig.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Build AWS SDK config
	var awsConfigOpts []func(*config.LoadOptions) error

	// Set region
	awsConfigOpts = append(awsConfigOpts, config.WithRegion(cfg.AWS.Region))

	// Set credentials if provided (otherwise uses default credential chain)
	if cfg.AWS.Credentials != nil {
		creds := credentials.NewStaticCredentialsProvider(
			cfg.AWS.Credentials.AccessKeyID,
			cfg.AWS.Credentials.SecretAccessKey,
			"",
		)
		awsConfigOpts = append(awsConfigOpts, config.WithCredentialsProvider(creds))
	}

	// Load AWS config
	awsConfig, err := config.LoadDefaultConfig(ctx, awsConfigOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create EC2 client
	ec2Client := ec2.NewFromConfig(awsConfig)

	// Override endpoint for LocalStack if configured
	if cfg.AWS.LocalStackEndpoint != "" {
		ec2Client = ec2.NewFromConfig(awsConfig, func(o *ec2.Options) {
			o.BaseEndpoint = aws.String(cfg.AWS.LocalStackEndpoint)
		})
	}

	return &Client{
		ec2Client: ec2Client,
		config:    cfg,
	}, nil
}

// CreateInstance creates an EC2 instance for a pod.
func (c *Client) CreateInstance(ctx context.Context, pod *corev1.Pod, instanceType string) (string, error) {
	if pod == nil {
		return "", fmt.Errorf("pod cannot be nil")
	}

	// Extract launch type from annotations
	launchType := c.config.Instances.DefaultLaunchType
	if lt, ok := pod.Annotations[annotationLaunchType]; ok {
		launchType = lt
	}

	// Build instance tags
	tags := c.buildInstanceTags(pod, instanceType)
	tagSpecs := []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeInstance,
			Tags:         tags,
		},
		{
			ResourceType: types.ResourceTypeVolume,
			Tags:         tags,
		},
	}

	// Build RunInstances input
	runInput := &ec2.RunInstancesInput{
		MaxCount:          aws.Int32(1),
		MinCount:          aws.Int32(1),
		InstanceType:      types.InstanceType(instanceType),
		SubnetId:          aws.String(c.config.AWS.SubnetID),
		SecurityGroupIds:  c.config.AWS.SecurityGroupIDs,
		TagSpecifications: tagSpecs,
	}

	// Set AMI (either from config or use latest Amazon Linux 2023)
	if c.config.AWS.AMIID != "" {
		runInput.ImageId = aws.String(c.config.AWS.AMIID)
	} else {
		// TODO: Implement AMI lookup for latest Amazon Linux 2023
		// For now, require AMI to be specified
		return "", fmt.Errorf("aws.amiID must be specified in config")
	}

	// Configure spot instances if requested
	if launchType == "spot" {
		runInput.InstanceMarketOptions = &types.InstanceMarketOptionsRequest{
			MarketType: types.MarketTypeSpot,
			SpotOptions: &types.SpotMarketOptions{
				SpotInstanceType: types.SpotInstanceTypeOneTime,
			},
		}

		// Set max spot price if configured
		if maxPrice, ok := c.config.Instances.MaxSpotPrices[instanceType]; ok {
			runInput.InstanceMarketOptions.SpotOptions.MaxPrice = aws.String(maxPrice)
		}
	}

	// Launch the instance
	result, err := c.ec2Client.RunInstances(ctx, runInput)
	if err != nil {
		return "", fmt.Errorf("failed to launch instance: %w", err)
	}

	if len(result.Instances) == 0 {
		return "", fmt.Errorf("no instances were created")
	}

	instanceID := *result.Instances[0].InstanceId

	// Wait for instance to be running (with timeout)
	waiter := ec2.NewInstanceRunningWaiter(c.ec2Client)
	if err := waiter.Wait(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}, 5*time.Minute); err != nil {
		// Cleanup: terminate the instance if it fails to start
		_ = c.TerminateInstance(ctx, instanceID)
		return "", fmt.Errorf("instance failed to reach running state: %w", err)
	}

	return instanceID, nil
}

// TerminateInstance terminates an EC2 instance.
func (c *Client) TerminateInstance(ctx context.Context, instanceID string) error {
	if instanceID == "" {
		return fmt.Errorf("instanceID cannot be empty")
	}

	_, err := c.ec2Client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return fmt.Errorf("failed to terminate instance %s: %w", instanceID, err)
	}

	return nil
}

// GetInstanceByPod retrieves an instance by pod namespace and name using tags.
func (c *Client) GetInstanceByPod(ctx context.Context, namespace, name string) (*Instance, error) {
	// Search for instance by pod tags
	result, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:orca.research/pod-namespace"),
				Values: []string{namespace},
			},
			{
				Name:   aws.String("tag:orca.research/pod-name"),
				Values: []string{name},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					string(types.InstanceStateNamePending),
					string(types.InstanceStateNameRunning),
					string(types.InstanceStateNameStopping),
					string(types.InstanceStateNameStopped),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe instances: %w", err)
	}

	// Find the instance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			return c.convertInstance(&instance), nil
		}
	}

	return nil, fmt.Errorf("instance not found for pod %s/%s", namespace, name)
}

// GetInstance retrieves an instance by ID.
func (c *Client) GetInstance(ctx context.Context, instanceID string) (*Instance, error) {
	if instanceID == "" {
		return nil, fmt.Errorf("instanceID cannot be empty")
	}

	result, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}

	return c.convertInstance(&result.Reservations[0].Instances[0]), nil
}

// ListInstances lists all ORCA-managed instances.
func (c *Client) ListInstances(ctx context.Context) ([]*Instance, error) {
	result, err := c.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("tag:ManagedBy"),
				Values: []string{"ORCA"},
			},
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					string(types.InstanceStateNamePending),
					string(types.InstanceStateNameRunning),
					string(types.InstanceStateNameStopping),
					string(types.InstanceStateNameStopped),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	var instances []*Instance
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			inst := instance // Create local copy for pointer
			instances = append(instances, c.convertInstance(&inst))
		}
	}

	return instances, nil
}

// buildInstanceTags builds EC2 tags for an instance.
func (c *Client) buildInstanceTags(pod *corev1.Pod, instanceType string) []types.Tag {
	// Get pod-specific tags from config
	tagMap := c.config.AWS.GetPodTags(pod.Namespace, pod.Name, instanceType)

	// Convert to EC2 tags
	tags := make([]types.Tag, 0, len(tagMap))
	for k, v := range tagMap {
		tags = append(tags, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	// Add Name tag for easy identification in console
	tags = append(tags, types.Tag{
		Key:   aws.String("Name"),
		Value: aws.String(fmt.Sprintf("orca-%s-%s", pod.Namespace, pod.Name)),
	})

	return tags
}

// convertInstance converts an EC2 instance to our internal representation.
func (c *Client) convertInstance(instance *types.Instance) *Instance {
	inst := &Instance{
		ID:    *instance.InstanceId,
		Type:  string(instance.InstanceType),
		State: string(instance.State.Name),
	}

	if instance.PublicIpAddress != nil {
		inst.PublicIP = *instance.PublicIpAddress
	}

	if instance.PrivateIpAddress != nil {
		inst.PrivateIP = *instance.PrivateIpAddress
	}

	if instance.LaunchTime != nil {
		inst.LaunchTime = *instance.LaunchTime
	}

	return inst
}
