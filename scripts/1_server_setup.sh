#!/bin/bash
# Hulta Pregnancy App - Server Setup Script
# This script handles basic server setup and software installation
# Version: 1.0

# Enhanced error handling and logging
set -euo pipefail

# Logging function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [SETUP] $*" | tee -a ~/hulta_setup.log
}

error_exit() {
    log "ERROR: $1"
    echo "ERROR: $1" >&2
    exit 1
}

# Trap any errors
trap 'error_exit "Command failed: $BASH_COMMAND"' ERR

# Configuration Variables
APP_DIR="/opt/hulta-pregnancy-app"

# Check if we're running as root - changed to just a warning
if [ "$(id -u)" = "0" ]; then
    log "WARNING: This script is running as root. Some commands might behave differently."
    # Not exiting, just warning
fi

# System update and preparation
prepare_system() {
    log "🛠️ Updating system packages"
    sudo apt-get update || error_exit "Failed to update package lists"
    sudo apt-get upgrade -y || error_exit "System upgrade failed"
    
    log "📦 Installing essential system tools"
    sudo apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release \
        software-properties-common \
        git \
        wget \
        unzip \
        jq \
        || error_exit "Failed to install essential tools"
        
    log "✅ System preparation complete"
}

# Firewall configuration - Preserving SSH key authentication
configure_firewall() {
    log "🔒 Configuring UFW firewall (preserving SSH access)"
    
    # Check if UFW is already enabled
    if sudo ufw status | grep -q "Status: active"; then
        log "Firewall is already active"
        
        # Make sure SSH is allowed
        sudo ufw allow 22/tcp || error_exit "Failed to allow SSH traffic"
        log "Ensured SSH access is preserved"
    else
        # First allow SSH to prevent lockout
        sudo ufw allow 22/tcp || error_exit "Failed to allow SSH traffic"
        
        # Then configure other rules
        sudo ufw default deny incoming
        sudo ufw default allow outgoing
        sudo ufw allow 80/tcp  # HTTP
        sudo ufw allow 443/tcp # HTTPS
        
        # Enable firewall without prompt
        echo "y" | sudo ufw enable || error_exit "Failed to enable firewall"
        
        # Verify SSH is allowed
        if ! sudo ufw status | grep -q "22/tcp.*ALLOW"; then
            log "⚠️ SSH port doesn't appear to be allowed. Adding it again to be safe."
            sudo ufw allow 22/tcp
        fi
    fi
    
    log "✅ Firewall configured with SSH access preserved"
}

# Docker installation
install_docker() {
    log "🐳 Installing Docker and Docker Compose"
    
    # Check if Docker is already installed
    if command -v docker &> /dev/null; then
        log "Docker is already installed: $(docker --version)"
    else
        # Remove any conflicting packages
        for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
            sudo apt-get remove -y $pkg &> /dev/null || true
        done

        # Add Docker's official GPG key
        if [ ! -f /etc/apt/keyrings/docker.gpg ]; then
            sudo install -m 0755 -d /etc/apt/keyrings
            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
            sudo chmod a+r /etc/apt/keyrings/docker.gpg
        fi

        # Add the repository to Apt sources
        echo \
            "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
            $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
            sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

        sudo apt-get update
        sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin \
            || error_exit "Docker installation failed"
            
        # Add current user to docker group
        sudo usermod -aG docker "$USER" || error_exit "Failed to add user to docker group"
        log "⚠️ You may need to log out and back in for docker group changes to take effect"
    fi
    
    log "✅ Docker installation complete"
}

# Install Nginx
install_nginx() {
    log "🌐 Installing and configuring Nginx"
    
    if command -v nginx &> /dev/null; then
        log "Nginx is already installed: $(nginx -v 2>&1)"
    else
        sudo apt-get install -y nginx || error_exit "Failed to install Nginx"
    fi
    
    log "✅ Nginx installation complete"
}

# Create application directory
create_app_dir() {
    log "📁 Creating application directory"
    
    if [ -d "$APP_DIR" ]; then
        log "Application directory already exists"
    else
        sudo mkdir -p "$APP_DIR" || error_exit "Failed to create application directory"
        sudo chown -R "$USER:$USER" "$APP_DIR" || error_exit "Failed to set directory ownership"
    fi
    
    log "✅ Application directory created"
}

# Main function
main() {
    log "🚀 Starting server setup"
    
    prepare_system
    configure_firewall
    install_docker
    install_nginx
    create_app_dir
    
    log "✅ Server setup completed successfully!"
    log "👉 Next step: Run 2_app_deploy.sh to deploy the application"
}

# Run main function
main

exit 0
