#!/bin/bash
set -e

# Configuration
# Change these if your Tailscale URLs or ports are different
BACKEND_URL="${BACKEND_URL:-https://property.napnap.work/api}" # Or use tailscale IP if bypassing Cloudflare for API too
RUSTFS_URL="${RUSTFS_URL:-http://rustfs.property.ras-pi.tail0684eb.ts.net:9000}"
RUSTFS_ACCESS_KEY="${RUSTFS_ACCESS_KEY:-rustfsadmin}"
RUSTFS_SECRET_KEY="${RUSTFS_SECRET_KEY:-rustfsadmin}" # Use your actual secret key
BUCKET_NAME="${BUCKET_NAME:-property-data}"

USERNAME="${ADMIN_USERNAME:-admin}"
PASSWORD="${ADMIN_PASSWORD:-admin}"

DATASET_FILE="${1:-dataset.csv}"

if [ ! -f "$DATASET_FILE" ]; then
    echo "Error: Dataset file '$DATASET_FILE' not found."
    echo "Usage: $0 [path_to_dataset.csv]"
    exit 1
fi

echo "--- 1. Authenticating ---"
COOKIE_JAR=$(mktemp)
LOGIN_RESPONSE=$(curl -s -X POST "$BACKEND_URL/auth/login" \
  -H "Content-Type: application/json" \
  -c "$COOKIE_JAR" \
  -d "{\"username\": \"$USERNAME\", \"password\": \"$PASSWORD\"}")

# The backend returns "status":true in the JSON body
STATUS=$(echo $LOGIN_RESPONSE | grep -oE '"status":true' || true)

if [ -z "$STATUS" ]; then
    echo "Failed to login. Response: $LOGIN_RESPONSE"
    rm -f "$COOKIE_JAR"
    exit 1
fi

echo "Successfully authenticated."

echo "--- 2. Resetting Backend State ---"
# Note: If this returns 404, please ensure the backend is redeployed with the new /admin/reset route
RESET_RESPONSE=$(curl -s -X POST "$BACKEND_URL/admin/reset" \
  -b "$COOKIE_JAR")

echo "Reset response: $RESET_RESPONSE"

echo "--- 3. Uploading Dataset to RustFS (S3) via Tailscale ---"
# Using AWS CLI if available, otherwise fallback to curl
FILENAME=$(basename "$DATASET_FILE")
BUCKET_KEY="uploads/$(date +%s)-$FILENAME"

if command -v aws &> /dev/null; then
    echo "Using AWS CLI for upload..."
    export AWS_ACCESS_KEY_ID="$RUSTFS_ACCESS_KEY"
    export AWS_SECRET_ACCESS_KEY="$RUSTFS_SECRET_KEY"
    export AWS_DEFAULT_REGION="us-east-1"

    # Ensure bucket exists (ignore error if it already exists)
    echo "Ensuring bucket $BUCKET_NAME exists..."
    aws s3 mb "s3://$BUCKET_NAME" --endpoint-url "$RUSTFS_URL" --no-verify-ssl || true
    
    aws s3 cp "$DATASET_FILE" "s3://$BUCKET_NAME/$BUCKET_KEY" \
        --endpoint-url "$RUSTFS_URL" \
        --no-verify-ssl
else
    echo "AWS CLI not found, using curl for upload (this might be slower for very large files)..."
    curl -v -X PUT -T "$DATASET_FILE" "$RUSTFS_URL/$BUCKET_NAME/$BUCKET_KEY"
fi

echo "Upload completed: $BUCKET_KEY"

echo "--- 4. Triggering Migration Job ---"
MIGRATE_RESPONSE=$(curl -s -X POST "$BACKEND_URL/admin/migrate-existing?bucketKey=$BUCKET_KEY&hasHeader=true" \
  -b "$COOKIE_JAR")

echo "Migration triggered: $MIGRATE_RESPONSE"

rm -f "$COOKIE_JAR"

echo "Done! You can monitor progress in the Admin Dashboard or via asynqmon (http://ras-pi.tail0684eb.ts.net:8090)."
