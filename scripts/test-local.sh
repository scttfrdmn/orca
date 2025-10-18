#!/bin/bash
# Test ORCA locally without AWS charges

set -e

echo "🧪 Testing ORCA in mock mode (no AWS charges)"
echo ""

# Check if kind cluster exists
if ! kind get clusters | grep -q "orca-dev"; then
    echo "❌ Kind cluster 'orca-dev' not found"
    echo "Run: ./scripts/setup-kind-cluster.sh"
    exit 1
fi

# Build ORCA
echo "📦 Building ORCA..."
make build

# Run ORCA in mock AWS mode
echo ""
echo "🚀 Starting ORCA in mock mode..."
echo "   (No real AWS instances will be created)"
echo ""

./bin/orca \
    --config config.dev.yaml \
    --kubeconfig ~/.kube/config \
    --log-level debug \
    --mock-aws true

echo ""
echo "✅ ORCA test complete!"
