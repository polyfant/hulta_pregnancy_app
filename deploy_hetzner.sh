#!/bin/bash

# Hetzner Deployment Script for Hulta Pregnancy App
# Version: 1.0
# Date: January 2025

# Exit on any error
set -e

# Server Preparation
echo " Preparing Hetzner Server for Deployment"

# Update and upgrade system
sudo apt update && sudo apt upgrade -y

# Install essential tools
sudo apt install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    nginx \
    certbot \
    python3-certbot-nginx

# Install Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Add current user to docker group
sudo usermod -aG docker $USER

# Create project directory
mkdir -p /opt/hulta-pregnancy-app

# Clone project (replace with your git clone command)
git clone https://github.com/yourusername/hulta_pregnancy_app.git /opt/hulta-pregnancy-app

# Navigate to project directory
cd /opt/hulta-pregnancy-app

# Copy environment files (ensure these are securely managed)
cp .env.example .env
cp frontend-react/.env.example frontend-react/.env

# Generate SSL Certificate
sudo certbot certonly --nginx --standalone

# Configure Nginx as Reverse Proxy
sudo tee /etc/nginx/sites-available/hulta-pregnancy <<EOF
server {
    listen 80;
    server_name _;

    location / {
        proxy_pass http://localhost:3000;  # Frontend
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }

    location /api {
        proxy_pass http://localhost:8080;  # Backend
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF

sudo ln -s /etc/nginx/sites-available/hulta-pregnancy /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Run Docker Compose
docker compose up -d

# Setup automatic updates and backups
(crontab -l 2>/dev/null; echo "0 2 * * * /opt/hulta-pregnancy-app/backup.sh") | crontab -
(crontab -l 2>/dev/null; echo "0 3 * * 0 /opt/hulta-pregnancy-app/update.sh") | crontab -

echo " Deployment Complete! "
