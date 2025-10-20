#!/bin/bash
# LocalStack initialization script for EC2 resources
# This script runs automatically when LocalStack starts

set -e

echo "Initializing LocalStack EC2 environment for ORCA..."

# Wait for LocalStack to be ready
echo "Waiting for LocalStack to be ready..."
until awslocal ec2 describe-regions > /dev/null 2>&1; do
  sleep 1
done

echo "LocalStack is ready. Setting up EC2 resources..."

# Create VPC
VPC_ID=$(awslocal ec2 create-vpc \
  --cidr-block 10.0.0.0/16 \
  --tag-specifications 'ResourceType=vpc,Tags=[{Key=Name,Value=orca-test-vpc},{Key=ManagedBy,Value=LocalStack}]' \
  --query 'Vpc.VpcId' \
  --output text)

echo "Created VPC: $VPC_ID"

# Enable DNS support
awslocal ec2 modify-vpc-attribute \
  --vpc-id "$VPC_ID" \
  --enable-dns-support

awslocal ec2 modify-vpc-attribute \
  --vpc-id "$VPC_ID" \
  --enable-dns-hostnames

# Create Internet Gateway
IGW_ID=$(awslocal ec2 create-internet-gateway \
  --tag-specifications 'ResourceType=internet-gateway,Tags=[{Key=Name,Value=orca-test-igw}]' \
  --query 'InternetGateway.InternetGatewayId' \
  --output text)

echo "Created Internet Gateway: $IGW_ID"

# Attach Internet Gateway to VPC
awslocal ec2 attach-internet-gateway \
  --vpc-id "$VPC_ID" \
  --internet-gateway-id "$IGW_ID"

# Create Subnet in us-west-2a
SUBNET_ID=$(awslocal ec2 create-subnet \
  --vpc-id "$VPC_ID" \
  --cidr-block 10.0.1.0/24 \
  --availability-zone us-west-2a \
  --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=orca-test-subnet-2a}]' \
  --query 'Subnet.SubnetId' \
  --output text)

echo "Created Subnet: $SUBNET_ID"

# Create Route Table
RTB_ID=$(awslocal ec2 create-route-table \
  --vpc-id "$VPC_ID" \
  --tag-specifications 'ResourceType=route-table,Tags=[{Key=Name,Value=orca-test-rtb}]' \
  --query 'RouteTable.RouteTableId' \
  --output text)

echo "Created Route Table: $RTB_ID"

# Create route to Internet Gateway
awslocal ec2 create-route \
  --route-table-id "$RTB_ID" \
  --destination-cidr-block 0.0.0.0/0 \
  --gateway-id "$IGW_ID"

# Associate Route Table with Subnet
awslocal ec2 associate-route-table \
  --subnet-id "$SUBNET_ID" \
  --route-table-id "$RTB_ID"

# Create Security Group
SG_ID=$(awslocal ec2 create-security-group \
  --group-name orca-test-sg \
  --description "ORCA test security group" \
  --vpc-id "$VPC_ID" \
  --tag-specifications 'ResourceType=security-group,Tags=[{Key=Name,Value=orca-test-sg},{Key=ManagedBy,Value=ORCA}]' \
  --query 'GroupId' \
  --output text)

echo "Created Security Group: $SG_ID"

# Add SSH ingress rule
awslocal ec2 authorize-security-group-ingress \
  --group-id "$SG_ID" \
  --protocol tcp \
  --port 22 \
  --cidr 0.0.0.0/0

# Add all egress (default in real AWS, but explicit in LocalStack)
awslocal ec2 authorize-security-group-egress \
  --group-id "$SG_ID" \
  --protocol -1 \
  --cidr 0.0.0.0/0 || true

# Create a dummy AMI for testing
# Note: LocalStack doesn't have real AMIs, so we create a simple one
echo "Creating test AMI..."
AMI_ID=$(awslocal ec2 register-image \
  --name "orca-test-ami" \
  --description "Test AMI for ORCA" \
  --architecture x86_64 \
  --virtualization-type hvm \
  --root-device-name /dev/xvda \
  --query 'ImageId' \
  --output text)

echo "Created AMI: $AMI_ID"

# Save resource IDs to a file for tests to use
cat > /tmp/localstack-orca-resources.env << EOF
LOCALSTACK_VPC_ID=$VPC_ID
LOCALSTACK_SUBNET_ID=$SUBNET_ID
LOCALSTACK_SG_ID=$SG_ID
LOCALSTACK_AMI_ID=$AMI_ID
EOF

echo "LocalStack EC2 initialization complete!"
echo "Resources saved to /tmp/localstack-orca-resources.env"
echo ""
echo "VPC ID:            $VPC_ID"
echo "Subnet ID:         $SUBNET_ID"
echo "Security Group ID: $SG_ID"
echo "AMI ID:            $AMI_ID"
