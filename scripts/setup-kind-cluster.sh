#!/bin/bash
# Setup kind cluster for ORCA development

set -e

CLUSTER_NAME="${1:-orca-dev}"

echo "Creating kind cluster: ${CLUSTER_NAME}"

# Create kind cluster with custom config
cat <<EOF | kind create cluster --name ${CLUSTER_NAME} --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
- role: worker
EOF

echo ""
echo "✅ Kind cluster created successfully!"
echo ""
echo "Next steps:"
echo "1. Verify cluster:"
echo "   kubectl cluster-info --context kind-${CLUSTER_NAME}"
echo "   kubectl get nodes"
echo ""
echo "2. Build and load ORCA image:"
echo "   make docker-build"
echo "   kind load docker-image orca:latest --name ${CLUSTER_NAME}"
echo ""
echo "3. Deploy ORCA:"
echo "   kubectl apply -f deploy/kubernetes/"
echo ""
echo "4. Test with a sample pod:"
echo "   kubectl apply -f examples/test-pod.yaml"
echo ""
