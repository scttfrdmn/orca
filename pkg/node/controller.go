package node

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/virtual-kubelet/virtual-kubelet/node"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/scttfrdmn/orca/pkg/config"
	"github.com/scttfrdmn/orca/pkg/provider"
)

// Controller manages the Virtual Kubelet node lifecycle.
type Controller struct {
	config      *config.Config
	provider    *provider.OrcaProvider
	nodeRunner  *node.NodeController
	kubeClient  kubernetes.Interface
	logger      zerolog.Logger
	version     string
	namespace   string
}

// NewController creates a new node controller.
func NewController(cfg *config.Config, kubeconfigPath, namespace, version string, logger zerolog.Logger) (*Controller, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Create ORCA provider
	logger.Info().Msg("Creating ORCA provider")
	orcaProvider, err := provider.NewProvider(cfg, cfg.Node.Name, namespace, version)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	// Create Kubernetes client
	logger.Info().Msg("Creating Kubernetes client")
	kubeConfig, err := buildKubeConfig(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &Controller{
		config:     cfg,
		provider:   orcaProvider,
		kubeClient: kubeClient,
		logger:     logger,
		version:    version,
		namespace:  namespace,
	}, nil
}

// Run starts the Virtual Kubelet node controller.
func (c *Controller) Run(ctx context.Context) error {
	c.logger.Info().
		Str("node_name", c.config.Node.Name).
		Str("namespace", c.namespace).
		Str("version", c.version).
		Msg("Starting ORCA Virtual Kubelet node")

	// Create Virtual Kubelet adapter
	adapter := NewVirtualKubeletAdapter(c.provider)

	// Create initial node object
	nodeObj := &corev1.Node{}
	nodeObj.Name = c.config.Node.Name
	adapter.ConfigureNode(ctx, nodeObj)

	// Get node interface
	nodeInterface := c.kubeClient.CoreV1().Nodes()

	// Configure node runner options
	nodeOpts := []node.NodeControllerOpt{
		node.WithNodeEnableLeaseV1(c.kubeClient.CoordinationV1().Leases(c.namespace), int32(40)),
		node.WithNodeStatusUpdateErrorHandler(func(ctx context.Context, err error) error {
			c.logger.Error().Err(err).Msg("Node status update failed")
			return err
		}),
	}

	// Create node runner
	c.logger.Info().Msg("Initializing Virtual Kubelet node controller")
	nodeRunner, err := node.NewNodeController(
		adapter,
		nodeObj,
		nodeInterface,
		nodeOpts...,
	)
	if err != nil {
		return fmt.Errorf("failed to create node controller: %w", err)
	}

	c.nodeRunner = nodeRunner

	// Start the node runner
	c.logger.Info().Msg("Starting Virtual Kubelet node controller")
	if err := c.nodeRunner.Run(ctx); err != nil {
		return fmt.Errorf("node controller error: %w", err)
	}

	c.logger.Info().Msg("Virtual Kubelet node controller stopped")
	return nil
}

// Shutdown gracefully shuts down the controller.
func (c *Controller) Shutdown(ctx context.Context) error {
	c.logger.Info().Msg("Shutting down ORCA node controller")
	// Node controller shutdown is handled by context cancellation
	return nil
}

// buildKubeConfig builds a Kubernetes REST config.
func buildKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath != "" {
		// Use specified kubeconfig file
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to default kubeconfig location
	return clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
}
