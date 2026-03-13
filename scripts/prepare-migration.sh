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
# Compress dataset if it's not already compressed
if [[ "$DATASET_FILE" != *.gz ]]; then
    echo "Compressing $DATASET_FILE (this might take a minute)..."
    # Use pigz if available for faster compression, otherwise fallback to gzip
    if command -v pigz &> /dev/null; then
        pigz -c "$DATASET_FILE" > "$DATASET_FILE.gz"
    else
        gzip -c "$DATASET_FILE" > "$DATASET_FILE.gz"
    fi
    UPLOAD_FILE="$DATASET_FILE.gz"
else
    UPLOAD_FILE="$DATASET_FILE"
fi

FILENAME=$(basename "$UPLOAD_FILE")
BUCKET_KEY="uploads/$(date +%s)-$FILENAME"

if command -v aws &> /dev/null; then
    echo "Using AWS CLI for upload..."
    export AWS_ACCESS_KEY_ID="$RUSTFS_ACCESS_KEY"
    export AWS_SECRET_ACCESS_KEY="$RUSTFS_SECRET_KEY"
    export AWS_DEFAULT_REGION="us-east-1"

    # Configure AWS CLI for better reliability on large files
    echo "Configuring AWS CLI for large file upload (concurrency: 10, chunksize: 32MB, timeout: 300s)..."
    aws configure set default.s3.max_concurrent_requests 10
    aws configure set default.s3.multipart_chunksize 32MB
    aws configure set default.s3.multipart_threshold 128MB
    aws configure set default.s3.max_attempts 20
    # Note: connect_timeout and read_timeout are set in config file, or use global options if supported

    # Ensure bucket exists (ignore error if it already exists)
    echo "Ensuring bucket $BUCKET_NAME exists..."
    aws s3 mb "s3://$BUCKET_NAME" --endpoint-url "$RUSTFS_URL" --no-verify-ssl || true
    
    # Calculate file size in bytes
    FILE_SIZE=$(wc -c < "$UPLOAD_FILE" | tr -d ' ')
    echo "Upload file size ($UPLOAD_FILE): $FILE_SIZE bytes"

    # Retry loop for upload
    MAX_RETRIES=5
    RETRY_COUNT=0
    SUCCESS=false
    while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
        echo "Uploading $UPLOAD_FILE (Attempt $((RETRY_COUNT+1))/$MAX_RETRIES)..."
        if aws s3 cp "$UPLOAD_FILE" "s3://$BUCKET_NAME/$BUCKET_KEY" \
            --endpoint-url "$RUSTFS_URL" \
            --no-verify-ssl \
            --cli-read-timeout 300 \
            --cli-connect-timeout 300 \
            --expected-size "$FILE_SIZE"; then
            SUCCESS=true
            break
        else
            RETRY_COUNT=$((RETRY_COUNT+1))
            echo "Upload failed. Waiting 20 seconds before retry..."
            sleep 20
        fi
    done

    if [ "$SUCCESS" = false ]; then
        echo "Error: Upload failed after $MAX_RETRIES attempts."
        # Clean up temporary gz file if we created it
        if [[ "$DATASET_FILE" != *.gz ]]; then rm -f "$UPLOAD_FILE"; fi
        exit 1
    fi
else
    echo "AWS CLI not found, using curl for upload (this might be slower for very large files)..."
    curl -v -X PUT -T "$UPLOAD_FILE" "$RUSTFS_URL/$BUCKET_NAME/$BUCKET_KEY"
fi

# Clean up temporary gz file if we created it
if [[ "$DATASET_FILE" != *.gz ]]; then
    rm -f "$UPLOAD_FILE"
fi

echo "Upload completed: $BUCKET_KEY"

echo "--- 4. Triggering Migration Job ---"
MIGRATE_RESPONSE=$(curl -s -X POST "$BACKEND_URL/admin/migrate-existing?bucketKey=$BUCKET_KEY&hasHeader=true" \
  -b "$COOKIE_JAR")

echo "Migration triggered: $MIGRATE_RESPONSE"

rm -f "$COOKIE_JAR"

echo "Done! You can monitor progress in the Admin Dashboard or via asynqmon (http://ras-pi.tail0684eb.ts.net:8090)."
