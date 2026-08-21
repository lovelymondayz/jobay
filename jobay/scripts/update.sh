#!/bin/bash
set -e
cd /root/hermes/jobay

echo "[jobay] Pulling latest code..."
git pull origin main

echo "[jobay] Stopping old containers..."
docker compose down

echo "[jobay] Building new image..."
docker compose build --no-cache

echo "[jobay] Starting new container..."
docker compose up -d

echo "[jobay] Waiting for container..."
sleep 5

echo "[jobay] Running health check..."
if curl -sf http://localhost:3010/api/health > /dev/null 2>&1; then
    echo "[jobay] ✓ Deployment successful! Dashboard: http://localhost:3010"
else
    echo "[jobay] WARNING: Health check failed. Check logs with: docker logs jobay"
    docker logs jobay --tail 20
    exit 1
fi
