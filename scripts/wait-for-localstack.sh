#!/bin/bash
# Wait for LocalStack to be ready for testing

set -e

LOCALSTACK_ENDPOINT="${LOCALSTACK_ENDPOINT:-http://localhost:4566}"
MAX_ATTEMPTS=60
ATTEMPT=0

echo "Waiting for LocalStack at $LOCALSTACK_ENDPOINT..."

until curl -s "$LOCALSTACK_ENDPOINT/_localstack/health" > /dev/null 2>&1; do
    ATTEMPT=$((ATTEMPT + 1))
    if [ $ATTEMPT -ge $MAX_ATTEMPTS ]; then
        echo "ERROR: LocalStack did not become ready after $MAX_ATTEMPTS attempts"
        exit 1
    fi
    echo "Attempt $ATTEMPT/$MAX_ATTEMPTS: LocalStack not ready yet..."
    sleep 2
done

echo "LocalStack is ready!"

# Wait for EC2 service to be available
echo "Checking EC2 service availability..."
until AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
    aws --endpoint-url="$LOCALSTACK_ENDPOINT" ec2 describe-regions --region us-west-2 > /dev/null 2>&1; do
    echo "EC2 service not ready yet..."
    sleep 2
done

echo "EC2 service is available!"

# Check if resources were initialized
if [ -f /tmp/localstack-orca-resources.env ]; then
    echo "LocalStack resources initialized:"
    cat /tmp/localstack-orca-resources.env
else
    echo "Warning: Resource initialization file not found. Resources may still be initializing."
    echo "Check 'make localstack-logs' for initialization progress."
fi

echo "LocalStack is ready for testing!"
