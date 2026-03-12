#!/bin/bash
set -e

# Ensure we are in the frontend directory
cd "$(dirname "$0")"

echo "Installing frontend dependencies..."
pnpm install

echo "Building frontend..."
pnpm run build

# Check if dist exists
if [ ! -d "dist" ]; then
  echo "Error: dist directory not found after build."
  exit 1
fi

echo "Uploading to RustFS..."
# Configure AWS CLI for S3-compatible storage
export AWS_ACCESS_KEY_ID="$RUSTFS_ACCESS_KEY"
export AWS_SECRET_ACCESS_KEY="$RUSTFS_SECRET_KEY"
export AWS_DEFAULT_REGION="us-east-1"

# Create bucket if it doesn't exist (ignore error if it exists)
aws s3 mb s3://"$BUCKET_NAME" --endpoint-url "$RUSTFS_URL" || echo "Bucket might already exist"

# Sync dist to bucket
aws s3 sync dist s3://"$BUCKET_NAME" --endpoint-url "$RUSTFS_URL" --delete

echo "Setting public bucket policy..."
# Make bucket public (read-only for anonymous users)
POLICY='{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::'"$BUCKET_NAME"'/*"]
    },
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": ["s3:ListBucket"],
      "Resource": ["arn:aws:s3:::'"$BUCKET_NAME"'"]
    }
  ]
}'
aws s3api put-bucket-policy --bucket "$BUCKET_NAME" --policy "$POLICY" --endpoint-url "$RUSTFS_URL"

echo "Frontend deployed successfully to RustFS!"
