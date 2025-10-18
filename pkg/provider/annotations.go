package provider

// Pod annotations used by ORCA for instance selection and configuration.
const (
	// AnnotationInstanceType specifies the exact EC2 instance type to use.
	// This is the highest priority selection method (explicit selection).
	// Example: "p5.48xlarge", "g5.4xlarge", "t3.small"
	AnnotationInstanceType = "orca.research/instance-type"

	// AnnotationLaunchType specifies whether to use on-demand or spot instances.
	// Valid values: "on-demand", "spot"
	// Default: "on-demand"
	AnnotationLaunchType = "orca.research/launch-type"

	// AnnotationMaxSpotPrice specifies the maximum spot price to pay ($/hour).
	// Only applicable when launch-type is "spot".
	// Example: "5.00" for $5/hour max
	AnnotationMaxSpotPrice = "orca.research/max-spot-price"

	// AnnotationWorkloadTemplate specifies a named workload template to use.
	// Templates are defined in configuration and provide default instance settings.
	// Example: "llm-training", "vision-training", "inference"
	AnnotationWorkloadTemplate = "orca.research/workload-template"

	// AnnotationDebug enables debug logging for this specific pod.
	AnnotationDebug = "orca.research/debug"

	// AnnotationBudgetNamespace specifies which budget namespace to charge.
	// Used for cost allocation across departments.
	// Example: "biology-dept", "cs-dept"
	AnnotationBudgetNamespace = "orca.research/budget-namespace"

	// AnnotationMaxLifetime specifies maximum instance lifetime (duration).
	// Instance will be terminated after this duration regardless of pod status.
	// Example: "4h", "24h", "7d"
	AnnotationMaxLifetime = "orca.research/max-lifetime"

	// AnnotationAMI specifies a custom AMI to use instead of the default.
	// Example: "ami-0123456789abcdef0"
	AnnotationAMI = "orca.research/ami"

	// AnnotationUserData specifies custom EC2 user data script.
	// Must be base64 encoded.
	AnnotationUserData = "orca.research/user-data"

	// AnnotationTerminationProtection enables termination protection.
	// Valid values: "true", "false"
	// Default: "false"
	AnnotationTerminationProtection = "orca.research/termination-protection"
)

// Node labels used by ORCA.
const (
	// LabelProvider identifies nodes provided by ORCA.
	// Value: "aws"
	LabelProvider = "orca.research/provider"

	// LabelEnvironment identifies the environment (dev, staging, prod).
	LabelEnvironment = "orca.research/environment"

	// LabelVersion is the ORCA version managing this node.
	LabelVersion = "orca.research/version"
)

// Node taints used by ORCA.
const (
	// TaintKeyBurstNode prevents regular pods from being scheduled on ORCA nodes.
	// Pods must explicitly tolerate this taint to burst to AWS.
	TaintKeyBurstNode = "orca.research/burst-node"

	// TaintValueBurstNode is the value for the burst node taint.
	TaintValueBurstNode = "true"
)

// EC2 instance tags used by ORCA.
const (
	// TagPod identifies which pod this instance is running.
	// Format: "namespace/name"
	TagPod = "orca.research/pod"

	// TagPodUID is the Kubernetes pod UID.
	TagPodUID = "orca.research/pod-uid"

	// TagCluster identifies which cluster this instance belongs to.
	TagCluster = "orca.research/cluster"

	// TagProvider identifies this as an ORCA-managed instance.
	// Value: "orca"
	TagProvider = "orca.research/provider"

	// TagVersion is the ORCA version that created this instance.
	TagVersion = "orca.research/version"

	// TagNamespace is the Kubernetes namespace of the pod.
	TagNamespace = "orca.research/namespace"

	// TagBudgetNamespace is the budget namespace for cost allocation.
	TagBudgetNamespace = "orca.research/budget-namespace"

	// TagCreatedAt is when the instance was created (RFC3339).
	TagCreatedAt = "orca.research/created-at"

	// TagMaxLifetime is the maximum lifetime of the instance.
	TagMaxLifetime = "orca.research/max-lifetime"
)
