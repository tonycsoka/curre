#!/bin/bash
# Uses parameters passed as positional arguments

ENV="${1:-dev}"
VERBOSE="${2:-true}"

echo "Deploying to environment: $ENV"
echo "Verbose mode: $VERBOSE"
echo ""

if [ "$VERBOSE" = "true" ]; then
    echo "[VERBOSE] Connecting to deployment server..."
    echo "[VERBOSE] Authenticating..."
    echo "[VERBOSE] Uploading artifacts..."
    echo "[VERBOSE] Running health checks..."
fi

echo "Deployment complete!"
echo "Environment $ENV is now live."
