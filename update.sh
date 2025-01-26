#!/bin/bash
# Update Script for Hulta Pregnancy App

cd /opt/hulta-pregnancy-app

# Pull latest changes
git pull origin main

# Rebuild and restart containers
docker compose down
docker compose build
docker compose up -d

# Prune old images
docker image prune -f

echo "Update completed successfully!"
