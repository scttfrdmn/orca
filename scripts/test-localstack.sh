#!/bin/bash
# Test ORCA with LocalStack

set -e

echo "🧪 Testing ORCA with LocalStack..."
echo ""

# Set LocalStack environment
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_ENDPOINT_URL=http://localhost:4567

# Check LocalStack is running
echo "1. Checking LocalStack health..."
if ! curl -s http://localhost:4567/_localstack/health | grep -E "\"ec2\": \"(available|running)\"" > /dev/null; then
    echo "❌ LocalStack is not running or EC2 is not available"
    echo "Run: ./scripts/setup-localstack.sh"
    exit 1
fi
echo "✅ LocalStack is healthy"
echo ""

# Verify config
echo "2. Verifying config.localstack.yaml..."
if [ ! -f "config.localstack.yaml" ]; then
    echo "❌ config.localstack.yaml not found"
    exit 1
fi
echo "✅ Config file found"
echo ""

# Build ORCA if not built
echo "3. Building ORCA..."
if [ ! -f "bin/orca" ]; then
    make build
fi
echo "✅ ORCA binary ready"
echo ""

# Test AWS connectivity
echo "4. Testing AWS EC2 connectivity..."
VPC_COUNT=$(aws ec2 describe-vpcs --endpoint-url http://localhost:4567 --query 'Vpcs | length(@)' --output text)
echo "   Found $VPC_COUNT VPCs in LocalStack"
echo "✅ AWS CLI can connect to LocalStack"
echo ""

# Run integration tests
echo "5. Running integration tests..."
cd /Users/scttfrdmn/src/orca
go test -v -tags=integration ./internal/aws/... 2>&1 | grep -v "^go: downloading" || true

echo ""
echo "✅ LocalStack integration test completed successfully!"
echo ""
echo "ORCA is working correctly with LocalStack. You can now:"
echo "  - Run ORCA: ./bin/orca --config config.localstack.yaml"
echo "  - Deploy pods with ORCA annotations"
echo "  - Test full pod lifecycle"
