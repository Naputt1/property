#!/bin/bash
echo "Running background API generation..."
swag init -g cmd/main.go --parseDependency --parseInternal > /dev/null
cd ../frontend
pnpm --silent generate:api > /dev/null
echo "API generation complete."
