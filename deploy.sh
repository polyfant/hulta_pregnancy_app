#!/bin/bash

# Exit on any error
set -e

# Pull latest changes
git pull origin main

# Build and deploy containers
docker-compose down
docker-compose build
docker-compose up -d

# Optional: Prune old images to save space
docker image prune -f

echo "Deployment completed successfully! 🚀"
