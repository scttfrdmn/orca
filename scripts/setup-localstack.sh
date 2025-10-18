#!/bin/bash
# Setup LocalStack for ORCA testing

set -e

echo "🐳 Setting up LocalStack for ORCA testing..."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop."
    exit 1
fi

echo "✅ Docker is running"

# Start LocalStack
echo ""
echo "Starting LocalStack container..."
docker run -d \
    --name orca-localstack \
    -p 4566:4566 \
    -p 4510-4559:4510-4559 \
    -e SERVICES=ec2 \
    -e DEBUG=1 \
    -e DOCKER_HOST=unix:///var/run/docker.sock \
    -v /var/run/docker.sock:/var/run/docker.sock \
    localstack/localstack:latest

echo "⏳ Waiting for LocalStack to be ready..."
sleep 10

# Wait for LocalStack to be healthy
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
    if curl -s http://localhost:4566/_localstack/health | grep -q "\"ec2\": \"available\""; then
        echo "✅ LocalStack is ready!"
        break
    fi
    attempt=$((attempt + 1))
    echo "   Attempt $attempt/$max_attempts..."
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "❌ LocalStack failed to start"
    exit 1
fi

# Configure AWS CLI for LocalStack
echo ""
echo "Configuring AWS CLI for LocalStack..."
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_ENDPOINT_URL=http://localhost:4566

# Create VPC
echo ""
echo "Creating VPC..."
VPC_ID=$(aws ec2 create-vpc \
    --cidr-block 10.0.0.0/16 \
    --endpoint-url http://localhost:4566 \
    --query 'Vpc.VpcId' \
    --output text)
echo "   VPC ID: $VPC_ID"

# Create subnet
echo "Creating subnet..."
SUBNET_ID=$(aws ec2 create-subnet \
    --vpc-id $VPC_ID \
    --cidr-block 10.0.1.0/24 \
    --endpoint-url http://localhost:4566 \
    --query 'Subnet.SubnetId' \
    --output text)
echo "   Subnet ID: $SUBNET_ID"

# Create security group
echo "Creating security group..."
SG_ID=$(aws ec2 create-security-group \
    --group-name orca-test-sg \
    --description "ORCA test security group" \
    --vpc-id $VPC_ID \
    --endpoint-url http://localhost:4566 \
    --query 'GroupId' \
    --output text)
echo "   Security Group ID: $SG_ID"

# Update config file
echo ""
echo "Updating config.localstack.yaml..."
sed -i.bak "s/vpcID: .*/vpcID: \"$VPC_ID\"/" config.localstack.yaml
sed -i.bak "s/subnetID: .*/subnetID: \"$SUBNET_ID\"/" config.localstack.yaml
sed -i.bak "s/- sg-12345/- \"$SG_ID\"/" config.localstack.yaml
rm config.localstack.yaml.bak

echo ""
echo "✅ LocalStack setup complete!"
echo ""
echo "Configuration:"
echo "  VPC ID:            $VPC_ID"
echo "  Subnet ID:         $SUBNET_ID"
echo "  Security Group ID: $SG_ID"
echo ""
echo "Environment variables:"
echo "  export AWS_ACCESS_KEY_ID=test"
echo "  export AWS_SECRET_ACCESS_KEY=test"
echo "  export AWS_DEFAULT_REGION=us-east-1"
echo "  export AWS_ENDPOINT_URL=http://localhost:4566"
echo ""
echo "To test ORCA with LocalStack:"
echo "  ./bin/orca --config config.localstack.yaml"
echo ""
echo "To stop LocalStack:"
echo "  docker stop orca-localstack"
echo "  docker rm orca-localstack"
