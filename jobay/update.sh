#!/bin/bash
# Jobay update script
set -e

echo "=== Jobay Update ==="
cd /root/hermes/jobay

echo "→ Git pull..."
git pull origin main

echo "→ Build Docker image..."
docker compose build --no-cache

echo "→ Recreate container..."
docker compose up -d --force-recreate

echo "→ Cleanup..."
docker image prune -f

echo "=== Done ==="
