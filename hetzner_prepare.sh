#!/bin/bash
# Hetzner Server Preparation Script
# Version: 1.0
# Date: January 2025

set -e

# Update and upgrade system
sudo apt update && sudo apt upgrade -y

# Install essential tools
sudo apt install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    ufw \
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

# Configure UFW (Uncomplicated Firewall)
sudo ufw allow OpenSSH
sudo ufw allow 'Nginx Full'
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw --force enable

# Create deployment directory
sudo mkdir -p /opt/hulta-pregnancy-app
sudo chown $USER:$USER /opt/hulta-pregnancy-app

# Setup basic security
sudo sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl restart sshd

# Install fail2ban for additional security
sudo apt install -y fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

echo "🚀 Hetzner Server Preparation Complete! 🔒"
echo "Next steps: Clone your project and deploy!"
