#!/bin/bash
echo "Running background API generation..."
swag init -g cmd/main.go --parseDependency --parseInternal
cd ../frontend
pnpm generate:api
echo "API generation complete."
