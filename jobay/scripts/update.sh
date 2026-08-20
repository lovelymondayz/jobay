#!/bin/bash
# Jobay update script - deploys the latest version safely

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "[jobay] Starting deployment..."

# Pull latest code
echo "[jobay] Pulling latest code..."
git pull origin main

# Build new image
echo "[jobay] Building new image..."
docker compose build

# Verify build succeeded
if [ $? -ne 0 ]; then
    echo "[jobay] ERROR: Build failed! Keeping old containers running."
    exit 1
fi

# Recreate containers with new image
echo "[jobay] Recreating containers..."
docker compose up -d --force-recreate

# Health check
echo "[jobay] Running health check..."
sleep 3
if curl -sf http://localhost:3001/api/status > /dev/null 2>&1; then
    echo "[jobay] ✓ Deployment successful! Dashboard: http://localhost:3001"
else
    echo "[jobay] WARNING: Health check failed. Check logs with: docker logs jobay"
    exit 1
fi
